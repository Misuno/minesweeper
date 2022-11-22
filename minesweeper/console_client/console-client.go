package consoleclient

import (
	"bufio"
	"fmt"
	"minesweeper/minesweeper"
	"strconv"
	"strings"
)

const (
	drawFormat = " %v |"
	drawHline  = "----"
)

var drawables = map[int]string{
	int(minesweeper.StateBomb):    "☢︎",
	int(minesweeper.StateEmpty):   " ",
	int(minesweeper.StateUnknown): "⎕",
	int(minesweeper.StateFlagged): "☂︎",
}

func DrawBoardOpen(board *minesweeper.Board) {
	cells := board.GetCells()
	for _, r := range *cells {
		for _, c := range r {
			c.Open()
		}
	}
	DrawBoard(board)
}

func DrawBoard(board *minesweeper.Board) {
	var numbers [][]int
	cells := board.GetBoard()

	for i, row := range cells {
		// First column - line numbers
		if i == 0 {
			var frstCol []int
			frstCol = append(frstCol, 0) // blank space for top left corner
			for j := 0; j < len(row); j++ {
				frstCol = append(frstCol, j)
			}
			numbers = append(numbers, frstCol)
		}

		// Other colums - first element is col number
		col := []int{i}
		for _, s := range row {
			col = append(col, int(s))
		}
		numbers = append(numbers, col)
	}

	for i := 0; i < len(numbers); i++ {
		for j := 0; j < len(numbers[i]); j++ {
			// Draw string
			s, ok := drawables[numbers[i][j]]
			if !ok {
				s = strconv.Itoa(numbers[i][j])
			}
			format := drawFormat
			if i*j == 0 {
				format = "(%v)|"
			}
			fmt.Printf(format, s)

			// If end - draw line
			if j == len(numbers[i])-1 {
				fmt.Println()
				for k := 0; k <= j; k++ {
					fmt.Print(drawHline)
				}
				fmt.Println()
			}
		}
	}
}

func ProcessInput(reader *bufio.Reader, board *minesweeper.Board) (minesweeper.State, error) {
	l, _ := reader.ReadString('\n')
	s := strings.Split(l, " ")
	if len(s) < 3 {
		return minesweeper.StateUnknown, fmt.Errorf("bad input")
	}

	for i := 0; i < len(s); i++ {
		s[i] = strings.Trim(s[i], " \n")
	}

	y, _ := strconv.Atoi(s[1])
	x, _ := strconv.Atoi(s[2])

	coord := minesweeper.Coord{X: x, Y: y}
	var result minesweeper.State
	switch s[0] {
	case "o", "O", "о", "О":
		result, _ = board.OpenCellCoord(coord)
	case "f", "F":
		result = board.FlagCellCoord(coord)
	default:
		return minesweeper.StateUnknown, fmt.Errorf("bad input")
	}

	return result, nil
}
