package graph

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

// Graph represents graph data - nodes and links.
type Graph struct {
	nodes []Node
	links []*Link

	nodeIndexes map[string]int
	nodeLinks   map[int]int
}

// NewGraph creates empty graph data.
func NewGraph() *Graph {
	return &Graph{}
}

// NewGraphFromJSON creates a graph from the given JSON file.
// TODO: add support for relfection and custom node types unmarshalling.
func NewGraphFromJSON(file string) (*Graph, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close() //nolint: errcheck

	return NewGraphFromJSONReader(fd)
}

// NewGraphFromJSONReader creates a graph from the given JSON file.
func NewGraphFromJSONReader(r io.Reader) (*Graph, error) {
	// decode into temporary struct to process
	var res struct {
		Nodes []*BasicNode `json:"nodes"`
		Links []*struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"links"`
	}
	err := json.NewDecoder(r).Decode(&res)
	if err != nil {
		return nil, err
	}

	// convert links IDs into indices
	g := &Graph{
		nodes: make([]Node, 0, len(res.Nodes)),
		links: make([]*Link, 0, len(res.Links)),
	}

	for _, node := range res.Nodes {
		g.AddNode(node)
	}

	for _, link := range res.Links {
		err := g.AddLinkByIDs(link.Source, link.Target)
		if err != nil {
			return nil, err
		}
	}

	g.prepare()

	if len(g.nodes) == 0 {
		return nil, errors.New("empty graph")
	}

	return g, err
}

// Nodes returns graph nodes
func (g *Graph) Nodes() []Node {
	return g.nodes
}

// Links returns graph links.
func (g *Graph) Links() []*Link {
	return g.links
}
