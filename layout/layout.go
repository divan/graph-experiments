package layout

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/divan/graph-experiments/graph"
	"gopkg.in/cheggaaa/pb.v1"
)

// stableThreshold determines the movement diff needed to
// call the system stable
const stableThreshold = 2.001

// Layout represents physical layout used to process graph.
type Layout interface {
	Nodes() []*Node
	Calculate()
	CalculateN(iterations int)
	Reset()

	AddForce(Force)
	ListForces() []Force
}

// LayoutWithDebug extends Layout interface with additional debug data methods.
// TODO(divan): this should be a static type, not an interface.
type LayoutWithDebug interface {
	Layout
	ForcesDebugData() ForcesDebugData
}

// Layout3D implements Layout interface for force-directed 3D graph.
type Layout3D struct {
	data *graph.Graph

	nodes  []*Node
	links  []*graph.Link
	forces []Force

	forceVectors    map[int]*ForceVector // cumulative force per node ID
	forcesDebugData ForcesDebugData
}

// New initializes layout with nodes data.
func New(data *graph.Graph, forces ...Force) LayoutWithDebug {
	l := &Layout3D{
		data:            data,
		links:           data.Links(),
		forces:          forces,
		forceVectors:    make(map[int]*ForceVector),
		forcesDebugData: make(ForcesDebugData),
	}

	l.Reset()

	return l
}

// Reset resets positions to the (semi)pseudorandom positions and cancels all
// forces and velocities.
func (l *Layout3D) Reset() {
	l.nodes = make([]*Node, 0, len(l.data.Nodes()))

	for i := range l.data.Nodes() {
		radius := 10 * math.Cbrt(float64(i))
		rollAngle := float64(float64(i) * math.Pi * (3 - math.Sqrt(5))) // golden angle
		yawAngle := float64(float64(i) * math.Pi / 24)                  // sequential (divan: wut?)

		var weight int = 1
		node := l.data.Nodes()[i]
		if wnode, ok := node.(graph.WeightedNode); ok {
			weight = wnode.Weight()
		}

		newNode := &Node{
			Point: &Point{
				Idx:  i,
				X:    int(radius * math.Cos(rollAngle)),
				Y:    int(radius * math.Sin(rollAngle)),
				Z:    int(radius * math.Sin(yawAngle)),
				Mass: weight,
			},
			ID: node.ID(),
		}
		l.nodes = append(l.nodes, newNode)
	}

	l.resetForces()
}

// Calculate runs positions' recalculations iteratively until the
// system minimizes it's energy.
func (l *Layout3D) Calculate() {
	// tx is the total movement, which should drop to the minimum
	// at the minimal energy state
	fmt.Println("Simulation started...")
	var (
		now    = time.Now()
		count  int
		prevTx float64
	)
	for tx := math.MaxFloat64; math.Abs(tx-prevTx) >= stableThreshold; {
		prevTx = tx
		tx = l.UpdatePositions()
		log.Println("PrevTx, tx:", tx, ", diff:", math.Abs(tx-prevTx))
		count++
		if count%1000 == 0 {
			since := time.Since(now)
			fmt.Printf("Iterations: %d, tx: %f, time: %v\n", count, tx, since)
		}
	}
	fmt.Printf("Simulation finished in %v, run %d iterations\n", time.Since(now), count)
}

// CalculateN run positions' recalculations exactly N times.
func (l *Layout3D) CalculateN(n int) {
	fmt.Println("Simulation started...")
	bar := pb.StartNew(n)
	for i := 0; i < n; i++ {
		l.UpdatePositions()
		bar.Increment()
	}
	bar.FinishPrint("Simulation finished")

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

// AddForce adds force to the internal list of forces.
func (l *Layout3D) AddForce(f Force) {
	l.forces = append(l.forces, f)
}

// ListForces returns the list of active forces.
func (l *Layout3D) ListForces() []Force {
	return l.forces
}

// Nodes returns nodes information.
func (l *Layout3D) Nodes() []*Node {
	return l.nodes
}

// ForcesDebugData returns debug data with forces.
func (l *Layout3D) ForcesDebugData() ForcesDebugData {
	return l.forcesDebugData
}

// Links returns graph data links.
func (l *Layout3D) Links() []*graph.Link {
	return l.links
}
