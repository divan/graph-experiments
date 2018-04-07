package layout

import "testing"

func TestInsert(t *testing.T) {
	o := NewOctree(NewGravityForce(-10, BarneHutMethod))
	p1 := &Point{Idx: 1, X: 1, Y: 1, Z: 1, Mass: 10}
	o.Insert(p1)

	if o.root == nil {
		t.Fatalf("Expected root node to be non-nil")
	}

	center := o.root.Center()
	if center != p1 {
		t.Fatalf("Expected center to be %v, but got %v", p1, center)
	}

	p2 := &Point{Idx: 2, X: 9, Y: 9, Z: 9, Mass: 10}
	o.Insert(p2)

	center = o.root.Center()
	expected := &Point{Idx: 0, X: 5, Y: 5, Z: 5, Mass: 20}
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

	o := newLeaf(&Point{X: 5, Y: 5, Z: 5, Mass: 1})
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			idx := findOctantIdx(o, test.p)
			if idx != test.idx {
				t.Fatalf("Expected idx %d, but got %d", test.idx, idx)
			}
		})
	}
}

func TestLeafInsert(t *testing.T) {
	p1 := &Point{X: 1, Y: 1, Z: 1, Mass: 1}
	p2 := &Point{Idx: 1, X: -1, Y: -1, Z: -1, Mass: 1}
	l := newLeaf(p1)
	center := l.Center()
	if center != p1 {
		t.Fatalf("center != p1")
	}
	node := l.Insert(p2)
	center = node.Center()
	expected := &Point{Mass: 2} // zero coords
	if *center != *expected {
		t.Fatalf("Expected %v, but got %v", expected, center)
	}
}

func TestBugCase1(t *testing.T) {
	o := NewOctree(NewGravityForce(-10, BarneHutMethod))
	points := []*Point{
		&Point{Idx: 0, X: -2, Y: 4, Z: 1, Mass: 2},
		&Point{Idx: 1, X: -6, Y: 4, Z: -1, Mass: 2},
		&Point{Idx: 2, X: -1, Y: -13, Z: 3, Mass: 2},
		&Point{Idx: 3, X: 14, Y: 14, Z: 5, Mass: 2},
		&Point{Idx: 4, X: -19, Y: -5, Z: 9, Mass: 2},
	}
	for i := 0; i < 5; i++ {
		o.Insert(points[i])
	}
	for i := 0; i < 5; i++ {
		leaf, err := findLeaf(o.root, i)
		if err != nil {
			t.Fatalf("Expected err to be non nil, got %v", err)
		}
		if leaf.point.Idx != i {
			t.Fatalf("Expected point index to be %d, got %d", i, leaf.point.Idx)
		}
		if leaf.point != points[i] {
			t.Fatalf("Expected point to be %v, got %v", points[i], leaf.point)
		}
	}
}
