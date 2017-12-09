package main

import (
	"flag"
	"log"

	"github.com/divan/graph-experiments/graph"
	"github.com/divan/graph-experiments/layout"
)

func main() {
	flag.Parse()

	data, err := graph.NewDataFromJSON("data.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded graph: %d nodes, %d links\n", len(data.Nodes), len(data.Links))

	log.Printf("Initializing layout...")
	layout := layout.New(data)

	ws := NewWSServer(layout)
	ws.updateGraph(data)

	startWeb(ws)
	select {}
}
