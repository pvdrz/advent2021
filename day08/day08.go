package main

import (
	"fmt"
	"io/ioutil"
	"math/bits"
	"strings"
)

type Signal uint8

const (
	A Signal = 1 << iota
	B
	C
	D
	E
	F
	G
)

func (signal Signal) String() string {
	switch signal {
	case A:
		return "A"
	case B:
		return "B"
	case C:
		return "C"
	case D:
		return "D"
	case E:
		return "E"
	case F:
		return "F"
	case G:
		return "G"
	default:
		panic("unreachable")
	}
}

const mask uint8 = 1<<7 - 1

var set0, _ = ParseSignalSet("abcefg")
var set1, _ = ParseSignalSet("cf")
var set2, _ = ParseSignalSet("acdeg")
var set3, _ = ParseSignalSet("acdfg")
var set4, _ = ParseSignalSet("bcdf")
var set5, _ = ParseSignalSet("abdfg")
var set6, _ = ParseSignalSet("abdefg")
var set7, _ = ParseSignalSet("acf")
var set8, _ = ParseSignalSet("abcdefg")
var set9, _ = ParseSignalSet("abcdfg")

type SignalSet struct {
	bits uint8
}

func ParseSignalSet(input string) (SignalSet, error) {
	set := SignalSet{}
	for _, char := range input {
		switch char {
		case 'a':
			set = set.Insert(A)
		case 'b':
			set = set.Insert(B)
		case 'c':
			set = set.Insert(C)
		case 'd':
			set = set.Insert(D)
		case 'e':
			set = set.Insert(E)
		case 'f':
			set = set.Insert(F)
		case 'g':
			set = set.Insert(G)
		default:
			return set, fmt.Errorf("invalid character %q", char)
		}
	}

	return set, nil
}

func ParseSignalSets(input string) ([]SignalSet, error) {
	input = strings.TrimSpace(input)
	inputs := strings.Split(input, " ")
	sets := make([]SignalSet, len(inputs))

	for i, input := range inputs {
		set, err := ParseSignalSet(input)
		if err != nil {
			return sets, err
		}
		sets[i] = set
	}

	return sets, nil
}

func (set *SignalSet) Len() int {
	return bits.OnesCount8(set.bits)
}

func (set SignalSet) Insert(sig Signal) SignalSet {
	return SignalSet{bits: set.bits | uint8(sig)}
}

func (set *SignalSet) Contains(sig Signal) bool {
	return set.bits&uint8(sig) != 0
}

func (set SignalSet) Remove(sig Signal) SignalSet {
	return SignalSet{bits: set.bits & (^uint8(sig) & mask)}
}

func (set SignalSet) IsSuperSet(other SignalSet) bool {
	return (other.Compl().bits | set.bits) == mask
}

func (set SignalSet) Compl() SignalSet {
	return SignalSet{
		bits: ^set.bits & mask,
	}
}

func (set SignalSet) Union(other SignalSet) SignalSet {
	return SignalSet{
		bits: set.bits | other.bits,
	}
}

func (set SignalSet) Inter(other SignalSet) SignalSet {
	return SignalSet{
		bits: set.bits & other.bits,
	}
}

func (set SignalSet) Diff(other SignalSet) SignalSet {
	return set.Inter(other.Compl())
}

func (set SignalSet) SymDiff(other SignalSet) SignalSet {
	return SignalSet{
		bits: (set.bits ^ other.bits) & mask,
	}
}

func (set SignalSet) Signals() []Signal {
	signals := make([]Signal, 0)
	for _, signal := range []Signal{A, B, C, D, E, F, G} {
		if set.Contains(signal) {
			signals = append(signals, signal)
		}
	}
	return signals
}

func (set SignalSet) AssertSingleton() (Signal, error) {
	var signal Signal

	if set.Len() != 1 {
		return signal, fmt.Errorf("set %v is not a singleton", set)
	}

	signal = set.Signals()[0]
	return signal, nil
}

func (set SignalSet) ToDigit() (uint8, error) {
	switch set {
	case set0:
		return 0, nil
	case set1:
		return 1, nil
	case set2:
		return 2, nil
	case set3:
		return 3, nil
	case set4:
		return 4, nil
	case set5:
		return 5, nil
	case set6:
		return 6, nil
	case set7:
		return 7, nil
	case set8:
		return 8, nil
	case set9:
		return 9, nil
	default:
		return 255, fmt.Errorf("the set %v does not match any digit", set)
	}
}

func (set SignalSet) String() string {
	return fmt.Sprint(set.Signals())
}

