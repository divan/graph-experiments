package layout

import (
	"fmt"
	"github.com/divan/graph-experiments/graph"
	"math"
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

type force struct {
	dx, dy, dz float64
}

func (f force) String() string {
	return fmt.Sprintf("f(%.03f, %.03f, %.03f)", f.dx, f.dy, f.dz)
}

// Add adds new force to f.
func (f *force) Add(f1 *force) *force {
	f.dx += f1.dx
	f.dy += f1.dy
	f.dz += f1.dz
	return f
}

// Sub substracts new force from f.
func (f *force) Sub(f1 *force) *force {
	f.dx -= f1.dx
	f.dy -= f1.dy
	f.dz -= f1.dz
	return f
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
	findNodeIdx := func(id string) int {
		for i := range l.nodes {
			if l.nodes[i].ID == id {
				return i
			}
		}
		fmt.Println("Can't find node for id:", id)
		return -1
	}
	for _, link := range l.links {
		fromIdx := findNodeIdx(link.Source)
		if fromIdx == -1 {
			continue
		}
		toIdx := findNodeIdx(link.Target)
		if toIdx == -1 {
			continue
		}

		f1, f2 := forces[fromIdx], forces[toIdx]
		f := springForce(l.nodes[fromIdx], l.nodes[toIdx])
		forces[fromIdx] = f1.Add(f)
		forces[toIdx] = f2.Sub(f)
	}

	l.integrate(forces)
}

func springForce(from, to *Node) *force {
	dx := float64(to.X - from.X)
	dy := float64(to.Y - from.Y)
	dz := float64(to.Z - from.Z)
	r := math.Sqrt(dx*dx + dy*dy + dz*dz)

	if r == 0 {
		r = 10
	}

	const (
		length = 20
		coeff  = 0.0008
		weight = 1
	)

	d := r - length
	c := coeff * d / r * weight

	return &force{
		dx: c * dx,
		dy: c * dy,
		dz: c * dz,
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

// integrate performs forces integration using Euler numerical
// integration method.
func (l *Layout3D) integrate(forces []*force) {
	const timeStep = float64(20) // FIXME: 20 what?
	for i := 0; i < len(l.nodes); i++ {
		body := l.nodes[i]
		coeff := timeStep / float64(body.Mass)

		vx := coeff * forces[i].dx
		vy := coeff * forces[i].dy
		vz := coeff * forces[i].dz
		v := math.Sqrt(vx*vx + vy*vy + vz*vz)

		if v > 1 {
			vx = vx / v
			vy = vy / v
			vz = vz / v
		}

		dx := timeStep * vx
		dy := timeStep * vy
		dz := timeStep * vz

		l.nodes[i].X += int32(dx)
		l.nodes[i].Y += int32(dy)
		l.nodes[i].Z += int32(dz)
	}
}
