package basic

import (
	"fmt"

	"github.com/divan/graph-experiments/graph"
)

// LineGenerator implements generator for line graph.
type LineGenerator struct {
	nodes int // number of nodes
}

// NewLineGenerator creates new line generator for N nodes graph.
func NewLineGenerator(n int) *LineGenerator {
	return &LineGenerator{
		nodes: n,
	}
}

// Generate generates the data for graph. Implements Generator interface.
func (l *LineGenerator) Generate() *graph.Data {
	data := graph.NewData()

	for i := 0; i < l.nodes; i++ {
		node := &graph.Node{
			ID: l.idxToID(i),
		}
		data.Nodes = append(data.Nodes, node)

		if i != l.nodes-1 {
			link := &graph.Link{
				Source: l.idxToID(i),
				Target: l.idxToID(i + 1),
			}
			data.Links = append(data.Links, link)
		}
	}

	return data
}

func (l *LineGenerator) idxToID(i int) string {
	return fmt.Sprintf("Node %d", i)
}
