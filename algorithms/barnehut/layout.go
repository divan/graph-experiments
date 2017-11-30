package main

import (
	"fmt"
	"math"
)

type Node struct {
	X         int32
	Y         int32
	Z         int32
	VelocityX float64
	VelocityY float64
	VelocityZ float64
}

func (n *Node) String() string {
	return fmt.Sprintf("[%.2f, %.2f, %.2f]", n.X, n.Y, n.Z)
}

type Layout interface {
	InitCoordinates(data []*NodeData)
	Nodes() []*Node
}

type Layout3D struct {
	nodes []*Node
}

func (l *Layout3D) InitCoordinates(data []*NodeData) {
	nodes := make([]*Node, 0, len(data))
	for i, _ := range data {
		radius := 10 * math.Cbrt(float64(i))
		rollAngle := float64(float64(i) * math.Pi * (3 - math.Sqrt(5))) // golden angle
		yawAngle := float64(float64(i) * math.Pi / 24)                  // sequential (divan: wut?)

		node := &Node{
			X: int32(radius * math.Cos(rollAngle)),
			Y: int32(radius * math.Sin(rollAngle)),
			Z: int32(radius * math.Sin(yawAngle)),
		}
		nodes = append(nodes, node)
	}
	l.nodes = nodes
}

func (l *Layout3D) Nodes() []*Node {
	return l.nodes
}
