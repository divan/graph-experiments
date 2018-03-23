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
	Forces() map[int][]*ForceVector
}

type Layout3D struct {
	nodes []*Node
	links []*graph.LinkData

	forces map[int][]*ForceVector
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

	l.forces = make(map[int][]*ForceVector)

	return l
}

func (l *Layout3D) Nodes() []*Node {
	return l.nodes
}

func (l *Layout3D) Forces() map[int][]*ForceVector {
	return l.forces
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
	ot := NewOctreeFromNodes(l.Nodes())

	forces := make([]*ForceVector, len(l.nodes))
	l.resetForces()

	l.applyRepulsion(ot, forces)
	l.applySprings(forces)

	l.integrate(forces)
}

func (l *Layout3D) resetForces() {
	l.forces = make(map[int][]*ForceVector)
}

func (l *Layout3D) applyRepulsion(ot *Octree, forces []*ForceVector) {
	// anti-gravity repelling
	for i := range l.nodes {
		f := &ForceVector{}
		gf, err := ot.CalcForce(i)
		if err != nil {
			fmt.Println("[ERROR] Force calc failed:", i, err)
			forces[i] = f
			continue
		}
		forces[i] = f.Add(gf)

		l.appendForce(i, forces[i])
	}
}

func (l *Layout3D) applySprings(forces []*ForceVector) {
	for _, link := range l.links {
		f := defaultSpringForce.Apply(l.nodes[link.FromIdx], l.nodes[link.ToIdx])
		forces[link.FromIdx] = forces[link.FromIdx].Add(f)
		forces[link.ToIdx] = forces[link.ToIdx].Sub(f)

		l.appendForce(link.FromIdx, f)
		l.appendForce(link.ToIdx, (&ForceVector{}).Sub(f))
	}
}

func (l *Layout3D) appendForce(nodeIdx int, f *ForceVector) {
	l.forces[nodeIdx] = append(l.forces[nodeIdx], f)
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
