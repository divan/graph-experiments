package main

import "fmt"

// Node represents whisper node in graph.
type Node struct {
	ID_     string `json:"id"`
	Group_  int    `json:"group,omitempty"`
	Weight_ int    `json:"weight,omitempty"`

	isClient bool
}

func NewNode(id string, client bool) *Node {
	var (
		group  = 1
		weight = 3
	)
	if client {
		group = 2
		weight = 1
	}
	return &Node{
		ID_:      id,
		Group_:   group,
		Weight_:  weight,
		isClient: client,
	}
}

// ID implements Node for Node.
func (n *Node) ID() string { return n.ID_ }

// Group implements GroupNode for Node.
func (n *Node) Group() int {
	return n.Group_
}

// Weight implements WeightedNode for Node.
func (n *Node) Weight() int {
	return n.Weight_
}

// String implements Stringer for Node.
func (n *Node) String() string {
	return fmt.Sprintf("%s [%d]", n.ID_, n.Group_)
}

func (n *Node) IsClient() bool {
	return n.isClient
}
