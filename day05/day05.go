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
	x0    int64
	y0    int64
	mx    int64
	my    int64
	t_max uint64
}

func (segment *Segment) Points() []Point {
	points_len := int64(segment.t_max) + 1
	points := make([]Point, points_len)

	for i := int64(0); i < points_len; i += 1 {
		points[i] = Point{
			x: segment.x0 + i*segment.mx,
			y: segment.y0 + i*segment.my,
		}
	}

	return points
}

func (segment *Segment) Intersection(other *Segment) []Point {
	intersection := make([]Point, 0)

	for _, point := range segment.Points() {
		if other.Contains(&point) {
			intersection = append(intersection, point)
		}
	}

	return intersection
}

func (segment *Segment) Contains(point *Point) bool {
	my_is_zero := segment.my == 0
	mx_is_zero := segment.mx == 0

	if my_is_zero {
		if mx_is_zero {
			// collapsed
			panic("unreachable")
		} else {
			// horizontal
			ty := point.y - segment.y0
			tx := (point.x - segment.x0) / segment.mx

			return (ty == 0) && (tx >= 0) && (tx <= int64(segment.t_max))
		}
	} else {
		if mx_is_zero {
			// vertical
			ty := (point.y - segment.y0) / segment.my
			tx := point.x - segment.x0

			return (tx == 0) && (ty >= 0) && (ty <= int64(segment.t_max))
		} else {
			// any diagonal
			tym := point.y - segment.y0
			txm := point.x - segment.x0

			if tym%segment.my != 0 && txm%segment.mx != 0 {
				return false
			}

			ty := tym / segment.my
			tx := txm / segment.mx

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

	if x, err := strconv.ParseInt(first_tuple[0], 10, 64); err != nil {
		return segment, err
	} else {
		x0 = x
	}
	if y, err := strconv.ParseInt(first_tuple[1], 10, 64); err != nil {
		return segment, err
	} else {
		y0 = y
	}
	if x, err := strconv.ParseInt(second_tuple[0], 10, 64); err != nil {
		return segment, err
	} else {
		x1 = x
	}
	if y, err := strconv.ParseInt(second_tuple[1], 10, 64); err != nil {
		return segment, err
	} else {
		y1 = y
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

	segment = Segment{x0: x0, y0: y0, mx: mx / t_max, my: my / t_max, t_max: uint64(t_max)}

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
		fmt.Println(err)
		return
	}

	segments, err := parseInput(string(input))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Part 2: %d\n", CountIntersects(segments))
}
