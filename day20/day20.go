package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Point struct {
	x int
	y int
}

type Map struct {
	rules          [512]bool
	pixels         map[Point]struct{}
	mins           Point
	maxs           Point
	horizon_is_lit bool
}

func (m Map) String() string {
	buf := "\n"
	// print the x coords
	buf += "  "
	for x := m.mins.x - 3; x <= m.maxs.x+3; x += 1 {
		if x < 0 {
			buf += fmt.Sprint(-x % 10)
		} else {
			buf += fmt.Sprint(x % 10)
		}
	}
	for y := m.mins.y - 3; y <= m.maxs.y+3; y += 1 {
		buf += "\n"
		if y < 0 {
			buf += fmt.Sprint(-y % 10)
		} else {
			buf += fmt.Sprint(y % 10)
		}

		buf += " "
		for x := m.mins.x - 3; x <= m.maxs.x+3; x += 1 {
			if m.contains(Point{x: x, y: y}) {
				buf += "#"
			} else {
				buf += "."
			}
		}
	}
	return buf
}

func (m *Map) len() int {
	return len(m.pixels)
}

func (m *Map) insert(point Point) {
	if m.mins.x > point.x {
		m.mins.x = point.x
	}
	if m.mins.y > point.y {
		m.mins.y = point.y
	}
	if m.maxs.x < point.x {
		m.maxs.x = point.x
	}
	if m.maxs.y < point.y {
		m.maxs.y = point.y
	}
	m.pixels[point] = struct{}{}

}

func (m *Map) contains(point Point) bool {
	_, ok := m.pixels[point]
	return ok || (m.horizon_is_lit && m.inHorizon(point))
}

func (m *Map) inHorizon(point Point) bool {
	return point.x <= m.mins.x-2 || point.x >= m.maxs.x+2 || point.y <= m.mins.y-2 || point.y >= m.maxs.y+2
}

func (m *Map) neighbors(point Point) int {
	value := 0
	for y := point.y - 1; y <= point.y+1; y += 1 {
		for x := point.x - 1; x <= point.x+1; x += 1 {
			value *= 2
			if m.contains(Point{x: x, y: y}) {
				value += 1
			}
		}
	}
	return value
}

func (m *Map) evolve() {
	new_m := Map{
		rules:          m.rules,
		pixels:         map[Point]struct{}{},
		mins:           Point{x: 0, y: 0},
		maxs:           Point{x: 0, y: 0},
		horizon_is_lit: m.rules[0],
	}

	if m.horizon_is_lit {
		new_m.horizon_is_lit = m.rules[511]
	}

	for y := m.mins.y - 2; y <= m.maxs.y+2; y += 1 {
		for x := m.mins.x - 2; x <= m.maxs.x+2; x += 1 {
			point := Point{x: x, y: y}
			code := m.neighbors(point)
			if m.rules[code] {
				new_m.insert(point)
			}
		}
	}

	new_m.mins.x += 1
	new_m.mins.y += 1
	new_m.maxs.x -= 1
	new_m.maxs.y -= 1

	*m = new_m
}

func parseInput(input string) (Map, error) {
	input = strings.Trim(input, "\n")

	m := Map{
		pixels:         map[Point]struct{}{},
		mins:           Point{x: 0, y: 0},
		maxs:           Point{x: 0, y: 0},
		horizon_is_lit: false,
	}

	chunks := strings.SplitN(input, "\n\n", 2)

	if len(chunks) != 2 {
		return m, fmt.Errorf("invalid input format")
	}

	head := strings.TrimSpace(chunks[0])
	if len(head) != 512 {
		return m, fmt.Errorf("invalid length for first line: expected 512, found: %d", len(head))
	}

	for i, char := range head {
		m.rules[i] = char == '#'
	}

	for y, line := range strings.Split(chunks[1], "\n") {
		line = strings.TrimSpace(line)
		for x, char := range line {
			if char == '#' {
				m.insert(Point{x: x, y: y})
			}
		}
	}

	return m, nil
}

func main() {
	// bytes, err := ioutil.ReadFile("./divi")
	bytes, err := ioutil.ReadFile("./input")
	// bytes, err := ioutil.ReadFile("./test")
	// bytes, err := ioutil.ReadFile("./test2")
	if err != nil {
		panic(err)
	}

	input := string(bytes)
	m, err := parseInput(input)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 2; i += 1 {
		m.evolve()
	}
	fmt.Printf("Part 1: %d\n", m.len())

	for i := 2; i < 50; i += 1 {
		m.evolve()
	}
	fmt.Printf("Part 2: %d\n", m.len())
}
