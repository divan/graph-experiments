package generation

import "github.com/divan/graph-experiments/graph"

// GraphGenerator represents generator that generates
// graph data with nodes and links.
type GraphGenerator interface {
	Generate() *graph.Data
}
