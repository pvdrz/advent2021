package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Dot struct {
	x uint
	y uint
}

type Fold struct {
	intersect     uint
	is_horizontal bool
}

func parseFold(input string) (Fold, error) {
	input = strings.TrimSpace(input)

	chunks := strings.SplitN(input, "=", 2)
	if len(chunks) != 2 {
		return (Fold{}), fmt.Errorf("invalid fold format: \"%s\"", input)
	}

	fst := chunks[0]
	if len(fst) == 0 {
		return (Fold{}), fmt.Errorf("no characters before `=`")
	}

	is_horizontal := fst[len(fst)-1] == 'y'

	intersect, err := strconv.ParseUint(chunks[1], 10, 0)
	if err != nil {
		return (Fold{}), err
	}

	return Fold{
		intersect:     uint(intersect),
		is_horizontal: is_horizontal,
	}, nil
}

type Page struct {
	dots map[Dot]struct{}
}

func parsePage(input string) (Page, error) {
	page := Page{dots: map[Dot]struct{}{}}
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		chunks := strings.SplitN(line, ",", 2)
		if len(chunks) != 2 {
			return page, fmt.Errorf("invalid dot format: \"%s\"", line)
		}

		x, err := strconv.ParseUint(chunks[0], 10, 0)
		if err != nil {
			return page, err
		}

		y, err := strconv.ParseUint(chunks[1], 10, 0)
		if err != nil {
			return page, err
		}

		page.add(Dot{x: uint(x), y: uint(y)})
	}

	return page, nil
}

func (page *Page) len() uint {
	return uint(len(page.dots))
}

func (page *Page) add(dot Dot) {
	page.dots[dot] = struct{}{}
}

func (page *Page) remove(dot Dot) {
	delete(page.dots, dot)
}

func (page *Page) contains(dot Dot) bool {
	_, ok := page.dots[dot]
	return ok
}

func (page *Page) verticalFold(x0 uint) error {
	for dot := range page.dots {
		if dot.x > x0 {
			dx := dot.x - x0
			if x0 < dx {
				return fmt.Errorf("the line `x = %d` is too close to the left", x0)
			}
			page.remove(dot)
			page.add(Dot{x: x0 - dx, y: dot.y})
		}
	}

	return nil
}

func (page *Page) horizontalFold(y0 uint) error {
	for dot := range page.dots {
		if dot.y > y0 {
			dy := dot.y - y0
			if y0 < dy {
				return fmt.Errorf("the line `y = %d` is too close to the top", y0)
			}
			page.remove(dot)
			page.add(Dot{x: dot.x, y: y0 - dy})
		}
	}

	return nil
}

func (page *Page) fold(fold Fold) error {
	if fold.is_horizontal {
		return page.horizontalFold(fold.intersect)
	} else {
		return page.verticalFold(fold.intersect)
	}
}

func (page *Page) display() {
	max_x := uint(0)
	max_y := uint(0)
	for dot := range page.dots {
		if max_x < dot.x {
			max_x = dot.x
		}
		if max_y < dot.y {
			max_y = dot.y
		}
	}

	lines := make([]string, max_y + 1)
    for j := uint(0); j <= max_y; j += 1 {
		for i := uint(0); i <= max_x; i += 1 {
			if page.contains(Dot{x: i, y: j}) {
				lines[j] += "██"
			} else {
				lines[j] += "  "
			}
		}
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}

func parseInput(input string) (Page, []Fold, error) {
	input = strings.Trim(input, "\n")

	chunks := strings.SplitN(input, "\n\n", 2)
	if len(chunks) != 2 {
		return (Page{}), nil, fmt.Errorf("invalid input format")
	}

	page, err := parsePage(chunks[0])
	if err != nil {
		return page, nil, err
	}

	snd := strings.TrimSpace(chunks[1])
	snd = strings.Trim(snd, "\n")
	lines := strings.Split(snd, "\n")
	folds := make([]Fold, len(lines))
	for i, line := range lines {
		fold, err := parseFold(line)
		if err != nil {
			return page, folds, err
		}
		folds[i] = fold
	}

	return page, folds, nil
}

func main() {
	bytes, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	input := string(bytes)
	page, folds, err := parseInput(input)
	if err != nil {
		panic(err)
	}

	err = page.fold(folds[0])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 1: %d\n", page.len())

	for _, fold := range folds[1:] {
		err := page.fold(fold)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Part 2")
	page.display()
}
