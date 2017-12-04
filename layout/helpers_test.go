package layout

import (
	"math"
	"testing"
)

func TestDistance(t *testing.T) {
	var tests = []struct {
		from, to *Node
		expected float64
	}{
		{
			from:     &Node{"", &Point{X: 0, Y: 0, Z: 0}},
			to:       &Node{"", &Point{X: 10, Y: 10, Z: 10}},
			expected: math.Sqrt(300),
		},
		{
			from:     &Node{"", &Point{X: 2, Y: 3, Z: 1}},
			to:       &Node{"", &Point{X: 8, Y: -5, Z: 0}},
			expected: math.Sqrt(101),
		},
	}
	for _, test := range tests {
		got := distance(test.from.Point, test.to.Point)
		if got != test.expected {
			t.Fatalf("Expected %.3f, but got %.3f", test.expected, got)
		}
	}
}
