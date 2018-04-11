package whisperv6

import (
	"fmt"
	"log"

	"github.com/divan/graph-experiments/graph"
	"github.com/divan/graph-experiments/simulation"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
)

// Simulator simulates WhisperV6 message propagation through the
// given p2p network.
type Simulator struct {
	data    *graph.Data
	network *simulations.Network
	nodes   []*simulations.Node
}

// NewSimulator intializes simulator for the given graph data.
func NewSimulator(data *graph.Data) *Simulator {
	services := map[string]adapters.ServiceFunc{
		"ping-pong": func(ctx *adapters.ServiceContext) (node.Service, error) {
			return newPingPongService(ctx.Config.ID), nil
		},
	}
	adapters.RegisterServices(services)

	adapter := adapters.NewSimAdapter(services)
	network := simulations.NewNetwork(adapter, &simulations.NetworkConfig{
		DefaultService: "ping-pong",
	})

	nodeCount := len(data.Nodes)
	sim := &Simulator{
		data:    data,
		network: network,
		nodes:   make([]*simulations.Node, nodeCount),
	}

	log.Println("Creating nodes...")
	for i := 0; i < nodeCount; i++ {
		node, err := sim.network.NewNodeWithConfig(nodeConfig(i))
		if err != nil {
			panic(err)
		}
		sim.nodes[i] = node
	}

	log.Println("Starting nodes...")
	for i := 0; i < nodeCount; i++ {
		err := network.Start(sim.nodes[i].ID())
		if err != nil {
			panic(err)
		}
	}

	log.Println("Connecting nodes...")
	for _, link := range data.Links {
		fromIdx, err := findNode(data.Nodes, link.Source)
		if err != nil {
			panic(err)
		}
		toIdx, err := findNode(data.Nodes, link.Target)
		if err != nil {
			panic(err)
		}
		node1 := sim.nodes[fromIdx]
		node2 := sim.nodes[toIdx]
		if err := network.Connect(node1.ID(), node2.ID()); err != nil {
			panic(err)
		}
	}

	return sim
}

// Stop stops simulator and frees all resources if any.
func (s *Simulator) Stop() error {
	log.Println("Shutting down simulation nodes...")
	return s.network.StopAll()
}

// SendMessage sends single message and tracks propagation. Implements simulator.Interface.
func (s *Simulator) SendMessage(startNodeIdx, ttl int) *simulation.Log {
	return &simulation.Log{}
}

// nodeConfig generates config for simulated node with random key.
func nodeConfig(idx int) *adapters.NodeConfig {
	key, err := crypto.GenerateKey()
	if err != nil {
		panic("unable to generate key")
	}
	var id discover.NodeID
	pubkey := crypto.FromECDSAPub(&key.PublicKey)
	copy(id[:], pubkey[1:])
	return &adapters.NodeConfig{
		ID:         id,
		PrivateKey: key,
		Name:       nodeIdxToName(idx),
	}
}

func nodeIdxToName(id int) string {
	return fmt.Sprintf("Node %d", id)
}

// findNode is a helper for finding node index by it's ID.
// TODO: remove this when links replaces into indexes.
func findNode(nodes []*graph.Node, ID string) (int, error) {
	for i := range nodes {
		if nodes[i].ID == ID {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Node with ID '%s' not found", ID)
}
