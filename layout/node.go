package layout

import "fmt"

// Node represents point in 3D space with some ID information
// attached to it.
type Node struct {
	ID string
	*Point
}

// String implements Stringer interface for Node.
func (n *Node) String() string {
	return fmt.Sprintf("[%d, %d, %d, m: %d]", n.X, n.Y, n.Z, n.Mass)
}

func newPointFromNode(idx int, n *Node) *Point {
	return &Point{
		Idx:  idx,
		X:    n.X,
		Y:    n.Y,
		Z:    n.Z,
		Mass: n.Mass,
	}
}

// Velocity represents velocity vector.
type Velocity struct {
	X float64
	Y float64
	Z float64
}
