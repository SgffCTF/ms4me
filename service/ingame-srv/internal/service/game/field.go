package game

import (
	"errors"
	"math/rand"
)

const fieldSize = 8
const mineCount = 10

var (
	ErrAlreadyOpen    = errors.New("Клетка уже открыта")
	ErrFlagOnOpenCell = errors.New("Нельзя поставить флаг на открытую клетку")
	ErrFieldSize      = errors.New("Выход за пределы поля")
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
		Rows:       fieldSize,
		Cols:       fieldSize,
		Mines:      mineCount,
		CellsOpen:  0,
		MineIsOpen: false,
	}
}

func (f *Field) OpenCell(row, col int) error {
	if row < 0 || row >= f.Rows || col < 0 || col >= f.Cols {
		return ErrFieldSize
	}

	if f.Grid[row][col].IsOpen {
		f.openCellsAround(row, col)
		return nil
	}

	if f.Grid[row][col].IsMine() {
		f.Grid[row][col].IsOpen = true
		f.MineIsOpen = true
		f.Grid[row][col].SetOpenValue()
		return nil
	}

	f.openNeighborCells(row, col)
	return nil
}

func (f *Field) openCellsAround(row, col int) {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			cellRow := row + i
			cellCol := col + j
			if cellRow < 0 || cellRow >= f.Rows || cellCol < 0 || cellCol >= f.Cols {
				continue
			}

			cell := f.Grid[cellRow][cellCol]
			if cell.IsOpen == false && cell.Value == CLOSED {
				if cell.IsMine() {
					cell.IsOpen = true
					f.MineIsOpen = true
					cell.SetOpenValue()
				} else {
					f.openNeighborCells(cellRow, cellCol)
				}
			}
		}
	}
}

func (f *Field) IsWin() bool {
	totalCells := f.Rows * f.Cols
	return f.CellsOpen == totalCells-f.Mines && !f.MineIsOpen
}

// openNeighborCells рекурсивно открывает соседние клетки
func (f *Field) openNeighborCells(row, col int) {
	if row < 0 || col < 0 || row >= f.Rows || col >= f.Cols {
		return
	}

	cell := f.Grid[row][col]

	if cell.IsOpen || cell.IsMine() {
		return
	}

	cell.IsOpen = true
	cell.SetOpenValue()
	f.CellsOpen++

	if cell.NeighborMines > 0 {
		return
	}

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			f.openNeighborCells(row+i, col+j)
		}
	}
}

func (f *Field) SetFlag(row int, col int) error {
	if row < 0 || row >= f.Rows || col < 0 || col >= f.Cols {
		return ErrFieldSize
	}

	cell := f.Grid[row][col]
	if cell.IsOpen {
		return ErrFlagOnOpenCell
	}
	if cell.Value == FLAG {
		cell.Value = CLOSED
	} else {
		cell.Value = FLAG
	}
	return nil
}

// CalculateFieldNeighborMines подсчитывает количество соседних мин в каждой клетке
func (f *Field) calculateFieldNeighborMines() {
	for row := 0; row < f.Rows; row++ {
		for col := 0; col < f.Cols; col++ {
			f.Grid[row][col].NeighborMines = f.calculateNeighborMines(row, col)
		}
	}
}

// calculateNeighborMines подсчитывает количество соседних мин в одной клетке
func (f *Field) calculateNeighborMines(row, col int) int {
	c := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if row+i < 0 || col+j < 0 || row+i >= fieldSize || col+j >= fieldSize || (i == 0 && j == 0) {
				continue
			}
			if f.Grid[row+i][col+j].IsMine() {
				c++
			}
		}
	}
	return c
}

// CreateGrid создаёт игровое поле
// firstRow, firstCol необходимы для того, чтобы генерировать поле после первого нажатия
func CreateField(firstRow, firstCol int) *Field {
	f := CreateClosedField()
	for i := 0; i < f.Mines; i++ {
		x, y := rand.Intn(f.Rows), rand.Intn(f.Cols)
		for f.Grid[x][y].Value == MINE || (firstRow-1 <= x && x <= firstRow+1 && firstCol-1 <= y && y <= firstCol+1) {
			x, y = rand.Intn(f.Rows), rand.Intn(f.Cols)
		}
		f.Grid[x][y].SetMine()
	}
	f.calculateFieldNeighborMines()
	return f
}

// CreateClosedField создаёт закрытое игровое поле без мин
func CreateClosedField() *Field {
	f := NewField()
	Grid := make([][]*Cell, f.Rows)
	for i := 0; i < f.Rows; i++ {
		row := make([]*Cell, f.Cols)
		for j := 0; j < f.Cols; j++ {
			row[j] = NewCell(CLOSED)
		}
		Grid[i] = row
	}
	f.Grid = Grid
	return f
}
