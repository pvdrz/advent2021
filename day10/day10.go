package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

func pop(stack []rune) (rune, []rune) {
	new_len := len(stack) - 1
	return stack[new_len], stack[:new_len]
}

func checkLine(line string) (int, int) {
	var head rune
	stack := make([]rune, 0)

	for _, char := range line {
		switch char {
		case '(':
			stack = append(stack, char)
		case '[':
			stack = append(stack, char)
		case '{':
			stack = append(stack, char)
		case '<':
			stack = append(stack, char)
		case ')':
			head, stack = pop(stack)
			if head != '(' {
				return 3, 0
			}
		case ']':
			head, stack = pop(stack)
			if head != '[' {
				return 57, 0
			}
		case '}':
			head, stack = pop(stack)
			if head != '{' {
				return 1197, 0
			}
		case '>':
			head, stack = pop(stack)
			if head != '<' {
				return 25137, 0
			}
		default:
			panic("invalid character")
		}
	}

	score := 0
	for i := len(stack) - 1; i >= 0; i -= 1 {
		switch stack[i] {
		case '(':
			score = 5*score + 1
		case '[':
			score = 5*score + 2
		case '{':
			score = 5*score + 3
		case '<':
			score = 5*score + 4
		default:
			panic("invalid character in stack")
		}
	}

	return 0, score
}

func main() {
	bytes, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	input := strings.Trim(string(bytes), "\n")
	lines := strings.Split(input, "\n")

	check_score := 0
    compl_scores := make([]int, 0)
	for _, line := range lines {
		check, compl := checkLine(line)
		check_score += check
        if compl != 0 {
            compl_scores = append(compl_scores, compl)
        }
	}

    sort.Ints(compl_scores)

	fmt.Printf("Part 1: %d\n", check_score)
	fmt.Printf("Part 2: %d\n", compl_scores[len(compl_scores)/2])
}
