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
	Value         CellType `json:"value"`
	NeighborMines int      `json:"neighbor_mines"` // кол-во мин по соседству
	IsOpen        bool     `json:"is_open"`
	HasMine       bool     `json:"has_mine"`
}

func NewCell(value CellType) *Cell {
	return &Cell{
		Value:         value,
		NeighborMines: 0,
		IsOpen:        false,
		HasMine:       false,
	}
}
