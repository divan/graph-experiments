package layout

import (
	"fmt"
	"math"

	"github.com/divan/graph-experiments/graph"
)

type Layout interface {
	Nodes() []*Node
	Calculate(iterations int)

	AddForce(Force)
	ListForces() []Force
}

type LayoutWithDebug interface {
	Layout
	ForceValues() map[int][]*ForceVector
}

type Layout3D struct {
	nodes  []*Node
	links  []*graph.LinkData
	forces []Force

	forceValues map[int][]*ForceVector
}

// Init initializes layout with nodes data. It assigns
// semi-random positions to nodes to facilitate further simulation.
func New(data *graph.Data, forces ...Force) LayoutWithDebug {
	nodes := generateRandomPositions(data.Nodes)

	return &Layout3D{
		nodes:       nodes,
		links:       data.Links,
		forces:      forces,
		forceValues: make(map[int][]*ForceVector),
	}

}

// generateRandomPositions returns new nodes array with (semi)random
// positions visually distributed in a 3D space.
func generateRandomPositions(nodes []*graph.NodeData) []*Node {
	ret := make([]*Node, 0, len(nodes))

	for i := range nodes {
		radius := 10 * math.Cbrt(float64(i))
		rollAngle := float64(float64(i) * math.Pi * (3 - math.Sqrt(5))) // golden angle
		yawAngle := float64(float64(i) * math.Pi / 24)                  // sequential (divan: wut?)

		node := &Node{
			Point: &Point{
				X:    int32(radius * math.Cos(rollAngle)),
				Y:    int32(radius * math.Sin(rollAngle)),
				Z:    int32(radius * math.Sin(yawAngle)),
				Mass: nodes[i].Weight + 2,
			},
			ID: nodes[i].ID,
		}
		ret = append(ret, node)
	}

	return ret
}

func (l *Layout3D) Nodes() []*Node {
	return l.nodes
}

func (l *Layout3D) ForceValues() map[int][]*ForceVector {
	return l.forceValues
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

	forceValues := make([]*ForceVector, len(l.nodes))
	l.resetForces()

	l.applyRepulsion(ot, forceValues)
	l.applySprings(forceValues)

	l.integrate(forceValues)
}

func (l *Layout3D) resetForces() {
	l.forceValues = make(map[int][]*ForceVector)
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
	l.forceValues[nodeIdx] = append(l.forceValues[nodeIdx], f)
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

func (l *Layout3D) AddForce(f Force) {
	l.forces = append(l.forces, f)
}

func (l *Layout3D) ListForces() []Force {
	return l.forces
}
