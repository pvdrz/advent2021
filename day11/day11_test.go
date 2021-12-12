package main

import (
	"testing"
)

const small_grid = `
    11111
    19991
    19191
    19991
    11111
`
const small_expected = `
    34543
    40004
    50005
    40004
    34543
`
const large_grid = `
    6594254334
    3856965822
    6375667284
    7252447257
    7468496589
    5278635756
    3287952832
    7993992245
    5957959665
    6394862637
`
const large_expected = `
    8807476555
    5089087054
    8597889608
    8485769600
    8700908800
    6600088989
    6800005943
    0000007456
    9000000876
    8700006848
`
const count_input = `
    5483143223
    2745854711
    5264556173
    6141336146
    6357385478
    4167524645
    2176841721
    6882881134
    4846848554
    5283751526
`

func checkEvolve(input_grid, input_expected string) func(t *testing.T) {
	return func(t *testing.T) {
		grid, err := parseGrid(input_grid)
		if err != nil {
			t.Fatalf("input is invalid")
		}
		expected, err := parseGrid(input_expected)
		if err != nil {
			t.Fatalf("input is invalid")
		}

		grid.evolve()

		if !grid.equals(&expected) {
			t.Fatalf("evolve mismatch")
		}
	}
}

func TestEvolve(t *testing.T) {
	t.Run("small evolve", checkEvolve(small_grid, small_expected))
	t.Run("large evolve", checkEvolve(large_grid, large_expected))
}

func TestCount(t *testing.T) {
	grid, err := parseGrid(count_input)
	if err != nil {
		t.Fatalf("input is invalid")
	}
    if grid.run(10) != 204 {
        t.Fatalf("invalid flash count for 10 steps")
    }
    if grid.run(90) + 204 != 1656 {
        t.Fatalf("invalid flash count for 100 steps")
    }
}
