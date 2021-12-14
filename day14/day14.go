package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

func parseInput(input string) (map[string]uint, map[string]string, error) {
	input = strings.TrimSpace(input)
	input = strings.Trim(input, "\n")

	chunks := strings.SplitN(input, "\n\n", 2)
	if len(chunks) != 2 {
		return nil, nil, fmt.Errorf("invalid input format")
	}

	fst := strings.Trim(chunks[0], "\n")
	fst = strings.TrimSpace(fst)

	template := map[string]uint{}
	for i := 0; i < len(fst)-1; i += 1 {
		input := fst[i : i+2]
		template[input] += 1
	}
	// include the last character plus ` ` to avoid counting it twice
	last_pair := fst[len(fst)-1:] + " "
	template[last_pair] = 1

	rules := map[string]string{}
	for _, line := range strings.Split(chunks[1], "\n") {
		line = strings.TrimSpace(line)
		chunks := strings.SplitN(line, " -> ", 2)
		if len(chunks) != 2 {
			return nil, nil, fmt.Errorf("invalid rule format")
		}
		input := chunks[0]
		if len(input) != 2 {
			return nil, nil, fmt.Errorf("invalid input format")
		}

		output := chunks[1]
		if len(output) != 1 {
			return nil, nil, fmt.Errorf("invalid output format")
		}

		rules[input] = output
	}

	return template, rules, nil
}

func step(template map[string]uint, rules map[string]string) map[string]uint {
	result := map[string]uint{}
	for input, count := range template {
		output, ok := rules[input]
		if !ok {
			result[input] = count
		} else {
			output1 := input[0:1] + output
			output2 := output + input[1:2]
			result[output1] += count
			result[output2] += count
		}
	}

	return result
}

func run(template map[string]uint, rules map[string]string, n_steps uint) map[string]uint {
	for i := uint(0); i < n_steps; i += 1 {
		template = step(template, rules)
	}

	return template
}

func getMaxMinCount(template map[string]uint) (uint, uint) {
	counts := map[rune]uint{}

	for input, count := range template {
		// just count the first character of each input pair to avoid repetitions
		for _, char := range input {
			counts[char] += count
			break
		}
	}

	min := uint(math.MaxUint)
	max := uint(0)

	for _, count := range counts {
		if min > count {
			min = count
		}
		if max < count {
			max = count
		}
	}

	return max, min
}

func main() {
	bytes, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	input := string(bytes)
	template, rules, err := parseInput(input)
	if err != nil {
		panic(err)
	}

	output := run(template, rules, 10)
	max, min := getMaxMinCount(output)

	fmt.Printf("Part 1: %d\n", max-min)

	output = run(output, rules, 30)
	max, min = getMaxMinCount(output)

	fmt.Printf("Part 2: %d\n", max-min)
}
