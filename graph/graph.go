package graph

import (
	"encoding/json"
	"os"
)

// Data represents graph data - nodes and links.
type Data struct {
	Nodes []*Node `json:"nodes"`
	Links []*Link `json:"links"`

	nodeHasLinks map[string]bool
}

// Node represents single node of the graph.
type Node struct {
	ID     string `json:"id"`
	Group  int    `json:"group,omitempty"`
	Weight int32  `json:"weight,omitempty"`
}

// Link represents single link between two nodes.
type Link struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	FromIdx int    `json:"-"`
	ToIdx   int    `json:"-"`
}

// NewData creates empty graph data.
func NewData() *Data {
	return &Data{}
}

// NewDataFromJSON creates a graph from the given JSON file.
func NewDataFromJSON(file string) (*Data, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var d *Data
	err = json.NewDecoder(fd).Decode(&d)
	if err != nil {
		return nil, err
	}

	d.prepare()

	return d, err
}

// prepare runs various optimization-related
// calculateions, caching etc.
func (d *Data) prepare() {
	// add node indexes to links
	nodeIndexes := make(map[string]int)
	for i, _ := range d.Nodes {
		nodeIndexes[d.Nodes[i].ID] = i
	}

	d.nodeHasLinks = make(map[string]bool)
	for i, link := range d.Links {
		d.nodeHasLinks[link.Source] = true
		d.Links[i].FromIdx = nodeIndexes[link.Source]
		d.Links[i].ToIdx = nodeIndexes[link.Target]
	}
}
