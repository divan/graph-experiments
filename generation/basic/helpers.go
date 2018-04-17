package basic

import (
	"fmt"

	"github.com/divan/graph-experiments/graph"
)

func addNode(g *graph.Graph, i int) {
	node := &graph.BasicNode{
		ID_: idxToID(i),
	}

	g.AddNode(node)
}

// idxToID creates ID to be used with Node struct from index.
func idxToID(i int) string {
	return fmt.Sprintf("Node %d", i)
}
