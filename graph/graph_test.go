package graph

import (
	"bytes"
	"io"
	"testing"
)

func TestNewJSONGraph(t *testing.T) {
	buf := bytes.NewBufferString(`{
		"nodes": [ {"id": "A", "weight": 10}, {"id": "B"}, {"id": "C"}, {"id": "D"} ],
		"links": [ {"source": "A", "target": "B"}, {"source": "C", "target": "D"}, {"source": "C", "target": "A"}]
	}`)
	graph, err := NewGraphFromJSONReader(buf)
	if err != nil {
		t.Fatal(err)
	}

	nodes := graph.Nodes()
	if len(nodes) != 4 {
		t.Fatalf("Expect graph to have %d nodes, but got %d", 4, len(nodes))
	}

	links := graph.Links()
	if len(links) != 3 {
		t.Fatalf("Expect graph to have %d links, but got %d", 3, len(links))
	}

	linksCounter := map[string]int{
		"A": 2,
		"B": 1,
		"C": 2,
		"D": 1,
	}
	for i, node := range nodes {
		got := graph.NodeLinks(i)
		expected := linksCounter[node.ID()]
		if got != expected {
			t.Fatalf("Expected number of links to be %d, but got %d for node '%s'",
				expected, got, node.ID())
		}
	}
}

func TestNewJSONGraphLarge(t *testing.T) {
	N := int(10e5)
	r := generateLargeGraphJSON(N, 2*N)

	graph, err := NewGraphFromJSONReader(r)
	if err != nil {
		t.Fatal(err)
	}

	nodes := graph.Nodes()
	if len(nodes) != N {
		t.Fatalf("Expect graph to have %v nodes, but got %d", N, len(nodes))
	}

	links := graph.Links()
	if len(links) != 2*N {
		t.Fatalf("Expect graph to have %v links, but got %d", 2*N, len(links))
	}
}

func BenchmarkImportJSONGraph(b *testing.B) {
	N := int(10e4)
	r := generateLargeGraphJSON(N, 2*N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		NewGraphFromJSONReader(r)
	}
}

// generate large graph with given number of nodes and links.
func generateLargeGraphJSON(nodes, links int) io.Reader {
	buf := new(bytes.Buffer)
	buf.WriteString(`{"nodes":[`)
	for i := 0; i < nodes; i++ {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(`{"id": "same"}`)
	}
	buf.WriteString(`],"links":[`)
	for i := 0; i < links; i++ {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(`{"source": "same", "target": "same"}`)
	}
	buf.WriteString(`]}`)
	return buf
}
