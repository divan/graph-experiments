package graph

import "fmt"

// NodeHasLinks implements fast check if given node has any links.
func (g *Graph) NodeHasLinks(idx int) bool {
	if g.nodeLinks == nil {
		g.prepare()
	}

	return g.nodeLinks[idx] > 0
}

// NodeLinks returns number of links for node.
func (g *Graph) NodeLinks(idx int) int {
	if g.nodeLinks == nil {
		g.prepare()
	}

	return g.nodeLinks[idx]
}

// NodeIDHasLinks implements fast check if given node by ID has any links.
func (g *Graph) NodeIDHasLinks(id string) bool {
	if g.nodeLinks == nil {
		g.prepare()
	}

	idx, err := g.NodeByID(id)
	if err != nil {
		return false
	}

	return g.nodeLinks[idx] > 0
}

// NodeIDLinks returns number of links for node ID.
func (g *Graph) NodeIDLinks(id string) int {
	if g.nodeLinks == nil {
		g.prepare()
	}

	idx, err := g.NodeByID(id)
	if err != nil {
		return 0
	}

	return g.nodeLinks[idx]
}

// LinkExists returns true if there is a link between source and target.
func (g *Graph) LinkExists(from, to int) bool {
	for _, link := range g.links {
		if link.From == from && link.To == to ||
			link.To == from && link.From == to {
			return true
		}
	}
	return false
}

// LinkExistsByID returns true if there is a link between source and target node IDs.
func (g *Graph) LinkExistsByID(source, target string) bool {
	for _, link := range g.links {
		from, ok := g.nodeIndexes[source]
		if !ok {
			continue
		}
		to, ok := g.nodeIndexes[target]
		if !ok {
			continue
		}
		if link.From == from && link.To == to ||
			link.To == from && link.From == to {
			return true
		}
	}
	return false
}

// NodeByID returns node index by its ID.
func (g *Graph) NodeByID(id string) (int, error) {
	for i, node := range g.nodes {
		if node.ID() == id {
			return i, nil
		}
	}
	return 0, fmt.Errorf("node %s not found", id)
}

// NodeIDByIdx returns node ID by its index.
func (g *Graph) NodeIDByIdx(idx int) (string, error) {
	if idx < 0 || idx > len(g.nodes)-1 {
		return "", fmt.Errorf("node for index %d not found", idx)
	}

	return g.nodes[idx].ID(), nil
}

// prepare runs various optimization-related
// calculations, caching etc.
func (g *Graph) prepare() {
	// add node indexes to links
	g.nodeIndexes = make(map[string]int)
	for i, node := range g.nodes {
		g.nodeIndexes[node.ID()] = i
	}

	g.nodeLinks = make(map[int]int)
	for _, link := range g.links {
		g.nodeLinks[link.From]++
		g.nodeLinks[link.To]++
	}
}
