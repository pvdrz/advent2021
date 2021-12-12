package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
)

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
	lines := strings.Split(input, "\n")

	if len(lines) == 0 {
		return (Grid{}), fmt.Errorf("input is empty")
	}

	n_cols := len(lines[0])

	rows := make([][]uint8, 0)
	for i, line := range lines {
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

func (grid *Grid) forEachNeighbor(i uint, j uint, f func(uint, uint)) {
	if j > 0 && grid.contains(i, j-1) {
		f(i, j-1)
	}
	if i > 0 && grid.contains(i-1, j) {
		f(i-1, j)
	}
	if j < math.MaxUint && grid.contains(i, j+1) {
		f(i, j+1)
	}
	if i < math.MaxUint && grid.contains(i+1, j) {
		f(i+1, j)
	}
}

func (grid *Grid) getMins() []Point {
	mins := make([]Point, 0)

	for j := uint(0); j < grid.nRows(); j += 1 {
		for i := uint(0); i < grid.nCols(); i += 1 {
			value := grid.get(i, j)
			is_min := true

			grid.forEachNeighbor(i, j, func(x uint, y uint) {
				is_min = is_min && value < grid.get(x, y)
			})

			if is_min {
				mins = append(mins, Point{x: i, y: j})
			}
		}
	}

	return mins
}

func (grid *Grid) getBasins() [][]Point {
	basins := make([][]Point, 0)

	visited := map[Point]struct{}{}

	for _, min := range grid.getMins() {
		basin := []Point{}
		queue := []Point{min}

		for len(queue) > 0 {
			l := len(queue) - 1
			point := queue[l]
			queue = queue[:l]

			if _, ok := visited[point]; ok {
				continue
			}

			if grid.get(point.x, point.y) >= 9 {
				continue
			}

			basin = append(basin, point)
			visited[point] = struct{}{}

			grid.forEachNeighbor(point.x, point.y, func(i uint, j uint) {
				queue = append(queue, Point{x: i, y: j})
			})
		}
		basins = append(basins, basin)
	}

	return basins
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

	mins := grid.getMins()
	risk := 0
	for _, point := range mins {
		risk += int(grid.get(point.x, point.y)) + 1
	}
	fmt.Printf("Part 1: %d\n", risk)

	basins := grid.getBasins()
	sort.Slice(basins, func(i int, j int) bool { return len(basins[i]) >= len(basins[j]) })

	prod := 1
	for _, basin := range basins[:3] {
		prod *= len(basin)
	}
	fmt.Printf("Part 2: %d\n", prod)
}
