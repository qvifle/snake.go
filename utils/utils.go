package utils

import "fmt"

var SquareSymbol rune = '■'
var EmptySymbol rune = '_'
var BonusSymbol = '@'

func ClearLines(count int) {
	for i := 0; i < count; i++ {
		fmt.Printf("\033[1A\033[K")
	}
}
func initMatrix(len int) [][]rune {
	matrix := make([][]rune, len)
	for i := 0; i < len; i++ {
		matrix[i] = make([]rune, len)
	}

	return matrix
}

func fillMatrix(matrix [][]rune, len int) {
	for i := 0; i < len; i++ {
		for j := 0; j < len; j++ {
			matrix[i][j] = EmptySymbol
		}
	}
}

func CreateMap(s int) [][]rune {
	matrix := initMatrix(s)
	fillMatrix(matrix, s)
	return matrix
}
