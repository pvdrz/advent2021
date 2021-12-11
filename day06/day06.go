package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const MAX_DAYS = 9

type Fishes struct {
	counts [MAX_DAYS]uint
}

func (fishes *Fishes) Total() uint {
	total := uint(0)
	for _, count := range fishes.counts {
		total += count
	}
	return total
}

func (fishes *Fishes) Evolve() {
	for day, count := range fishes.counts {
		if day == 0 {
			fishes.counts[0] -= count
			fishes.counts[6] += count
			fishes.counts[MAX_DAYS-1] += count
		} else {
			fishes.counts[day] -= count
			fishes.counts[day-1] += count
		}
	}
}

func ParseInput(input string) (Fishes, error) {
	var fishes Fishes

	for _, chunk := range strings.Split(input, ",") {
		day, err := strconv.Atoi(chunk)
		if err != nil {
			return fishes, err
		}
		fishes.counts[day] += 1
	}

	return fishes, nil
}

func main() {
	bytes, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	input := strings.Trim(string(bytes), "\n")

	state, err := ParseInput(input)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 80; i += 1 {
		state.Evolve()
	}
	fmt.Printf("Part 1: %d\n", state.Total())

	for i := 80; i < 256; i += 1 {
		state.Evolve()
	}
	fmt.Printf("Part 2: %d\n", state.Total())
}