func FindSignalSet(sets []SignalSet, f func(SignalSet) bool) (SignalSet, error) {
	for _, set := range sets {
		if f(set) {
			return set, nil
		}
	}
	return SignalSet{}, fmt.Errorf("not found")
}

func GetDecode(sets []SignalSet) func(SignalSet) (uint8, error) {
	one_in, err := FindSignalSet(sets, func(set SignalSet) bool { return set.Len() == 2 })
	if err != nil {
		panic("set for one is missing")
	}

	seven_in, err := FindSignalSet(sets, func(set SignalSet) bool { return set.Len() == 3 })
	if err != nil {
		panic("set for seven is missing")
	}

	a_in, err := seven_in.Diff(one_in).AssertSingleton()
	if err != nil {
		panic(err)
	}

	sixes := SignalSet{}
	for _, set := range sets {
		if set.Len() == 6 {
			sixes = sixes.Union(set.Compl())
		}
	}

	two_in, err := FindSignalSet(sets, func(set SignalSet) bool {
		return set.Len() == 5 && set.Contains(a_in) && set.IsSuperSet(sixes)
	})
	if err != nil {
		panic("set for two is missing")
	}

	g_in, err := two_in.Diff(sixes).Remove(a_in).AssertSingleton()
	if err != nil {
		panic(err)
	}

	three_in, err := FindSignalSet(sets, func(set SignalSet) bool {
		return set.Len() == 5 && set.Contains(g_in) && set.IsSuperSet(seven_in)
	})
	if err != nil {
		panic("set for three is missing")
	}

	d_in, err := three_in.Diff(seven_in).Remove(g_in).AssertSingleton()
	if err != nil {
		panic(err)
	}

	five_in, err := FindSignalSet(sets, func(set SignalSet) bool {
		return set.Len() == 5 && set != two_in && set != three_in
	})
	if err != nil {
		panic("set for five is missing")
	}

	b_and_f := five_in.Remove(a_in).Remove(g_in).Remove(d_in)

	f_in, err := one_in.Inter(b_and_f).AssertSingleton()
	if err != nil {
		panic(err)
	}

	b_in, err := b_and_f.Remove(f_in).AssertSingleton()
	if err != nil {
		panic(err)
	}

	c_in, err := one_in.Remove(f_in).AssertSingleton()
	if err != nil {
		panic(err)
	}

	eight_in, err := FindSignalSet(sets, func(set SignalSet) bool { return set.Len() == 7 })
	if err != nil {
		panic("set for eight is missing")
	}

	e_in, err := eight_in.Remove(a_in).Remove(b_in).Remove(c_in).Remove(d_in).Remove(f_in).Remove(g_in).AssertSingleton()
	if err != nil {
		panic(err)
	}

	mappings := map[Signal]Signal{
		a_in: A,
		b_in: B,
		c_in: C,
		d_in: D,
		e_in: E,
		f_in: F,
		g_in: G,
	}

	return func(set_in SignalSet) (uint8, error) {
		set_out := SignalSet{}
		for _, sig_in := range set_in.Signals() {
			sig_out := mappings[sig_in]
			set_out = set_out.Insert(sig_out)
		}
		return set_out.ToDigit()
	}
}

type Entry struct {
	clues  []SignalSet
	values []SignalSet
}

func ParseEntry(input string) (Entry, error) {
	input = strings.TrimSpace(input)

	chunks := strings.SplitN(input, " | ", 2)
	if len(chunks) != 2 {
		return (Entry{}), fmt.Errorf("invalid entry format")
	}

	clues, err := ParseSignalSets(chunks[0])
	if err != nil {
		return (Entry{}), err
	}

	values, err := ParseSignalSets(chunks[1])
	if err != nil {
		return (Entry{}), err
	}

	return Entry{clues: clues, values: values}, nil
}

func main() {
	bytes, err := ioutil.ReadFile("./input")
	if err != nil {
		panic(err)
	}

	input := strings.Trim(string(bytes), "\n")
	lines := strings.Split(input, "\n")

	count := 0
	sum := 0
	for _, line := range lines {
		entry, err := ParseEntry(line)
		if err != nil {
			panic(err)
		}

		decode := GetDecode(entry.clues)

		value := 0
		pows := [4]int{1000, 100, 10, 1}
		for i, set_in := range entry.values {
			digit, err := decode(set_in)
			if err != nil {
				panic(err)
			}
			if digit == 1 || digit == 4 || digit == 7 || digit == 8 {
				count += 1
			}
			value += int(digit) * pows[i]
		}
		sum += value
	}

	fmt.Printf("Part 1: %d\n", count)
	fmt.Printf("Part 2: %d\n", sum)
}
