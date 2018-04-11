package whisperv6

import (
	"log"
	"time"

	"github.com/divan/graph-experiments/graph"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
)

// Simulator simulates WhisperV6 message propagation through the
// given p2p network.
type Simulator struct {
	data *graph.Data
}

// NewSimulator intializes simulator for the given graph data.
func NewSimulator(data *graph.Data) *Simulator {
	return &Simulator{
		data: data,
	}
}

// SendMessage sends single message and tracks propagation. Implements simulator.Interface.
func (s *Simulator) SendMessage(startNodeIdx, ttl int) []*LogEntry {
	services := map[string]adapters.ServiceFunc{
		"ping-pong": func(ctx *adapters.ServiceContext) (node.Service, error) {
			return newPingPongService(ctx.Config.ID), nil
		},
	}
	adapters.RegisterServices(services)

	var adapter adapters.NodeAdapter
	adapter = adapters.NewSimAdapter(services)

	network := simulations.NewNetwork(adapter, &simulations.NetworkConfig{
		DefaultService: "ping-pong",
	})

	log.Println("Creating nodes...")
	node1, err := network.NewNodeWithConfig(nodeConfig("Node 1"))
	if err != nil {
		panic(err)
	}
	node2, err := network.NewNodeWithConfig(nodeConfig("Node 2"))
	if err != nil {
		panic(err)
	}

	log.Println("Starting nodes...")
	if err := network.Start(node1.ID()); err != nil {
		panic(err)
	}
	if err := network.Start(node2.ID()); err != nil {
		panic(err)
	}
	defer network.StopAll()

	log.Println("Connecting nodes...")
	if err := network.Connect(node1.ID(), node2.ID()); err != nil {
		panic(err)
	}

	log.Println("Sleeping 10 secs...")
	time.Sleep(10 * time.Second)
}

// nodeConfig generates config for simulated node with random key.
func nodeConfig(name string) *adapters.NodeConfig {
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
		Name:       name,
	}
}
