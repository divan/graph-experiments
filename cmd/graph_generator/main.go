package main

import (
	"flag"
	"log"
	"os"

	"github.com/divan/graph-experiments/export"
	"github.com/divan/graph-experiments/generation"
	"github.com/divan/graph-experiments/generation/basic"
	"github.com/divan/graph-experiments/generation/net"
)

func main() {
	var (
		genType = flag.String("type", "net", "Generator type (net, split-brain, line, circle, grid, grid3d, small-world, king)")
		nodes   = flag.Int("n", 20, "Number of nodes")
		conns   = flag.Int("conns", 4, "Number of connections between hosts for net generator")
		output  = flag.String("o", "network.json", "Output filename for network data")
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
		gen = net.NewNetGenerator(*nodes, *conns, "192.168.1.1", net.Exact)
	case "split-brain":
		gen = net.NewSplitBrainGenerator(*nodes, *conns, "192.168.1.1")
	case "line":
		gen = basic.NewLineGenerator(*nodes)
	case "circle":
		gen = basic.NewCircleGenerator(*nodes)
	case "grid":
		gen = basic.NewGrid2DGeneratorN(*nodes)
	case "grid3d":
		gen = basic.NewGrid3DGeneratorN(*nodes)
	case "small-world":
		gen = basic.NewWattsStrogatzGenerator(*nodes, *conns)
	case "king":
		gen = basic.NewKingGeneratorN(*nodes)
	}

	log.Printf("Generating %s graph with %d nodes...\n", *genType, *nodes)
	data := gen.Generate()
	err = export.NewJSON(fd, true).ExportGraph(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Written graph into", *output)
}
