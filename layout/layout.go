package layout

import (
	"fmt"
	"math"

	"github.com/divan/graph-experiments/graph"
)

type Node struct {
	ID string
	*Point
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
			Point: &Point{
				X:    int32(radius * math.Cos(rollAngle)),
				Y:    int32(radius * math.Sin(rollAngle)),
				Z:    int32(radius * math.Sin(yawAngle)),
				Mass: data.Nodes[i].Weight + 2,
			},
			ID: data.Nodes[i].ID,
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
		fmt.Println("Inserting point", p)
		ot.Insert(p)
	}

	forces := make([]*force, len(l.nodes))

	l.applyRepulsion(ot, forces)
	l.applySprings(forces)

	l.integrate(forces)
	for i, _ := range l.nodes {
		fmt.Printf("i==> Node [%d]: %v\n", i, l.nodes[i])
	}
}

func (l *Layout3D) applyRepulsion(ot *Octree, forces []*force) {
	// anti-gravity repelling
	for i := range l.nodes {
		f := &force{}
		gf, err := ot.CalcForce(i)
		if err != nil {
			fmt.Println("[ERROR] Force calc failed:", i, err)
			forces[i] = f
			continue
		}
		forces[i] = f.Add(gf)
	}
}

func (l *Layout3D) applySprings(forces []*force) {
	for _, link := range l.links {
		f := springForce(l.nodes[link.FromIdx], l.nodes[link.ToIdx])
		forces[link.FromIdx] = forces[link.FromIdx].Add(f)
		forces[link.ToIdx] = forces[link.ToIdx].Sub(f)
	}
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
