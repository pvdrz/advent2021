package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const EXAMPLE = `0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2`

type Point struct {
	x int64
	y int64
}

type Segment struct {
	start Point
	slope Point
	t_max int64
}

func (segment *Segment) isDiagonal() bool {
	return segment.slope.x != 0 && segment.slope.y != 0
}

func (segment *Segment) forEachPoint(f func(int64, Point)) {
	for i := int64(0); i <= segment.t_max; i += 1 {
		f(i, Point{
			x: segment.start.x + i*segment.slope.x,
			y: segment.start.y + i*segment.slope.y,
		})
	}
}

func (segment *Segment) Points() []Point {
	points := make([]Point, segment.t_max+1)

	segment.forEachPoint(func(i int64, point Point) {
		points[i] = point
	})

	return points
}

func (segment *Segment) Intersection(other *Segment) []Point {
	if segment.t_max > other.t_max {
		return other.Intersection(segment)
	}

	intersection := make([]Point, 0)

	segment.forEachPoint(func(_ int64, point Point) {
		if other.Contains(&point) {
			intersection = append(intersection, point)
		}
	})

	return intersection
}

func (segment *Segment) Contains(point *Point) bool {
	my_is_zero := segment.slope.y == 0
	mx_is_zero := segment.slope.x == 0

	if my_is_zero {
		if mx_is_zero {
			// collapsed
			panic("unreachable")
		} else {
			// horizontal
			ty := point.y - segment.start.y
			tx := (point.x - segment.start.x) / segment.slope.x

			return (ty == 0) && (tx >= 0) && (tx <= int64(segment.t_max))
		}
	} else {
		if mx_is_zero {
			// vertical
			ty := (point.y - segment.start.y) / segment.slope.y
			tx := point.x - segment.start.x

			return (tx == 0) && (ty >= 0) && (ty <= int64(segment.t_max))
		} else {
			// any diagonal
			tym := point.y - segment.start.y
			txm := point.x - segment.start.x

			if tym%segment.slope.y != 0 && txm%segment.slope.x != 0 {
				return false
			}

			ty := tym / segment.slope.y
			tx := txm / segment.slope.x

			return (tx == ty) && (tx >= 0) && (tx <= int64(segment.t_max))
		}
	}
}

func ParseSegment(input string) (Segment, error) {
	segment := Segment{}

	tuples := strings.SplitN(input, " -> ", 2)

	first_tuple := strings.SplitN(tuples[0], ",", 2)
	second_tuple := strings.SplitN(tuples[1], ",", 2)

	var x0 int64
	var y0 int64
	var x1 int64
	var y1 int64
	var err error

	if x0, err = strconv.ParseInt(first_tuple[0], 10, 64); err != nil {
		return segment, err
	}
	if y0, err = strconv.ParseInt(first_tuple[1], 10, 64); err != nil {
		return segment, err
	}
	if x1, err = strconv.ParseInt(second_tuple[0], 10, 64); err != nil {
		return segment, err
	}
	if y1, err = strconv.ParseInt(second_tuple[1], 10, 64); err != nil {
		return segment, err
	}

	mx := x1 - x0
	my := y1 - y0

	var t_max int64
	if my == 0 {
		if mx == 0 {
			return segment, errors.New("segment is collapsed")
		} else if mx < 0 {
			t_max = -mx
		} else {
			t_max = mx
		}
	} else if my < 0 {
		t_max = -my
	} else {
		t_max = my
	}

	segment = Segment{start: Point{x: x0, y: y0}, slope: Point{x: mx / t_max, y: my / t_max}, t_max: t_max}

	return segment, nil
}

func parseInput(input string) ([]Segment, error) {
	input = strings.TrimSpace(input)
	lines := strings.Split(input, "\n")
	segments := make([]Segment, len(lines))

	for i, line := range lines {
		if segment, err := ParseSegment(line); err != nil {
			return segments, err
		} else {
			segments[i] = segment
		}
	}

	return segments, nil
}

func CountIntersects(segments []Segment) int {
	n := len(segments)
	points := make(map[Point]struct{})

	for i, segment_a := range segments {
		for j := i + 1; j < n; j += 1 {
			segment_b := segments[j]
			for _, point := range segment_a.Intersection(&segment_b) {
				points[point] = struct{}{}
			}
		}
	}

	return len(points)
}

func main() {
	input, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	segments, err := parseInput(string(input))
	if err != nil {
		panic(err)
	}

    non_diag := make([]Segment, 0)
    for _, segment := range segments {
        if !segment.isDiagonal() {
            non_diag = append(non_diag, segment)
        }
    }

	fmt.Printf("Part 1: %d\n", CountIntersects(non_diag))
	fmt.Printf("Part 2: %d\n", CountIntersects(segments))
}
