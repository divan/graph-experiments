package net

import (
	"errors"
	"math/rand"

	"github.com/divan/graph-experiments/graph"
)

// DummyGenerator implements dummy network generator,
// where network consits from N hosts, with M connections
// each. Nodes are represented as IPv4 addresses.
type DummyGenerator struct {
	hosts        int
	connections  int
	startIP      string
	distribution ConnectionsDistribution
}

// NewDummyGenerator creates new dummy network generator with given parameters.
func NewDummyGenerator(hosts, conns int, startIP string, distribution ConnectionsDistribution) *DummyGenerator {
	return &DummyGenerator{
		hosts:        hosts,
		connections:  conns,
		startIP:      startIP,
		distribution: distribution,
	}
}

// ConnectionsDistribution represents types of distributions of
// connections between nodes.
type ConnectionsDistribution int

// Predefined types of connections distributions.
const (
	Exact ConnectionsDistribution = iota
	Uniform
)

type Node struct {
	IP string `json:"id"`
}

// NewNode generates new Node with givan IP address.
func NewNode(ip string) *Node {
	return &Node{
		IP: ip,
	}
}

// ID implements graph.Node interface and returns node ID.
func (n *Node) ID() string {
	return n.IP
}

// Generate generates dummy network with known number of
// hosts with known number of connections each. Implements Generator
// interface.
func (d *DummyGenerator) Generate() *graph.Graph {
	g := graph.NewGraph()

	// generate hosts
	gen := NewIPGenerator(d.startIP)
	for i := 0; i < d.hosts; i++ {
		ip := gen.NextAddress()
		node := NewNode(ip)
		g.AddNode(node)
	}

	// generate links
	for i := 0; i < d.hosts; i++ {
		for j := 0; j < d.conns(); j++ {
			idx, err := d.nextIdx(g, i)
			if err == nil {
				g.AddLink(i, idx)
			}
		}
	}

	return g
}

// conns return the number of connections based on the
// actual distrubtion.
func (d *DummyGenerator) conns() int {
	switch d.distribution {
	case Uniform:
		n := rand.Intn(d.connections)
		if n == 0 {
			n = 1
		}
		return n
	case Exact:
		return d.connections
	}

	return d.connections
}

// nextIdx returns next node idx to connect to.
// For now it uses uniform distribution and retries two times
// to minimize the probability of choosing the existing link
// (it doesn't guarantee, but it's cheap).
// TODO: make it more efficient and faster
func (d *DummyGenerator) nextIdx(g *graph.Graph, i int) (int, error) {
	// use uniform distribution for now
	idx := rand.Intn(d.hosts - 1)
	if idx == i || g.LinkExists(idx, i) {
		idx = rand.Intn(d.hosts - 1)
		if idx == i || g.LinkExists(idx, i) {
			idx = rand.Intn(d.hosts - 1)
			if idx == i || g.LinkExists(idx, i) {
				return 0, errors.New("too many colissions")
			}
		}
	}
	return idx, nil
}
