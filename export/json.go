package export

import (
	"encoding/json"
	"io"

	"github.com/divan/graph-experiments/graph"
)

// JSON implements GraphExporter into JSON format.
type JSON struct {
	output   io.Writer
	indented bool
}

func NewJSON(w io.Writer, indented bool) *JSON {
	return &JSON{
		output:   w,
		indented: indented,
	}
}

func (j *JSON) ExportGraph(g *graph.Graph) error {
	type link struct {
		Source string `json:"source"`
		Target string `json:"target"`
	}
	var data struct {
		Nodes []graph.Node `json:"nodes"`
		Links []*link      `json:"links"`
	}

	nodes := make(map[int]string)
	for i := range g.Nodes() {
		nodes[i] = g.Nodes()[i].ID()
	}

	data.Nodes = g.Nodes()
	data.Links = make([]*link, len(g.Links()))
	for i, l := range g.Links() {
		data.Links[i] = &link{
			Source: nodes[l.From],
			Target: nodes[l.To],
		}
	}

	enc := json.NewEncoder(j.output)
	if j.indented {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(data)
}
