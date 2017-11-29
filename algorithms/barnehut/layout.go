package main

import (
	"fmt"
	"math"
)

type Node struct {
	X         float64
	Y         float64
	Z         float64
	VelocityX float64
	VelocityY float64
	VelocityZ float64
}

func (n *Node) String() string {
	return fmt.Sprintf("[%.2f, %.2f, %.2f]", n.X, n.Y, n.Z)
}

type Layout interface {
	InitCoordinates(data []*NodeData) []*Node
}

type Layout3D struct{}

func (l *Layout3D) InitCoordinates(data []*NodeData) []*Node {
	nodes := make([]*Node, 0, len(data))
	for i, _ := range data {
		radius := 10 * math.Cbrt(float64(i))
		rollAngle := float64(float64(i) * math.Pi * (3 - math.Sqrt(5))) // golden angle
		yawAngle := float64(float64(i) * math.Pi / 24)                  // sequential (divan: wut?)

		node := &Node{
			X: radius * math.Cos(rollAngle),
			Y: radius * math.Sin(rollAngle),
			Z: radius * math.Sin(yawAngle),
		}
		nodes = append(nodes, node)
	}
	return nodes
}
