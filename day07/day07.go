package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func ParseInput(input string) ([]uint, error) {
	input = strings.Trim(input, "\n")
	input = strings.TrimSpace(input)
	chunks := strings.Split(input, ",")

	positions := make([]uint, len(chunks))

	for i, chunk := range chunks {
		position, err := strconv.ParseUint(chunk, 10, 0)
		if err != nil {
			return positions, err
		}
		positions[i] = uint(position)
	}

	return positions, nil
}

func GetRegularFuel(positions []uint, goal uint) uint {
	fuel := uint(0)
	for _, position := range positions {
		if position >= goal {
			fuel += position - goal
		} else {
			fuel += goal - position
		}
	}
	return fuel
}

func GetCrabFuel(positions []uint, goal uint) uint {
	fuel := uint(0)
	for _, position := range positions {
		var dist uint
		if position >= goal {
			dist = position - goal
		} else {
			dist = goal - position
		}
		fuel += dist * (dist + 1) / 2
	}
	return fuel
}

func LeastFuel(positions []uint, fuel_fn func([]uint, uint) uint) uint {
	min_pos := uint(math.MaxUint)
	max_pos := uint(0)

	for _, pos := range positions {
		if min_pos > pos {
			min_pos = pos
		}

		if max_pos < pos {
			max_pos = pos
		}
	}

	min_fuel := uint(math.MaxUint)
	for pos := min_pos; pos <= max_pos; pos += 1 {
		fuel := fuel_fn(positions, pos)
		if fuel < min_fuel {
			min_fuel = fuel
		}
	}

	return min_fuel
}

func main() {
	bytes, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	input := strings.Trim(string(bytes), "\n")

	positions, err := ParseInput(input)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Part 1: %d\n", LeastFuel(positions, GetRegularFuel))
	fmt.Printf("Part 2: %d\n", LeastFuel(positions, GetCrabFuel))
}
