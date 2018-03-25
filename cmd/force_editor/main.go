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
	repelling := layout.NewGravityForce(-100.2, layout.BarneHutMethod)
	springs := layout.NewSpringForce(0.011, 20.0, layout.BarneHutMethod)
	layout3D := layout.New(data, repelling, springs)

	ws := NewWSServer(layout3D)
	ws.updateGraph(data)

	log.Printf("Starting web server...")
	startWeb(ws)
	select {}
}
