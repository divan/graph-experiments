package main

import (
	"sync"
	"time"

	"github.com/divan/graph-experiments/graph"
)

type Stats struct {
	mu sync.RWMutex

	// current data
	Clients    []string
	ClientsNum int
	Servers    []string
	ServersNum int
	LinksNum   int
	Nodes      []*NodeStats
	LastUpdate time.Time

	// history data
	Timestamps  []string
	ServersHist []int
	ClientsHist []int
}

type NodeStats struct {
	ID         string
	Peers      []string
	Clients    []string
	PeersNum   int
	ClientsNum int
	IsClient   bool
}

func (s *Stats) Stats() *Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s
}

func (s *Stats) Update(g *graph.Graph) {
	var servers, clients []string
	for _, node := range g.Nodes() {
		n := node.(*Node)
		if n.IsClient() {
			clients = append(clients, n.ID())
		} else {
			servers = append(servers, n.ID())
		}
	}

	findLinks := func(idx int) []*graph.Link {
		var ret []*graph.Link
		for _, link := range g.Links() {
			if link.From == idx || link.To == idx {
				ret = append(ret, link)
			}
		}
		return ret
	}

	ns := make([]*NodeStats, 0, len(g.Nodes()))
	for i, node := range g.Nodes() {
		n := node.(*Node)

		var peers, clients int
		var peersS, clientsS []string
		links := findLinks(i)
		for _, link := range links {
			var peer graph.Node
			if i == link.From {
				peer = g.Nodes()[link.To]
			} else {
				peer = g.Nodes()[link.From]
			}
			p := peer.(*Node)
			if p.IsClient() {
				clients++
				clientsS = append(clientsS, p.ID())
			} else {
				peers++
				peersS = append(peersS, p.ID())
			}
		}

		nodeStat := &NodeStats{
			ID:         n.ID(),
			IsClient:   n.IsClient(),
			Peers:      peersS,
			Clients:    clientsS,
			PeersNum:   peers,
			ClientsNum: clients,
		}
		ns = append(ns, nodeStat)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Clients = clients
	s.ClientsNum = len(clients)
	s.Servers = servers
	s.ServersNum = len(servers)
	s.LinksNum = len(g.Links())

	s.Nodes = ns

	now := time.Now()
	s.LastUpdate = now

	// update historic data
	JavascriptISOString := "2006-01-02T15:04:05.999Z07:00"
	jsNow := now.UTC().Format(JavascriptISOString)
	s.Timestamps = append(s.Timestamps, jsNow)
	s.ServersHist = append(s.ServersHist, s.ServersNum)
	s.ClientsHist = append(s.ClientsHist, s.ClientsNum)
}
