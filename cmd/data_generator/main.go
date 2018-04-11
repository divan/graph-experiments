package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/divan/graph-experiments/generator/basic"
	"github.com/divan/graph-experiments/generator/net"
	"github.com/divan/graph-experiments/graph"
	"github.com/divan/graph-experiments/simulation"
	"github.com/divan/graph-experiments/simulation/naivep2p"
	"github.com/divan/graph-experiments/simulation/whisperv6"
)

type Generator interface {
	Generate() *graph.Data
}

func main() {
	var (
		dataKind       = flag.String("type", "net", "Example random IPs network (net, line, circle, grid, grid3d)")
		propagate      = flag.Bool("simulate", false, "Run propagation simulation")
		simulator      = flag.String("sim.type", "whisperv6", "Types of simulators (naivep2p, whisperv6)")
		nodes          = flag.Int("n", 20, "Number of nodes")
		netConns       = flag.Int("net.connections", 4, "Number of connections between hosts for net generator")
		naiveP2P_N     = flag.Int("naivep2p.N", 3, "Number of peers to propagate (0..N of peers)")
		naiveP2P_Delay = flag.Duration("naivep2p.delay", 10*time.Millisecond, "Delay for each step")
		simTTL         = flag.Int("sim.ttl", 10, "Message TTL for simulation")
		simStartNode   = flag.String("sim.startNode", "192.168.1.2", "IP address of node initiating the sending")
		netOutput      = flag.String("o", "data.json", "Output filename for network data")
		simOutput      = flag.String("p", "propagation.json", "Output filename for p2p sending data")
	)
	flag.Parse()

	// Prepare output files for writing
	netFd, err := os.Create(*netOutput)
	if err != nil {
		log.Fatal("Open file for writing failed:", err)
	}
	defer netFd.Close()

	simFd, err := os.Create(*simOutput)
	if err != nil {
		log.Fatal("Open file for writing failed:", err)
	}
	defer simFd.Close()

	var generator Generator
	switch *dataKind {
	case "net":
		generator = net.NewDummyGenerator(*nodes, *netConns, "192.168.1.1", net.Exact)
	case "line":
		generator = basic.NewLineGenerator(*nodes)
	case "circle":
		generator = basic.NewCircleGenerator(*nodes)
	case "grid":
		generator = basic.NewGrid2DGeneratorN(*nodes)
	case "grid3d":
		generator = basic.NewGrid3DGeneratorN(*nodes)
	}

	data := generator.Generate()
	err = json.NewEncoder(netFd).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Written graph into", *netOutput)

	if *propagate {
		var sim simulation.Simulator
		switch *simulator {
		case "naivep2p":
			sim = naivep2p.NewSimulator(data, *naiveP2P_N, *naiveP2P_Delay)
		case "whisperv6":
			sim = whisperv6.NewSimulator(data)
		}
		startNodeIdx, err := findNode(data.Nodes, *simStartNode)
		if err != nil {
			log.Fatalf("Can't find node '%s' in graph data. Check your -sim.startNode option", *simStartNode)
		}

		// Start simulation by sending single message
		sendData := sim.SendMessage(startNodeIdx, *simTTL)
		err = json.NewEncoder(simFd).Encode(sendData)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Written p2p sending data into", *netOutput)
	}
}

// findNode is a helper for finding node index by it's ID.
func findNode(nodes []*graph.Node, ID string) (int, error) {
	for i := range nodes {
		if nodes[i].ID == ID {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Node with ID '%s' not found", ID)
}
