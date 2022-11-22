package minesweeper

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Board struct {
	cells [][]*Cell
}

func MakeBoard(width, height, bombs int) *Board {
	b := Board{}
	for i := 0; i < width; i++ {
		col := make([]*Cell, 0)
		for j := 0; j < height; j++ {
			cell := MakeCell()
			col = append(col, &cell)
		}
		b.cells = append(b.cells, col)
	}
	b.PutBombs(bombs)
	b.PutNumbers()
	return &b
}

func (b *Board) GetCells() *[][]*Cell {
	return &b.cells
}

func (b *Board) GetBoard() [][]State {
	var res [][]State
	for i, col := range b.cells {
		res = append(res, make([]State, 0))
		for _, cell := range col {
			c := cell.GetState()
			res[i] = append(res[i], c)
		}
	}
	return res
}

func (b *Board) IsAllOpened() bool {
	counter := 0
	for _, r := range b.cells {
		for _, c := range r {
			if c.isFlagged || c.isOpen {
				counter++
			}
		}
	}
	return counter == len(b.cells)*len(b.cells[0])
}

func (b *Board) FlagCellCoord(c Coord) State {
	return b.cells[c.X][c.Y].FlagSwitch()
}

func (b *Board) OpenCellCoord(c Coord) (State, error) {
	return b.OpenCell(c.X, c.Y)
}

func (b *Board) OpenCell(x, y int) (State, error) {
	c := b.cells[x][y]
	if c.isFlagged {
		return 0, errors.New("Cell is flagged")
	}

	if c.state == StateEmpty { // If empty open all adjacent empties and numbers
		// First get all blank cells and open them
		var blankCells []*Cell
		blankCells = append(blankCells, b.cells[x][y])

		i := 0
		for {
			n := blankCells[i].neighbours
			n = filter(n, func(a *Cell) bool {
				return a.state == StateEmpty
			})

			n = filter(n, func(a *Cell) bool {
				return contains(blankCells, a) < 0
			})

			blankCells = append(blankCells, n...)

			if i == len(blankCells)-1 {
				break
			}

			i++
		}

		// Then get and open all cells with numbers
		numbers := make([]*Cell, 0)
		for _, bl := range blankCells {
			bl.Open()
			for _, nbr := range bl.neighbours {
				if nbr.state > StateEmpty {
					numbers = append(numbers, nbr)
				}
			}
		}
		for _, n := range numbers {
			n.Open()
		}
	}

	return c.Open(), nil
}

func (b *Board) Draw(debug ...bool) {
	for _, row := range b.cells {
		for _, cell := range row {
			cell.Draw(debug...)
		}
		fmt.Println()
	}
}

func (b *Board) PutBombs(count int) error {
	rows := len(b.cells)
	if rows == 0 {
		return errors.New("empty board")
	}

	cols := len(b.cells[0])
	if cols == 0 {
		return errors.New("empty column")
	}

	if count >= rows*cols {
		return errors.New("too many bombs")
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		for {
			x := rand.Intn(rows)
			y := rand.Intn(cols)
			c := b.cells[x][y]
			if c.state != StateBomb {
				c.state = StateBomb
				break
			}
		}
	}

	return nil
}

func (b *Board) PutNumbers() error {
	maxX := len(b.cells) - 1
	for x := 0; x <= maxX; x++ {
		maxY := len(b.cells[x]) - 1
		for y := 0; y <= maxY; y++ {
			cell := b.cells[x][y]
			if cell.state != StateBomb {
				cell.state = b.GetNumberForCell(x, y, maxX, maxY)
			}
		}
	}

	return nil
}

func (b *Board) GetNumberForCell(x, y, maxX, maxY int) State {
	bombsCount := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			i := x + dx
			j := y + dy

			if isInvalidIndex(i, j, maxX, maxY) {
				continue
			}

			if b.cells[i][j].state == StateBomb {
				bombsCount++
			}

			b.cells[x][y].neighbours = append(b.cells[x][y].neighbours, b.cells[i][j])
		}
	}

	return State(bombsCount)
}

func (b *Board) GetAdjacentBlanks(x, y, maxX, maxY int) []Coord {
	var result []Coord
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if (dx == 0 && dy == 0) || dx*dy != 0 { // If base cell or diagonal cell do nothing
				continue
			}

			i := x + dx
			j := y + dy

			if isInvalidIndex(i, j, maxX, maxY) {
				continue
			}

			if b.cells[i][j].state == StateEmpty {
				result = append(result, Coord{X: i, Y: j})
			}
		}
	}

	return result
}

func isInvalidIndex(i int, j int, maxX int, maxY int) bool {
	return i < 0 || j < 0 || i > maxX || j > maxY
}

func contains[T comparable](s []T, v T) int {
	for i, a := range s {
		if a == v {
			return i
		}
	}
	return -1
}

func filter[T comparable](s []T, f func(T) bool) []T {
	var res []T
	for _, v := range s {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}
