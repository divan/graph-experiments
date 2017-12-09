package main

import (
	"flag"
	"log"
	"time"

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

	ws := NewWSServer()
	startWeb(ws)
	ws.updateGraph(data)
	for i := 0; i < 10; i++ {
		log.Printf("Calculating layout...")
		layout.Calculate(1)
		nodes := layout.Nodes()
		ws.updatePositions(nodes)

		time.Sleep(1 * time.Second)
	}
	select {}
}
