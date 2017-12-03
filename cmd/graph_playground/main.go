package main

import (
	"flag"
	"github.com/divan/graph-experiments/export/ngraph_binary"
	"github.com/divan/graph-experiments/graph"
	"github.com/divan/graph-experiments/layout"
	"log"
)

func main() {
	iterations := flag.Int("steps", 100, "Number of iterations for force-directed layout simulation")
	flag.Parse()

	data, err := graph.NewDataFromJSON("static/data.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded graph: %d nodes, %d links\n", len(data.Nodes), len(data.Links))

	log.Printf("Initializing layout...")
	layout := layout.New(data)
	log.Printf("Calculating layout...")
	layout.Calculate(*iterations)

	log.Printf("Writing output...")
	ngraph := ngraph_binary.NewExporter("./static/data")
	err = ngraph.Save(layout, data)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(startWeb())
}
