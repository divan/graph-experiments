package layout

import "testing"

func TestInsert(t *testing.T) {
	o := NewOctree()
	p1 := &Point{1, 1, 1, 1, 10}
	o.Insert(p1)

	if o.root == nil {
		t.Fatalf("Expected root node to be non-nil")
	}

	center := o.root.Center()
	if center != p1 {
		t.Fatalf("Expected center to be %v, but got %v", p1, center)
	}

	p2 := &Point{2, 9, 9, 9, 10}
	o.Insert(p2)

	center = o.root.Center()
	expected := &Point{0, 5, 5, 5, 20}
	if *center != *expected {
		t.Fatalf("Expected center to be %v, but got %v", expected, center)
	}
}

func TestFindOctantIdx(t *testing.T) {
	var tests = []struct {
		name string
		p    *Point
		idx  int
	}{
		{
			name: "bottom back right",
			p:    &Point{X: 9, Y: 9, Z: 9},
			idx:  7,
		},
		{
			name: "top front left",
			p:    &Point{X: 1, Y: 1, Z: 1},
			idx:  0,
		},
		{
			name: "bottom front right",
			p:    &Point{X: 9, Y: 2, Z: 9},
			idx:  5,
		},
	}

	o := newLeaf(&Point{0, 5, 5, 5, 1})
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			idx := findOctantIdx(o, test.p)
			if idx != test.idx {
				t.Fatalf("Expected idx %d, but got %d", test.idx, idx)
			}
		})
	}
}
