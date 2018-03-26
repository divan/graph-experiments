package layout

import (
	"fmt"
	"math"

	"github.com/divan/graph-experiments/graph"
)

// stableThreshold determines the movement diff needed to
// call the system stable
const stableThreshold = 20.0

type Layout interface {
	Nodes() []*Node
	Calculate()
	CalculateN(iterations int)

	AddForce(Force)
	ListForces() []Force
}

type LayoutWithDebug interface {
	Layout
	ForcesDebugData() ForcesDebugData
}

type Layout3D struct {
	nodes  []*Node
	links  []*graph.LinkData
	forces []Force

	forceVectors    map[int]*ForceVector // cumulative force per node ID
	forcesDebugData ForcesDebugData
}

// Init initializes layout with nodes data. It assigns
// semi-random positions to nodes to facilitate further simulation.
func New(data *graph.Data, forces ...Force) LayoutWithDebug {
	nodes := generateRandomPositions(data.Nodes)

	return &Layout3D{
		nodes:           nodes,
		links:           data.Links,
		forces:          forces,
		forceVectors:    make(map[int]*ForceVector),
		forcesDebugData: make(ForcesDebugData),
	}
}

// Calculate runs positions' recalculations iteratively until the
// system minimizes it's energy.
func (l *Layout3D) Calculate() {
	// tx is the total movement, which should drop to the minimum
	// at the minimal energy state
	var count int
	for tx := math.MaxFloat64; tx >= stableThreshold; {
		tx = l.UpdatePositions()
		count++
		if count%1000 == 0 {
			fmt.Printf("Iterations: %d, tx: %f\n", count, tx)
		}
	}
	fmt.Println("Total Iterations:", count)
}

// CalculateN run positions' recalculations exactly N times.
func (l *Layout3D) CalculateN(n int) {
	for i := 0; i < n; i++ {
		l.UpdatePositions()
	}
}

// UpdatePositions recalculates nodes' positions, applying all the forces.
// It returns average amount of movement generated by this step.
func (l *Layout3D) UpdatePositions() float64 {
	l.resetForces()

	for _, f := range l.forces {
		applyForces := f.ByRule() // TODO: rename to Rule
		applyForces(f, l.nodes, l.links, l.forceVectors, l.forcesDebugData)
	}

	return l.integrate()
}

func (l *Layout3D) resetForces() {
	l.forceVectors = make(map[int]*ForceVector)
	l.forcesDebugData = make(ForcesDebugData)
}

func (l *Layout3D) AddForce(f Force) {
	l.forces = append(l.forces, f)
}

func (l *Layout3D) ListForces() []Force {
	return l.forces
}

func (l *Layout3D) Nodes() []*Node {
	return l.nodes
}

func (l *Layout3D) ForcesDebugData() ForcesDebugData {
	return l.forcesDebugData
}

func (l *Layout3D) Links() []*graph.LinkData {
	return l.links
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
				Idx:  i,
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
