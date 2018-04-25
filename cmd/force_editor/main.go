package main

import (
	"flag"
	"log"

	"github.com/divan/graph-experiments/graph"
	"github.com/divan/graph-experiments/layout"
)

func main() {
	flag.Parse()

	data, err := graph.NewGraphFromJSON("network.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded graph: %d nodes, %d links\n", len(data.Nodes()), len(data.Links()))

	log.Printf("Initializing layout...")
	repelling := layout.NewGravityForce(-200.0, layout.BarneHutMethod)
	springs := layout.NewSpringForce(0.01, 5.0, layout.ForEachLink)
	drag := layout.NewDragForce(0.4, layout.ForEachNode)
	layout3D := layout.NewWithDebug(data, repelling, springs, drag)

	ws := NewWSServer(layout3D)
	ws.updateGraph(data)

	log.Printf("Starting web server...")
	startWeb(ws)
	select {}
}
