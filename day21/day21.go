package main

import "fmt"

type Die struct {
	val   int
	rolls int
}

func newDie() Die {
	return Die{val: 0, rolls: 0}
}

func (die *Die) roll() int {
	die.val += 1
	die.val %= 100

	die.rolls += 1

	return die.val
}

func (dice *Die) totalRolls() int {
	return dice.rolls
}

func main() {
	p1 := 6
	p2 := 8

	s1 := 0
	s2 := 0

	die := newDie()
    losing_score := 0

	for i := 0; true; i += 1 {
		fmt.Printf("Turn %d\n", i)

		roll := die.roll() + die.roll() + die.roll()

		p1 = 1 + (p1+roll-1)%10
		s1 += p1
		fmt.Printf("Player 1 rolls %d and moves to space %d for a total score of %d.\n", roll, p1, s1)

		if s1 >= 1000 {
			losing_score = s2
            break
		}

		roll = die.roll() + die.roll() + die.roll()

		p2 = 1 + (p2+roll-1)%10
		s2 += p2
		fmt.Printf("Player 2 rolls %d and moves to space %d for a total score of %d.\n", roll, p2, s2)

		if s2 >= 1000 {
			losing_score = s1
			break
		}
	}
    fmt.Printf("The die rolled %d times.\n", die.totalRolls())
    fmt.Printf("Part 1: %d\n", losing_score * die.totalRolls())
}
