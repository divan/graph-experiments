package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/divan/graph-experiments/graph"
	"github.com/divan/graph-experiments/simulation"
	"github.com/divan/graph-experiments/simulation/naivep2p"
	"github.com/divan/graph-experiments/simulation/whisperv6"
)

func main() {
	var (
		simType       = flag.String("type", "whisperv6", "Type of simulators (naivep2p, whisperv6)")
		ttl           = flag.Int("ttl", 10, "Message TTL for simulation")
		startNode     = flag.String("startNode", "192.168.1.2", "ID(name) of node initiating the message propagation")
		naiveP2PN     = flag.Int("naivep2p.N", 3, "Number of peers to propagate (0..N of peers)")
		naiveP2PDelay = flag.Duration("naivep2p.delay", 10*time.Millisecond, "Delay for each step")
		input         = flag.String("i", "network.json", "Input filename for pregenerated data to be used with simulation")
		output        = flag.String("o", "propagation.json", "Output filename for p2p sending data")
	)
	flag.Parse()

	data, err := graph.NewDataFromJSON(*input)
	if err != nil {
		log.Fatal("Opening input file failed:", err)
	}

	fd, err := os.Create(*output)
	if err != nil {
		log.Fatal("Opening output file failed:", err)
	}
	defer fd.Close()

	var sim simulation.Simulator
	switch *simType {
	case "naivep2p":
		sim = naivep2p.NewSimulator(data, *naiveP2PN, *naiveP2PDelay)
	case "whisperv6":
		sim = whisperv6.NewSimulator(data)
	}
	startNodeIdx, err := findNode(data.Nodes, *startNode)
	if err != nil {
		log.Fatalf("Can't find node '%s' in graph data. Check your -startNode option", *startNode)
	}

	// Start simulation by sending single message
	log.Printf("Starting message sending %s simulation for graph with %d nodes...", *simType, len(data.Nodes))
	sendData := sim.SendMessage(startNodeIdx, *ttl)
	err = json.NewEncoder(fd).Encode(sendData)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Written %s propagation data into %s", *simType, *output)
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
