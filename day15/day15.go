package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type Point struct {
	x uint
	y uint
}

func (point *Point) equals(other *Point) bool {
	return point.x == other.x && point.y == other.y
}

type Queued struct {
	point    Point
	distance uint
}

func insert(queue []Queued, item Queued) []Queued {
	// we always search in the interval [min, max)
	min := uint(0)
	max := uint(len(queue))

	// this will hold the index of the element we will be checking in the
	// current iteration
	i := uint(0)
	// iterate while the interval is not empty
	for min < max {
		// update the index to the half of the interval
		i = (max + min) / 2
		// get the element to be checked
		current := queue[i]
		if current.distance < item.distance {
			// if the current element has a distance smaller than `item`, we
			// move to the interval [min, i)
			max = i
		} else if current.distance > item.distance {
			// if the current element has a distance larger than `item`, we
			// move to the interval [i + 1, max)
			min = i + 1
		} else {
			// if both elements have the same distance, we stop iterating
			// because it means we can insert `item` in the current index
			// without disordering the queue.
			break
		}
	}

	queue = append(queue[:i+1], queue[i:]...)
	queue[i] = item

	return queue
}

func pop(queue []Queued) ([]Queued, Queued) {
	n := len(queue) - 1
	return queue[:n], queue[n]
}

type Grid struct {
	rows   [][]uint
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

	rows := make([][]uint, 0)
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) != n_cols {
			return (Grid{}), fmt.Errorf("line %d has an invalid length: expected %d, found %d", i, n_cols, len(line))
		}

		row := make([]uint, n_cols)
		for i, char := range line {
			digit, err := strconv.ParseUint(string(char), 10, 0)
			if err != nil {
				return (Grid{}), err
			}
			row[i] = uint(digit)
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

func (grid *Grid) contains(i, j, tiles uint) bool {
	return i < tiles*grid.nCols() && j < tiles*grid.nRows()
}

func (grid *Grid) get(i, j, tiles uint) uint {
	value := grid.rows[j%grid.nRows()][i%grid.nCols()]

	extra := i/grid.nCols() + j/grid.nRows()

	for a := uint(0); a < extra; a += 1 {
		value += 1
		if value > 9 {
			value = 1
		}
	}

	return value
}

func (grid *Grid) set(i, j, val uint) {
	grid.rows[j][i] = val
}

func (grid *Grid) getNeighbors(i, j, tiles uint) []Point {
	neighbors := make([]Point, 0)

	up := j > 0 && grid.contains(i, j-1, tiles)
	left := i > 0 && grid.contains(i-1, j, tiles)
	down := j < math.MaxUint && grid.contains(i, j+1, tiles)
	right := i < math.MaxUint && grid.contains(i+1, j, tiles)
	if up {
		neighbors = append(neighbors, Point{x: i, y: j - 1})
	}
	if left {
		neighbors = append(neighbors, Point{x: i - 1, y: j})
	}
	if down {
		neighbors = append(neighbors, Point{x: i, y: j + 1})
	}
	if right {
		neighbors = append(neighbors, Point{x: i + 1, y: j})
	}

	return neighbors
}

func (grid *Grid) findPath(tiles uint) (uint, error) {
	start := Point{x: 0, y: 0}
	end := Point{x: tiles*grid.nCols() - 1, y: tiles*grid.nRows() - 1}

	dist := map[Point]uint{start: 0}

	queue := []Queued{{point: start, distance: dist[start]}}

	var curr Queued
	for len(queue) > 0 {
		queue, curr = pop(queue)
		if curr.point.equals(&end) {
			break
		}

		for _, neighbor := range grid.getNeighbors(curr.point.x, curr.point.y, tiles) {
			alt := curr.distance + grid.get(neighbor.x, neighbor.y, tiles)
			dst, ok := dist[neighbor]

			if !ok || alt < dst {
				dist[neighbor] = alt
				queue = insert(queue, Queued{point: neighbor, distance: alt})
			}
		}
	}

	if total, ok := dist[end]; ok {
		return total, nil
	} else {
		return total, fmt.Errorf("path not found")
	}
}

func main() {
	bytes, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	input := string(bytes)
	grid, err := parseGrid(input)
	if err != nil {
		panic(err)
	}

	path_len, err := grid.findPath(1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Part 1: %d\n", path_len)

	path_len, err = grid.findPath(5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Part 2: %d\n", path_len)
}
