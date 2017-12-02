package main

import (
	"fmt"
	"math"

	"github.com/divan/graph-experiments/algorithms/barnehut/octree"
)

type Node struct {
	X         int32
	Y         int32
	Z         int32
	Mass      int32
	VelocityX float64
	VelocityY float64
	VelocityZ float64
}

func (n *Node) String() string {
	return fmt.Sprintf("[%d, %d, %d, m: %d]", n.X, n.Y, n.Z, n.Mass)
}

type Layout interface {
	Init(data []*NodeData)
	Nodes() []*Node
	Calculate(iterations int)
}

type Layout3D struct {
	nodes []*Node
}

// Init initializes layout with nodes data. It assigns
// semi-random positions to nodes to facilitate further simulation.
func (l *Layout3D) Init(data []*NodeData) {
	nodes := make([]*Node, 0, len(data))
	for i, _ := range data {
		radius := 10 * math.Cbrt(float64(i))
		rollAngle := float64(float64(i) * math.Pi * (3 - math.Sqrt(5))) // golden angle
		yawAngle := float64(float64(i) * math.Pi / 24)                  // sequential (divan: wut?)

		node := &Node{
			X:    int32(radius * math.Cos(rollAngle)),
			Y:    int32(radius * math.Sin(rollAngle)),
			Z:    int32(radius * math.Sin(yawAngle)),
			Mass: 1,
		}
		nodes = append(nodes, node)
	}
	l.nodes = nodes
}

func (l *Layout3D) Nodes() []*Node {
	return l.nodes
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
func (f *force) Add(f1 *force) {
	f.dx += f1.dx
	f.dy += f1.dy
	f.dz += f1.dz
}

func (l *Layout3D) updatePositions() {
	// insert current node positions into octree
	ot := octree.New()
	for idx, node := range l.Nodes() {
		p := newPointFromNode(idx, node)
		ot.Insert(p)
	}

	var forces []force
	for i, _ := range l.Nodes() {
		// calculate force for i-th node
		var f force
		f.dx, f.dy, f.dz = ot.CalcForce(i)
		forces = append(forces, f)
	}

	l.integrate(forces)
}

func newPointFromNode(idx int, n *Node) *octree.Point {
	return &octree.Point{
		Idx:  idx,
		X:    n.X,
		Y:    n.Y,
		Z:    n.Z,
		Mass: n.Mass,
	}
}

// integrate performs forces integration using Euler numerical
// integration method.
func (l *Layout3D) integrate(forces []force) {
	const timeStep = float64(20) // FIXME: 20 what?
	for i := 0; i < len(l.Nodes()); i++ {
		body := l.Nodes()[i]
		coeff := timeStep / float64(body.Mass)

		body.VelocityX += coeff * forces[i].dx
		body.VelocityY += coeff * forces[i].dy
		body.VelocityZ += coeff * forces[i].dz
		v := math.Sqrt(body.VelocityX*body.VelocityX + body.VelocityY*body.VelocityY + body.VelocityZ*body.VelocityZ)

		if v > 1 {
			body.VelocityX = body.VelocityX / v
			body.VelocityY = body.VelocityY / v
			body.VelocityZ = body.VelocityZ / v
		}

		dx := timeStep * body.VelocityX
		dy := timeStep * body.VelocityY
		dz := timeStep * body.VelocityZ

		l.Nodes()[i].X += int32(dx)
		l.Nodes()[i].Y += int32(dy)
		l.Nodes()[i].Z += int32(dz)
	}
}
