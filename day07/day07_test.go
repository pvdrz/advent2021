package main

import "testing"

type testCase struct {
	goal uint
	fuel uint
}

func TestEvolve(t *testing.T) {
	t.Run("get regular fuel", func(t *testing.T) {
		positions, err := ParseInput("16,1,2,0,4,2,7,1,2,14")
		if err != nil {
			t.Fatalf("invalid input: %s", err)
		}

		test_cases := []testCase{
			{goal: 1, fuel: 41},
			{goal: 3, fuel: 39},
			{goal: 10, fuel: 71},
			{goal: 2, fuel: 37},
		}

		for _, test_case := range test_cases {
			fuel := GetRegularFuel(positions, test_case.goal)
			if fuel != test_case.fuel {
				t.Errorf("invalid fuel for goal %d: expected %d, found %d", test_case.goal, test_case.fuel, fuel)
			}
		}
	})

	t.Run("regular fuel example", func(t *testing.T) {
		positions, err := ParseInput("16,1,2,0,4,2,7,1,2,14")
		if err != nil {
			t.Fatalf("invalid input: %s", err)
		}

		fuel := LeastFuel(positions, GetRegularFuel)

		if fuel != 37 {
			t.Fatalf("invalid least fuel: expected 37, found %d", fuel)
		}
	})

	t.Run("crab fuel example", func(t *testing.T) {
		positions, err := ParseInput("16,1,2,0,4,2,7,1,2,14")
		if err != nil {
			t.Fatalf("invalid input: %s", err)
		}

		fuel := LeastFuel(positions, GetCrabFuel)

		if fuel != 168 {
			t.Fatalf("invalid least fuel: expected 168, found %d", fuel)
		}
	})
}
