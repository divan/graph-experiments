package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/divan/graph-experiments/generation"
	"github.com/divan/graph-experiments/generation/basic"
	"github.com/divan/graph-experiments/generation/net"
)

func main() {
	var (
		genType  = flag.String("type", "net", "Generator type (net, line, circle, grid, grid3d)")
		nodes    = flag.Int("n", 20, "Number of nodes")
		netConns = flag.Int("conns", 4, "Number of connections between hosts for net generator")
		output   = flag.String("o", "network.json", "Output filename for network data")
	)
	flag.Parse()

	// Prepare output file for writing
	fd, err := os.Create(*output)
	if err != nil {
		log.Fatal("Opening output file failed:", err)
	}
	defer func(fd *os.File) {
		if err := fd.Close(); err != nil {
			log.Fatal("Closing output file failed:", err)
		}
	}(fd)

	var gen generation.GraphGenerator
	switch *genType {
	case "net":
		gen = net.NewDummyGenerator(*nodes, *netConns, "192.168.1.1", net.Exact)
	case "line":
		gen = basic.NewLineGenerator(*nodes)
	case "circle":
		gen = basic.NewCircleGenerator(*nodes)
	case "grid":
		gen = basic.NewGrid2DGeneratorN(*nodes)
	case "grid3d":
		gen = basic.NewGrid3DGeneratorN(*nodes)
	}

	log.Printf("Generating %s graph with %d nodes...\n", *genType, *nodes)
	data := gen.Generate()
	err = json.NewEncoder(fd).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Written graph into", *output)

}
