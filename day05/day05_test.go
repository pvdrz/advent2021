package main

import (
	"testing"
)

func checkPoints(t *testing.T, input string, expected_points []Point) func(*testing.T) {
	return func(t *testing.T) {
		segment, err := ParseSegment(input)
		if err != nil {
			t.Fatalf("format is not OK: %s", err)
		}

		points := segment.Points()

		for i, expected_point := range expected_points {
			point := points[i]
			if expected_point != point {
				t.Errorf("point mismatch: expected %v, found %v", expected_point, point)
			}
		}
	}
}

func checkContains(t *testing.T, input string, expected_points []Point) func(*testing.T) {
	return func(t *testing.T) {
		segment, err := ParseSegment(input)
		if err != nil {
			t.Fatalf("format is not OK: %s", err)
		}

		for _, point := range expected_points {
			if !segment.Contains(&point) {
				t.Errorf("point %v is not in the segment", point)
			}
		}

		point := expected_points[0]
		point.x += 1
		point.y += 2

		if segment.Contains(&point) {
			t.Errorf("point %v should not be in the segment", point)
		}
	}
}

func checkIntersect(t *testing.T, input_a string, input_b string, expected_points func(*Segment, *Segment) []Point) func(*testing.T) {
	return func(t *testing.T) {
		segment_a, err := ParseSegment(input_a)
		if err != nil {
			t.Fatalf("format is not OK: %s", err)
		}

		segment_b, err := ParseSegment(input_b)
		if err != nil {
			t.Fatalf("format is not OK: %s", err)
		}

		points := segment_a.Intersection(&segment_b)

		for i, expected_point := range expected_points(&segment_a, &segment_b) {
			point := points[i]
			if expected_point != point {
				t.Errorf("point mismatch: expected %v, found %v", expected_point, point)
			}
		}
	}
}

func TestContains(t *testing.T) {
	t.Run("horizontal", checkContains(t, "0,9 -> 5,9", []Point{
		{x: 0, y: 9},
		{x: 1, y: 9},
		{x: 2, y: 9},
		{x: 3, y: 9},
		{x: 4, y: 9},
		{x: 5, y: 9},
	}))
	t.Run("vertical", checkContains(t, "1,2 -> 1,4", []Point{
		{x: 1, y: 2},
		{x: 1, y: 3},
		{x: 1, y: 4},
	}))
	t.Run("diagonal", checkContains(t, "2,1 -> 4,3", []Point{
		{x: 2, y: 1},
		{x: 3, y: 2},
		{x: 4, y: 3},
	}))
	t.Run("antidiagonal", checkContains(t, "1,5 -> 3,3", []Point{
		{x: 1, y: 5},
		{x: 2, y: 4},
		{x: 3, y: 3},
	}))
}

func TestPoints(t *testing.T) {
	t.Run("horizontal", checkPoints(t, "0,9 -> 5,9", []Point{
		{x: 0, y: 9},
		{x: 1, y: 9},
		{x: 2, y: 9},
		{x: 3, y: 9},
		{x: 4, y: 9},
		{x: 5, y: 9},
	}))
	t.Run("vertical", checkPoints(t, "1,2 -> 1,4", []Point{
		{x: 1, y: 2},
		{x: 1, y: 3},
		{x: 1, y: 4},
	}))
	t.Run("diagonal", checkPoints(t, "2,1 -> 4,3", []Point{
		{x: 2, y: 1},
		{x: 3, y: 2},
		{x: 4, y: 3},
	}))
	t.Run("antidiagonal", checkPoints(t, "1,5 -> 3,3", []Point{
		{x: 1, y: 5},
		{x: 2, y: 4},
		{x: 3, y: 3},
	}))
}

func TestIntersect(t *testing.T) {
	t.Run("same line", checkIntersect(t, "0,0 -> 4,6", "0,0 -> 4,6", func(a *Segment, _ *Segment) []Point { return a.Points() }))
	t.Run("contained b in a", checkIntersect(t, "0,9 -> 5,9", "0,9 -> 2,9", func(_ *Segment, b *Segment) []Point { return b.Points() }))
	t.Run("horizontal and vertical", checkIntersect(t, "0,9 -> 5,9", "3,0 -> 3,10", func(_ *Segment, b *Segment) []Point { return []Point{{x: 3, y: 9}} }))
	t.Run("horizontal and diagonal", checkIntersect(t, "0,9 -> 15,9", "9,9 -> 7,7", func(_ *Segment, b *Segment) []Point { return []Point{{x: 9, y: 9}} }))
	t.Run("singleton intersection", checkIntersect(t, "9,4 -> 3,4", "3,4 -> 1,4", func(a *Segment, b *Segment) []Point { return []Point{{x: 3, y: 4}} }))
}
