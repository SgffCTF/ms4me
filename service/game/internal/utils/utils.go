package utils

func MineFunc(rows, cols int) int {
	return int((float64(rows) * float64(cols) / 100) * 16) // count of mines should be 16% of field
}
