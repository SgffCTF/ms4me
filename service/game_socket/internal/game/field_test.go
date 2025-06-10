package game

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestField(t *testing.T) {
	testCases := []struct {
		name   string
		field  *Field
		result *Field
	}{
		{
			name:   "",
			field:  readGrid("testcases/test_001.json"),
			result: readGrid("testcases/result_001.json"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.field.OpenCell(0, 0)
			require.Equal(t, tc.field.CellsOpen, tc.result.CellsOpen)
			for row := 0; row < FIELD_SIZE; row++ {
				for col := 0; col < FIELD_SIZE; col++ {
					require.Equal(t, tc.result.Grid[row][col].IsOpen, tc.field.Grid[row][col].IsOpen)
					if tc.field.Grid[row][col].IsOpen {
						fmt.Print("* ")
						continue
					}
					if tc.field.Grid[row][col].HasMine {
						fmt.Print("1 ")
					} else {
						fmt.Print("0 ")
					}
				}
				fmt.Println()
			}
		})
	}
}

func readGrid(file string) *Field {
	data, err := os.ReadFile(file)
	if err != nil {
		panic("can't read testcase file: " + file)
	}
	var f Field
	err = json.Unmarshal(data, &f)
	if err != nil {
		panic("can't unmarshal testcase file: " + err.Error())
	}
	return &f
}
