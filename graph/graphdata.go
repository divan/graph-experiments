package graph

import (
	"encoding/json"
	"os"
)

type Data struct {
	Nodes []*NodeData `json:"nodes"`
	Links []*LinkData `json:"links"`

	nodeHasLinks map[string]bool
}

type NodeData struct {
	ID    string `json:"id"`
	Group int    `json:"group"`
}

type LinkData struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	FromIdx int    `json:"-"`
	ToIdx   int    `json:"-"`
}

func NewData() *Data {
	return &Data{}
}

func NewDataFromJSON(file string) (*Data, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var g *Data
	err = json.NewDecoder(fd).Decode(&g)
	if err != nil {
		return nil, err
	}

	// add node indexes to links
	nodeIndexes := make(map[string]int)
	for i, _ := range g.Nodes {
		nodeIndexes[g.Nodes[i].ID] = i
	}

	g.nodeHasLinks = make(map[string]bool)
	for i, link := range g.Links {
		g.nodeHasLinks[link.Source] = true
		g.Links[i].FromIdx = nodeIndexes[link.Source]
		g.Links[i].ToIdx = nodeIndexes[link.Target]
	}

	return g, err
}

func (g *Data) NodeHasLinks(nodeID string) bool {
	return g.nodeHasLinks[nodeID]
}
