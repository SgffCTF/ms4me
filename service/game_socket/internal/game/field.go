package game

import "math/rand"

const FIELD_SIZE = 8
const MINE_COUNT = 10

type Field struct {
	rows       int
	cols       int
	mines      int
	cellsOpen  int
	mineIsOpen bool
	grid       [][]*Cell
}

func New() *Field {
	return &Field{
		rows:       FIELD_SIZE,
		cols:       FIELD_SIZE,
		mines:      MINE_COUNT,
		cellsOpen:  0,
		mineIsOpen: false,
	}
}

// CreateGrid создаёт игровое поле
// firstRow, firstCol необходимы для того, чтобы генерировать поле после первого нажатия
func (f *Field) CreateGrid(firstRow, firstCol int) {
	f.createClosedGrid()
	for i := 0; i < MINE_COUNT; i++ {
		x, y := rand.Intn(8), rand.Intn(8)
		for f.grid[x][y].IsMine() || (firstRow-1 <= x && x <= firstRow+1 && firstCol-1 <= y && y <= firstCol+1) {
			x, y = rand.Intn(8), rand.Intn(8)
		}
		f.grid[x][y].Set(MINE)
	}
}

func (f *Field) createClosedGrid() {
	grid := make([][]*Cell, FIELD_SIZE)
	for i := 0; i < FIELD_SIZE; i++ {
		row := make([]*Cell, FIELD_SIZE)
		for j := 0; j < FIELD_SIZE; j++ {
			row[j] = NewCell(CLOSED)
		}
		grid[i] = row
	}
	f.grid = grid
}
