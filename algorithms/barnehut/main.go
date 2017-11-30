package main

import (
	"log"
)

func main() {
	data, err := NewGraphDataFromJSON("static/data.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Loaded graph: %d nodes, %d links\n", len(data.Nodes), len(data.Links))

	layout := &Layout3D{}
	layout.InitCoordinates(data.Nodes)

	out := NewNgraphBinaryOutput("./static/data")
	err = out.Save(layout, data)
	if err != nil {
		log.Fatal(err)
	}

	err = startWeb()
	if err != nil {
		log.Fatal(err)
	}
}
