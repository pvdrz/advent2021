package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

type Basket struct {
	value    uint
	is_right bool
	depth    uint
}

func displayBaskets(baskets []Basket) string {
	buf := ""
	depth := uint(0)

	for _, basket := range baskets {
		for depth > basket.depth {
			buf += "]"
			depth -= 1
		}
		if basket.is_right {
			buf += ","
			buf += fmt.Sprint(basket.value)
			buf += "]"
			depth -= 1
		} else {
			if depth == basket.depth {
				buf += "]"
				depth -= 1
			}
			buf += ","
			for depth < basket.depth {
				buf += "["
				depth += 1
			}
			buf += fmt.Sprint(basket.value)
		}
	}

	for depth > 0 {
		buf += "]"
		depth -= 1
	}

	return buf[1:]
}

func (basket Basket) String() string {
	char := "i"
	if basket.is_right {
		char = "d"
	}

	return fmt.Sprintf("(%d%s%d)", basket.value, char, basket.depth)
}

func parseInput(input string) ([]Basket, error) {
	input = strings.TrimSpace(input)

	depth := uint(0)
	is_right := false

	baskets := []Basket{}
	for _, char := range input {
		switch char {
		case '[':
			depth += 1
			is_right = false
		case ']':
			if depth == 0 {
				return baskets, fmt.Errorf("parentheses mismatch")
			}
			depth -= 1
		case ',':
			is_right = true
		default:
			if char >= '0' && char <= '9' {
				basket := Basket{
					value:    uint(char - '0'),
					depth:    depth,
					is_right: is_right,
				}
				baskets = append(baskets, basket)
			} else {
				return baskets, fmt.Errorf("invalid character")
			}
		}
	}

	return baskets, nil
}

func add(a []Basket, b []Basket) []Basket {
	res := make([]Basket, len(a)+len(b))

	for i, basket := range a {
		basket.depth += 1
		res[i] = basket
	}

	for i, basket := range b {
		basket.depth += 1
		res[i+len(a)] = basket
	}

	return res
}

func explode(baskets []Basket) ([]Basket, bool) {
	for i, left := range baskets {
		if left.depth > 4 && !left.is_right {
			right := baskets[i+1]
			if right.is_right && right.depth == left.depth {
				new_basket := Basket{
					value:    0,
					is_right: false,
					depth:    right.depth - 1,
				}

				if i-1 >= 0 {
					left_neighbor := baskets[i-1]
					left_neighbor.value += left.value
					if left_neighbor.depth == new_basket.depth && !left_neighbor.is_right {
						new_basket.is_right = true
					}
					baskets[i-1] = left_neighbor
				}

				if i+2 < len(baskets) {
					right_neighbor := baskets[i+2]
					right_neighbor.value += right.value
					if right_neighbor.depth == new_basket.depth && right_neighbor.is_right {
						if new_basket.is_right {
							panic("here be dragons")
						}
						new_basket.is_right = false
					}
					baskets[i+2] = right_neighbor
				}

				baskets = append(baskets[:i], baskets[i+1:]...)
				baskets[i] = new_basket
				return baskets, true
			}
		}
	}

	return baskets, false
}

func split(baskets []Basket) ([]Basket, bool) {
	for i, basket := range baskets {
		if basket.value >= 10 {
			value := float64(basket.value) / 2.0
			floor := uint(math.Floor(value))
			ceil := uint(math.Ceil(value))

			left := Basket{
				value:    floor,
				is_right: false,
				depth:    basket.depth + 1,
			}

			right := Basket{
				value:    ceil,
				is_right: true,
				depth:    basket.depth + 1,
			}

			baskets = append(baskets[:i+1], baskets[i:]...)
			baskets[i] = left
			baskets[i+1] = right
			return baskets, true
		}
	}

	return baskets, false
}

func reduce(baskets []Basket) []Basket {
	var did_explode bool
	var did_split bool

	for {
		baskets, did_explode = explode(baskets)
		if did_explode {
			continue
		}
		baskets, did_split = split(baskets)
		if !did_split {
			break
		}
	}

	return baskets
}
func magnitude(baskets []Basket) uint {
	for len(baskets) > 1 {
		for i, left := range baskets {
			if !left.is_right && i+1 < len(baskets) {
				right := baskets[i+1]
				if right.is_right && right.depth == left.depth {
					new_basket := Basket{
						value:    3*left.value + 2*right.value,
						is_right: false,
						depth:    right.depth - 1,
					}

					if i-1 >= 0 {
						left_neighbor := baskets[i-1]
						if left_neighbor.depth == new_basket.depth && !left_neighbor.is_right {
							new_basket.is_right = true
						}
					}

					if i+2 < len(baskets) {
						right_neighbor := baskets[i+2]
						if right_neighbor.depth == new_basket.depth && right_neighbor.is_right {
							if new_basket.is_right {
								panic("here be dragons")
							}
							new_basket.is_right = false
						}
					}

					baskets = append(baskets[:i], baskets[i+1:]...)
					baskets[i] = new_basket
					break
				}
			}
		}
	}

	return baskets[0].value
}

func main() {
	bytes, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	input := strings.Trim(string(bytes), "\n")

	lines := [][]Basket{}
	for _, line := range strings.Split(input, "\n") {
		baskets, err := parseInput(line)
		if err != nil {
			panic(err)
		}
		lines = append(lines, baskets)
	}

	if len(lines) == 0 {
		panic("input is empty")
	}

	acc := lines[0]
	for _, baskets := range lines[1:] {
		acc = add(acc, baskets)
		acc = reduce(acc)
	}

	fmt.Printf("Part 1: %d\n", magnitude(acc))

    max := uint(0)
	for i, a := range lines {
		for j := i + 1; j < len(lines); j += 1 {
			b := lines[j]
			sum1 := magnitude(reduce(add(a, b)))
			sum2 := magnitude(reduce(add(b, a)))

			for _, sum := range []uint{sum1, sum2} {
                if sum > max {
                    max = sum
                }
			}
		}
	}

	fmt.Printf("Part 2: %d\n", max)
}
