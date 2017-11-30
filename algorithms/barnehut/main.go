package main

import (
	"flag"
	"log"
)

func main() {
	var webOnly = flag.Bool("web", false, "Don't generate data, just serve web page")
	flag.Parse()

	if !*webOnly {
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
	}

	log.Fatal(startWeb())
}
