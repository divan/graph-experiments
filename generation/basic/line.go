package basic

import (
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
		addNode(data, i)

		if i != l.nodes-1 {
			addLink(data, i, i+1)
		}
	}

	return data
}
