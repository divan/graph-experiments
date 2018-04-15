package graph

// Node defines the graph node. Any type implementing this
// interface can be used as a graph node.
type Node interface {
	ID() string
}

// GroupedNode represents node that have 'group' attribute.
type GroupedNode interface {
	Group() int
}

// WeightedNode represents node that have 'weight' attribute.
type WeightedNode interface {
	Weight() int
}

// BasicNode represents basic built-in node type for simple cases.
type BasicNode struct {
	id     string
	group  int
	weight int
}

// ID implements Node for BasicNode.
func (b *BasicNode) ID() string { return b.id }

// Group implements GroupNode for BasicNode.
func (b *BasicNode) Group() int { return b.group }

// Weight implements WeightedNode for BasicNode.
func (b *BasicNode) Weight() int { return b.weight }
