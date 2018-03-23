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
