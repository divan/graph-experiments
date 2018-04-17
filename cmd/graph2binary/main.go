package main

import (
	"flag"
	"log"
	"os"

	"github.com/divan/graph-experiments/export/ngraph_binary"
	"github.com/divan/graph-experiments/graph"
	"github.com/divan/graph-experiments/layout"
)

func main() {
	iterations := flag.Int("iterations", 1, "Number of iterations for force-directed layout simulation")
	dir := flag.String("dir", "data/", "Output dir")
	flag.Parse()

	file := "data.json"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	data, err := graph.NewGraphFromJSON(file)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded graph: %d nodes, %d links\n", len(data.Nodes()), len(data.Links()))
	ngraph := ngraph_binary.NewExporter(*dir)

	log.Printf("Initializing layout...")
	layout := layout.New(data)
	log.Printf("Calculating layout...")
	layout.CalculateN(*iterations)

	log.Printf("Writing output to %s...\n", *dir)
	err = ngraph.Save(layout, data)
	if err != nil {
		log.Fatal(err)
	}
}
