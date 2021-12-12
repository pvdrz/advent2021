package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func pop(stack []Point) (Point, []Point) {
	new_len := len(stack) - 1
	return stack[new_len], stack[:new_len]
}

type Point struct {
	x uint
	y uint
}

type Grid struct {
	rows   [][]uint8
	n_cols uint
}

func parseGrid(input string) (Grid, error) {
	input = strings.Trim(input, "\n")
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")

	if len(lines) == 0 {
		return (Grid{}), fmt.Errorf("input is empty")
	}

	n_cols := len(lines[0])

	rows := make([][]uint8, 0)
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) != n_cols {
			return (Grid{}), fmt.Errorf("line %d has an invalid length: expected %d, found %d", i, n_cols, len(line))
		}

		row := make([]uint8, n_cols)
		for i, char := range line {
			digit, err := strconv.ParseUint(string(char), 10, 8)
			if err != nil {
				return (Grid{}), err
			}
			row[i] = uint8(digit)
		}
		rows = append(rows, row)
	}

	return Grid{rows: rows, n_cols: uint(n_cols)}, nil
}

func (grid *Grid) nCols() uint {
	return grid.n_cols
}

func (grid *Grid) nRows() uint {
	return uint(len(grid.rows))
}

func (grid *Grid) contains(i uint, j uint) bool {
	return i < grid.nCols() && j < grid.nRows()
}

func (grid *Grid) get(i uint, j uint) uint8 {
	return grid.rows[j][i]
}

func (grid *Grid) set(i uint, j uint, val uint8) {
	grid.rows[j][i] = val
}

func (grid *Grid) forEachNeighbor(i uint, j uint, f func(uint, uint)) {
	up := j > 0 && grid.contains(i, j-1)
	left := i > 0 && grid.contains(i-1, j)
	down := j < math.MaxUint && grid.contains(i, j+1)
	right := i < math.MaxUint && grid.contains(i+1, j)
	if up {
		f(i, j-1)
		if left {
			f(i-1, j-1)
		}
		if right {
			f(i+1, j-1)
		}
	}
	if left {
		f(i-1, j)
	}
	if down {
		f(i, j+1)
		if left {
			f(i-1, j+1)
		}
		if right {
			f(i+1, j+1)
		}
	}
	if right {
		f(i+1, j)
	}
}

func (grid *Grid) evolve() uint {
	flashing := make([]Point, 0)
	for j := uint(0); j < grid.nRows(); j += 1 {
		for i := uint(0); i < grid.nCols(); i += 1 {
			val := (grid.get(i, j) + 1) % 10
			grid.set(i, j, val)
			if val == 0 {
				flashing = append(flashing, Point{x: i, y: j})
			}
		}
	}

	count := uint(len(flashing))
	var point Point
	for len(flashing) > 0 {
		point, flashing = pop(flashing)
		grid.forEachNeighbor(point.x, point.y, func(i, j uint) {
			val := grid.get(i, j)
			if val != 0 {
				val = (val + 1) % 10
				grid.set(i, j, val)
				if val == 0 {
					flashing = append(flashing, Point{x: i, y: j})
					count += 1
				}
			}
		})
	}

	return count
}

func (grid *Grid) run(steps int) uint {
	count := uint(0)
	for i := 0; i < steps; i += 1 {
		count += grid.evolve()
	}
	return count
}

func (grid *Grid) runUntilSync() uint {
	step := uint(0)
	size := grid.nRows() * grid.nCols()
	for {
		count := grid.evolve()
		step += 1

		if count == size {
			return step
		}
	}
}

func (grid *Grid) equals(other *Grid) bool {
	if grid.nRows() != other.nRows() || grid.nCols() != other.nCols() {
		return false
	}

	for j := uint(0); j < grid.nRows(); j += 1 {
		for i := uint(0); i < grid.nCols(); i += 1 {
			if grid.get(i, j) != other.get(i, j) {
				return false
			}
		}
	}

	return true
}

func main() {
	bytes, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	input := strings.Trim(string(bytes), "\n")
	grid, err := parseGrid(input)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 1: %d\n", grid.run(100))
	fmt.Printf("Part 2: %d\n", 100+grid.runUntilSync())
}
