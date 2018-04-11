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

	plog, err := LoadPropagationData("propagation.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded propagation data: %d timestamps\n", len(plog.Timestamps))

	log.Printf("Initializing layout...")
	repelling := layout.NewGravityForce(-80.0, layout.BarneHutMethod)
	springs := layout.NewSpringForce(0.02, 5.0, layout.ForEachLink)
	drag := layout.NewDragForce(0.7, layout.ForEachNode)
	layout3D := layout.New(data, repelling, springs, drag)

	ws := NewWSServer(layout3D)
	ws.layout.CalculateN(600)
	ws.updateGraph(data)
	ws.updatePropagationData(plog)

	log.Printf("Starting web server...")
	startWeb(ws)
	select {}
}
