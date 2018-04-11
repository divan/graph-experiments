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
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"
)

// Simulator simulates WhisperV6 message propagation through the
// given p2p network.
type Simulator struct {
	data     *graph.Data
	network  *simulations.Network
	nodes    []*simulations.Node
	whispers map[discover.NodeID]*whisper.Whisper
}

// NewSimulator intializes simulator for the given graph data.
func NewSimulator(data *graph.Data) *Simulator {
	whispers := make(map[discover.NodeID]*whisper.Whisper)

	cfg := &whisper.Config{
		MaxMessageSize:     whisper.DefaultMaxMessageSize,
		MinimumAcceptedPOW: 0.001,
	}
	services := map[string]adapters.ServiceFunc{
		"shh": func(ctx *adapters.ServiceContext) (node.Service, error) {
			// it's important to init whisper service here, as it
			// be initialized for each peer
			id := ctx.Config.ID
			service := whisper.New(cfg)
			whispers[id] = service
			return service, nil
		},
	}
	adapters.RegisterServices(services)

	adapter := adapters.NewSimAdapter(services)
	network := simulations.NewNetwork(adapter, &simulations.NetworkConfig{
		DefaultService: "shh",
	})

	nodeCount := len(data.Nodes)
	sim := &Simulator{
		data:     data,
		network:  network,
		nodes:    make([]*simulations.Node, nodeCount),
		whispers: whispers,
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
	if err := network.StartAll(); err != nil {
		panic(err)
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
		// if connection already exists, skip it, as network.Connect will fail
		if network.GetConn(node1.ID(), node2.ID()) != nil {
			continue
		}
		if err := network.Connect(node1.ID(), node2.ID()); err != nil {
			panic(err)
		}
	}

	return sim
}

// Stop stops simulator and frees all resources if any.
func (s *Simulator) Stop() error {
	log.Println("Shutting down simulation nodes...")
	s.network.Shutdown()
	return nil
}

// SendMessage sends single message and tracks propagation. Implements simulator.Interface.
func (s *Simulator) SendMessage(startNodeIdx, ttl int) *simulation.Log {
	node := s.nodes[startNodeIdx]
	service, ok := s.whispers[node.ID()]
	if !ok {
		log.Fatalf("Whisper service for node %d not found", startNodeIdx)
	}

	log.Println(" Sending Whisper message...")
	msg := generateMessage(ttl)
	err := service.Send(msg)
	log.Println(" Error:", err)

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
