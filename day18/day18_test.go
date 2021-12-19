package main

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	inputs := []string{
		"[[[[[9,8],1],2],3],4]",
		"[[1,2],[[3,4],5]]",
		"[[6,[5,[4,[3,2]]]],1]",
	}
	for i, input := range inputs {
		t.Run(fmt.Sprintf("parse %d", i+1), func(t *testing.T) {
			baskets, err := parseInput(input)
			if err != nil {
				t.Fatalf("parsing failed: %s", err)
			}
			output := displayBaskets(baskets)
			if output != input {
				t.Errorf("invalid basket: expected `%s`, found `%s`", input, output)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		a, err := parseInput("[1,2]")
		if err != nil {
			t.Fatalf("parsing failed: %s", err)
		}
		b, err := parseInput("[[3,4],5]")
		if err != nil {
			t.Fatalf("parsing failed: %s", err)
		}
		output := displayBaskets(add(a, b))

		if output != "[[1,2],[[3,4],5]]" {
			t.Errorf("invalid basket: expected `[[1,2],[[3,4],5]]`, found `%s`", output)
		}
	})
}

func TestExplode(t *testing.T) {
	inputs := []string{
		"[[[[[9,8],1],2],3],4]",
		"[[6,[5,[4,[3,2]]]],1]",
	}
	outputs := []string{
		"[[[[0,9],2],3],4]",
		"[[6,[5,[7,0]]],3]",
	}
	for i, input := range inputs {
		t.Run(fmt.Sprintf("explode %d", i+1), func(t *testing.T) {
			baskets, err := parseInput(input)
			if err != nil {
				t.Fatalf("parsing failed: %s", err)
			}

			result, ok := explode(baskets)
			if !ok {
				t.Fatalf("`%s` did not explode", displayBaskets(result))
			}

			output := displayBaskets(result)
			if output != outputs[i] {
				t.Errorf("invalid basket: expected `%s`, found `%s`", outputs[i], output)
			}
		})
	}
}

func TestMagnitude(t *testing.T) {
	inputs := []string{
		"[[1,2],[[3,4],5]]",
		"[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
		"[[[[1,1],[2,2]],[3,3]],[4,4]]",
		"[[[[3,0],[5,3]],[4,4]],[5,5]]",
		"[[[[5,0],[7,4]],[5,5]],[6,6]]",
		"[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]",
	}
	outputs := []uint{
		143,
		1384,
		445,
		791,
		1137,
		3488,
	}
	for i, input := range inputs {
		t.Run(fmt.Sprintf("magnitude %d", i+1), func(t *testing.T) {
			baskets, err := parseInput(input)
			if err != nil {
				t.Fatalf("parsing failed: %s", err)
			}

			output := magnitude(baskets)
			if output != outputs[i] {
				t.Errorf("invalid magnitude: expected `%d`, found `%d`", outputs[i], output)
			}
		})
	}
}
