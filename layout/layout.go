package layout

import (
	"fmt"
	"math"

	"github.com/divan/graph-experiments/graph"
)

type Node struct {
	X    int32
	Y    int32
	Z    int32
	Mass int32
	ID   string
}

func (n *Node) String() string {
	return fmt.Sprintf("[%d, %d, %d, m: %d]", n.X, n.Y, n.Z, n.Mass)
}

type Layout interface {
	Nodes() []*Node
	Calculate(iterations int)
}

type Layout3D struct {
	nodes []*Node
	links []*graph.LinkData
}

// Init initializes layout with nodes data. It assigns
// semi-random positions to nodes to facilitate further simulation.
func New(data *graph.Data) Layout {
	l := &Layout3D{}

	nodes := make([]*Node, 0, len(data.Nodes))
	for i, _ := range data.Nodes {
		radius := 10 * math.Cbrt(float64(i))
		rollAngle := float64(float64(i) * math.Pi * (3 - math.Sqrt(5))) // golden angle
		yawAngle := float64(float64(i) * math.Pi / 24)                  // sequential (divan: wut?)

		node := &Node{
			X:    int32(radius * math.Cos(rollAngle)),
			Y:    int32(radius * math.Sin(rollAngle)),
			Z:    int32(radius * math.Sin(yawAngle)),
			Mass: 1,
			ID:   data.Nodes[i].ID,
		}
		nodes = append(nodes, node)
	}
	l.nodes = nodes
	l.links = data.Links

	return l
}

func (l *Layout3D) Nodes() []*Node {
	return l.nodes
}

func (l *Layout3D) Links() []*graph.LinkData {
	return l.links
}

func (l *Layout3D) Calculate(n int) {
	for i := 0; i < n; i++ {
		l.updatePositions()
	}
}

func (l *Layout3D) updatePositions() {
	// insert current node positions into octree
	ot := NewOctree()
	for idx, node := range l.Nodes() {
		p := newPointFromNode(idx, node)
		ot.Insert(p)
	}

	forces := make([]*force, len(l.nodes))
	// anti-gravity repelling
	for i := range l.nodes {
		f := &force{}
		f.dx, f.dy, f.dz = ot.CalcForce(i)
		forces[i] = f
	}

	// link springs
	for _, link := range l.links {
		f1, f2 := forces[link.FromIdx], forces[link.ToIdx]
		f := springForce(l.nodes[link.FromIdx], l.nodes[link.ToIdx])
		forces[link.FromIdx] = f1.Add(f)
		forces[link.ToIdx] = f2.Sub(f)
	}

	l.integrate(forces)
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
