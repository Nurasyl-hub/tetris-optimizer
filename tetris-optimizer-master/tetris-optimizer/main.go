package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type Point struct {
	x, y int
}

func main() {
	if len(os.Args) != 2 {
		return
	}

	tetrominoes, err := parseFile(os.Args[1])
	if err != nil {
		fmt.Println("ERROR")
		return
	}

	size := int(math.Ceil(math.Sqrt(float64(len(tetrominoes) * 4))))

	for {
		board := make([][]rune, size)
		for i := range board {
			board[i] = make([]rune, size)
			for j := range board[i] {
				board[i][j] = '.'
			}
		}

		if solve(board, tetrominoes, 0) {
			for _, row := range board {
				fmt.Println(string(row))
			}
			return
		}
		size++
	}
}

func parseFile(filename string) ([][]Point, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	content := string(data)
	content = strings.ReplaceAll(content, "\r\n", "\n")
	
	chunks := strings.Split(content, "\n\n")
	var shapes [][]Point

	for _, chunk := range chunks {
		if strings.TrimSpace(chunk) == "" {
			continue
		}
		
		lines := strings.Split(strings.Trim(chunk, "\n"), "\n")
		if len(lines) != 4 {
			return nil, fmt.Errorf("bad format")
		}

		var shape []Point
		connections := 0
		for y, line := range lines {
			if len(line) != 4 {
				return nil, fmt.Errorf("bad width")
			}
			for x, char := range line {
				if char == '#' {
					shape = append(shape, Point{x, y})
					if x > 0 && line[x-1] == '#' { connections++ }
					if x < 3 && line[x+1] == '#' { connections++ }
					if y > 0 && lines[y-1][x] == '#' { connections++ }
					if y < 3 && lines[y+1][x] == '#' { connections++ }
				}
			}
		}

		if len(shape) != 4 || (connections != 6 && connections != 8) {
			return nil, fmt.Errorf("invalid shape")
		}
		shapes = append(shapes, normalize(shape))
	}

	if len(shapes) == 0 {
		return nil, fmt.Errorf("no shapes")
	}
	return shapes, nil
}

func normalize(shape []Point) []Point {
	minX, minY := 4, 4
	for _, p := range shape {
		if p.x < minX { minX = p.x }
		if p.y < minY { minY = p.y }
	}
	res := make([]Point, 4)
	for i, p := range shape {
		res[i] = Point{p.x - minX, p.y - minY}
	}
	return res
}

func solve(board [][]rune, shapes [][]Point, index int) bool {
	if index == len(shapes) {
		return true
	}

	size := len(board)
	for y := 0; y <= size; y++ {
		for x := 0; x <= size; x++ {
			if canPlace(board, shapes[index], x, y) {
				place(board, shapes[index], x, y, rune('A'+index))
				if solve(board, shapes, index+1) {
					return true
				}
				place(board, shapes[index], x, y, '.')
			}
		}
	}
	return false
}

func canPlace(board [][]rune, shape []Point, x, y int) bool {
	size := len(board)
	for _, p := range shape {
		nx, ny := x+p.x, y+p.y
		if nx < 0 || nx >= size || ny < 0 || ny >= size || board[ny][nx] != '.' {
			return false
		}
	}
	return true
}

func place(board [][]rune, shape []Point, x, y int, char rune) {
	for _, p := range shape {
		board[y+p.y][x+p.x] = char
	}
}