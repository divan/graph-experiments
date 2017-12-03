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
	Source string `json:"source"`
	Target string `json:"target"`
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

	g.nodeHasLinks = make(map[string]bool)
	for _, link := range g.Links {
		g.nodeHasLinks[link.Source] = true
	}

	return g, err
}

func (g *Data) NodeHasLinks(nodeID string) bool {
	return g.nodeHasLinks[nodeID]
}
