package main

import (
	"flag"
	"log"
)

func main() {
	var (
		webOnly    = flag.Bool("web", false, "Don't generate data, just serve web page")
		iterations = flag.Int("iterations", 1, "Number of iterations for force-directed layout simulation")
	)
	flag.Parse()

	if !*webOnly {
		data, err := NewGraphDataFromJSON("static/example1.json")
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Loaded graph: %d nodes, %d links\n", len(data.Nodes), len(data.Links))

		layout := &Layout3D{}
		log.Printf("Initializing layout...")
		layout.Init(data)
		log.Printf("Calculating layout: nodes: %v, links: %v", layout.nodes, layout.links[0])
		_ = *iterations
		//layout.Calculate(*iterations)

		log.Printf("Writing output...")
		out := NewNgraphBinaryOutput("./static/data")
		err = out.Save(layout, data)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Fatal(startWeb())
}
