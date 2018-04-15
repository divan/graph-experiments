package graph

import "fmt"

// Link represents single link between two nodes.
type Link struct {
	From int
	To   int
}

// NewLink constructs new Link object.
// Note, this function doesn't know actual nodes, so it doesn't
// check for indices validity. If you need validate indices, use
// graph.AddLink() instead.
func NewLink(from, to int) *Link {
	return &Link{
		From: from,
		To:   to,
	}
}

// AddLink adds new link to the graph and validates input
// indices.
func (g *Graph) AddLink(from, to int) error {
	if from > len(g.nodes) {
		return fmt.Errorf("Node not found: %v", from)
	}
	if to > len(g.nodes) {
		return fmt.Errorf("Node not found: %v", to)
	}

	link := NewLink(from, to)
	g.links = append(g.links, link)
	return nil
}

// AddLinkByIDs adds new link to the graph by node IDs and
// valides the input.
func (g *Graph) AddLinkByIDs(fromID, toID string) error {
	from, err := g.findNodeByID(fromID)
	if err != nil {
		return err
	}
	to, err := g.findNodeByID(toID)
	if err != nil {
		return err
	}

	link := NewLink(from, to)
	g.links = append(g.links, link)
	return nil
}

func (g *Graph) findNodeByID(id string) (int, error) {
	for i, node := range g.nodes {
		if node.ID() == id {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Node %s not found", id)
}
