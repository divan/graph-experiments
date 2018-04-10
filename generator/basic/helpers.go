package basic

import (
	"fmt"

	"github.com/divan/graph-experiments/graph"
)

func addNode(data *graph.Data, i int) {
	node := &graph.Node{
		ID: idxToID(i),
	}
	data.Nodes = append(data.Nodes, node)
}

func addLink(data *graph.Data, i, j int) {
	link := &graph.Link{
		Source: idxToID(i),
		Target: idxToID(j),
	}
	data.Links = append(data.Links, link)
}

// idxToID creates ID to be used with Node struct from index.
func idxToID(i int) string {
	return fmt.Sprintf("Node %d", i)
}
