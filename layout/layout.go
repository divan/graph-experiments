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
	ForceDebugInfo() map[int][]*ForceDebugInfo
}

type Layout3D struct {
	nodes  []*Node
	links  []*graph.LinkData
	forces []Force

	forceVectors   map[int]*ForceVector // cumulative force per node ID
	forceDebugInfo map[int][]*ForceDebugInfo
}

// Init initializes layout with nodes data. It assigns
// semi-random positions to nodes to facilitate further simulation.
func New(data *graph.Data, forces ...Force) LayoutWithDebug {
	nodes := generateRandomPositions(data.Nodes)

	return &Layout3D{
		nodes:          nodes,
		links:          data.Links,
		forces:         forces,
		forceVectors:   make(map[int]*ForceVector),
		forceDebugInfo: make(map[int][]*ForceDebugInfo),
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

func (l *Layout3D) ForceDebugInfo() map[int][]*ForceDebugInfo {
	return l.forceDebugInfo
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
	//ot := NewOctreeFromNodes(l.Nodes())

	l.resetForces()

	for _, f := range l.forces {
		_ = f
		//f.Apply()
	}
	//l.applyRepulsion(ot, forceValues)
	//l.applySprings(forceValues)

	//l.integrate(forceValues)
}

func (l *Layout3D) resetForces() {
	l.forceVectors = make(map[int]*ForceVector)
	l.forceDebugInfo = make(map[int][]*ForceDebugInfo)
}

func (l *Layout3D) applyRepulsion(ot *Octree) {
	// anti-gravity repelling
	for i, node := range l.nodes {
		gf, err := ot.CalcForce(i)
		if err != nil {
			fmt.Println("[ERROR] Force calc failed:", i, err)
			continue
		}
		l.forceVectors[node.Idx] = gf

		// attach debug information
		f := &ForceDebugInfo{}
		l.appendForce(node.Idx, f)
	}
}

func (l *Layout3D) applySprings() {
	for _, link := range l.links {
		f := defaultSpringForce.Apply(l.nodes[link.FromIdx].Point, l.nodes[link.ToIdx].Point)
		l.forceVectors[link.FromIdx] = l.forceVectors[link.FromIdx].Add(f)
		l.forceVectors[link.ToIdx] = l.forceVectors[link.ToIdx].Sub(f)

		//l.appendForce(link.FromIdx, &ForceDebugInfo{"spring", f})
		//l.appendForce(link.ToIdx, (&ForceVector{}).Sub(f))
	}
}

// appendForce append debug information about applied force to debug info.
func (l *Layout3D) appendForce(nodeIdx int, f *ForceDebugInfo) {
	l.forceDebugInfo[nodeIdx] = append(l.forceDebugInfo[nodeIdx], f)
}

func (l *Layout3D) AddForce(f Force) {
	l.forces = append(l.forces, f)
}

func (l *Layout3D) ListForces() []Force {
	return l.forces
}
