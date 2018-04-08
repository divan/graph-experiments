package basic

import (
	"fmt"

	"github.com/divan/graph-experiments/graph"
)

// CircleGenerator implements generator for circle graph.
type CircleGenerator struct {
	nodes int // number of nodes
}

// NewCircleGenerator creates new line generator for N nodes graph.
func NewCircleGenerator(n int) *CircleGenerator {
	return &CircleGenerator{
		nodes: n,
	}
}

// Generate generates the data for graph. Implements Generator interface.
func (l *CircleGenerator) Generate() *graph.Data {
	data := graph.NewData()

	for i := 0; i < l.nodes; i++ {
		node := &graph.Node{
			ID: l.idxToID(i),
		}
		data.Nodes = append(data.Nodes, node)

		targetIdx := i + 1
		if i == l.nodes-1 {
			targetIdx = 0
		}
		link := &graph.Link{
			Source: l.idxToID(i),
			Target: l.idxToID(targetIdx),
		}
		data.Links = append(data.Links, link)
	}

	return data
}

func (l *CircleGenerator) idxToID(i int) string {
	return fmt.Sprintf("Node %d", i)
}
