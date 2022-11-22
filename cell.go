package minesweeper

import (
	"fmt"
	"strconv"
)

type State int8

const (
	StateEmpty   State = 0
	StateBomb    State = -1
	StateUnknown State = -2
	StateFlagged State = -3
)

type Cell struct {
	state      State
	isOpen     bool
	isFlagged  bool
	neighbours []*Cell
}

func MakeCell() Cell {
	return Cell{
		state:      StateEmpty,
		isOpen:     false,
		isFlagged:  false,
		neighbours: make([]*Cell, 0),
	}
}

func (c *Cell) GetState() State {
	switch {
	case c.isFlagged:
		return StateFlagged
	case c.isOpen:
		return c.state
	default:
		return StateUnknown
	}
}

func (c *Cell) IsOpen() bool {
	return c.isOpen
}

func (c *Cell) Open() State {
	c.isOpen = true
	return c.state
}

func (c *Cell) FlagSwitch() State {
	if c.isOpen {
		return c.state
	}
	c.isFlagged = !c.isFlagged
	return StateFlagged
}

func (c *Cell) Draw(debug ...bool) {
	var dbg bool
	if len(debug) == 0 {
		dbg = false
	} else {
		dbg = debug[0]
	}

	var sym string
	if dbg || c.isOpen {
		switch c.state {
		case StateEmpty:
			sym = " "
		case StateBomb:
			sym = "*"
		default:
			sym = strconv.Itoa(int(c.state))
		}
	} else {
		sym = "◻︎"
	}
	fmt.Printf(" %s ", sym)
}
