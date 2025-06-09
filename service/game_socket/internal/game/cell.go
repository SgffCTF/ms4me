package game

type CellType rune

const (
	MINE    CellType = 'm'
	EMPTY   CellType = '0'
	CLOSED  CellType = 'c'
	FLAG    CellType = 'f'
	MINES_1 CellType = '1'
	MINES_2 CellType = '2'
	MINES_3 CellType = '3'
	MINES_4 CellType = '4'
	MINES_5 CellType = '5'
	MINES_6 CellType = '6'
	MINES_7 CellType = '7'
	MINES_8 CellType = '8'
)

type Cell struct {
	value         CellType
	neighborMines int // кол-во мин по соседству
}

func NewCell(value CellType) *Cell {
	return &Cell{
		value:         value,
		neighborMines: 0,
	}
}

func (c *Cell) IsOpen() bool {
	return c.value != CLOSED
}

func (c *Cell) IsMine() bool {
	return c.value == MINE
}

func (c *Cell) Set(value CellType) {
	c.value = value
}
