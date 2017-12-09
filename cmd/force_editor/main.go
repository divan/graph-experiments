package main

import (
	"flag"
	"log"

	"github.com/divan/graph-experiments/graph"
)

func main() {
	flag.Parse()

	data, err := graph.NewDataFromJSON("web/data/data.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded graph: %d nodes, %d links\n", len(data.Nodes), len(data.Links))

	/*
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
	*/

	log.Fatal(startWeb())
}
