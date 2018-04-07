package p2p

import (
	"fmt"
	"log"
	"time"

	"github.com/divan/graph-experiments/cmd/data_generator/net"
)

// PropagationLog represnts log of p2p message propagation
// with relative timestamps (starting from T0).
type PropagationLog struct {
	Timestamps []int   // timestamps in milliseconds starting from T0
	Indices    [][]int // indices of links for each step, len should be equal to len of Timestamps field
	Nodes      [][]int // indices of nodes involved in each step
}

// SimulatePropagation start simulation of message propagation through the network
// defined by data, starting from node startNode, and sending message to N peers
// each delay.
// Simulation assumes that links are all of equal length and propagation takes
// fixed amount of time for now.
func SimulatePropagation(data *net.Data, N, ttl int, delay time.Duration, startNodeIP string) *PropagationLog {
	startIdx, err := findNode(data.Nodes, startNodeIP)
	if err != nil {
		log.Fatal(err)
	}
	s := NewSimulator(data, N, delay)
	logEntries := s.Run(startIdx, ttl)
	return s.LogEntries2PropagationLog(logEntries)
}

// findNode is a helper for finding node index by it's IP address.
func findNode(nodes []*net.Node, IP string) (int, error) {
	for i := range nodes {
		if nodes[i].IP == IP {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Node with IP '%s' not found", IP)
}
