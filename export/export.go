package export

import (
	"github.com/divan/graph-experiments/graph"
	//"github.com/divan/graph-experiments/layout"
)

// GraphExporter defines exporting capability for graph.
type GraphExporter interface {
	ExportGraph(*graph.Graph) error
}

/*
// GraphExporter defines exporting capability for layouts.
type LayoutExporter interface {
	ExportLayout(*layout.Layout) error
}
*/
