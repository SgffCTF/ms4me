package game

import (
	"errors"
	"math/rand"
)

const FIELD_SIZE = 8
const MINE_COUNT = 10

var (
	ErrAlreadyOpen    = errors.New("Клетка уже открыта")
	ErrFlagOnOpenCell = errors.New("Нельзя поставить флаг на открытую клетку")
)

type Field struct {
	Rows       int       `json:"rows"`
	Cols       int       `json:"cols"`
	Mines      int       `json:"mines"`
	CellsOpen  int       `json:"cells_open"`
	MineIsOpen bool      `json:"mine_is_open"`
	Grid       [][]*Cell `json:"grid"`
}

func NewField() *Field {
	return &Field{
		Rows:       FIELD_SIZE,
		Cols:       FIELD_SIZE,
		Mines:      MINE_COUNT,
		CellsOpen:  0,
		MineIsOpen: false,
	}
}

func (f *Field) OpenCell(row, col int) error {
	if f.Grid[row][col].IsOpen {
		return ErrAlreadyOpen
	}

	if f.Grid[row][col].HasMine {
		f.MineIsOpen = true
		return nil
	}

	f.openNeighborCells(row, col)
	return nil
}

// openNeighborCells рекурсивно открывает соседние клетки
func (f *Field) openNeighborCells(row, col int) {
	if row < 0 || col < 0 || row >= FIELD_SIZE || col >= FIELD_SIZE {
		return
	}

	cell := f.Grid[row][col]

	if cell.IsOpen || cell.HasMine {
		return
	}

	cell.IsOpen = true
	f.CellsOpen++

	if cell.NeighborMines > 0 {
		return
	}

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			f.openNeighborCells(row+i, col+j)
		}
	}
}

func (f *Field) SetFlag(row int, col int) error {
	if f.Grid[row][col].IsOpen {
		return ErrFlagOnOpenCell
	}
	f.Grid[row][col].Value = FLAG
	return nil
}

// CalculateFieldNeighborMines подсчитывает количество соседних мин в каждой клетке
func (f *Field) calculateFieldNeighborMines() {
	for row := 0; row < FIELD_SIZE; row++ {
		for col := 0; col < FIELD_SIZE; col++ {
			f.Grid[row][col].NeighborMines = f.calculateNeighborMines(row, col)
		}
	}
}

// calculateNeighborMines подсчитывает количество соседних мин в одной клетке
func (f *Field) calculateNeighborMines(row, col int) int {
	c := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if row+i < 0 || col+j < 0 || row+i >= FIELD_SIZE || col+j >= FIELD_SIZE {
				continue
			}

			c++
		}
	}
	return c
}

// CreateGrid создаёт игровое поле
// firstRow, firstCol необходимы для того, чтобы генерировать поле после первого нажатия
func CreateField(firstRow, firstCol int) *Field {
	f := CreateClosedField()
	for i := 0; i < MINE_COUNT; i++ {
		x, y := rand.Intn(8), rand.Intn(8)
		for f.Grid[x][y].Value == MINE || (firstRow-1 <= x && x <= firstRow+1 && firstCol-1 <= y && y <= firstCol+1) {
			x, y = rand.Intn(8), rand.Intn(8)
		}
		f.Grid[x][y].HasMine = true
	}
	f.calculateFieldNeighborMines()
	return f
}

// CreateClosedField создаёт закрытое игровое поле без мин
func CreateClosedField() *Field {
	f := NewField()
	Grid := make([][]*Cell, FIELD_SIZE)
	for i := 0; i < FIELD_SIZE; i++ {
		row := make([]*Cell, FIELD_SIZE)
		for j := 0; j < FIELD_SIZE; j++ {
			row[j] = NewCell(CLOSED)
		}
		Grid[i] = row
	}
	f.Grid = Grid
	return f
}
