package main

import (
	"fmt"

	"github.com/divan/graph-experiments/graph"
)

func printStats(g *graph.Graph) {
	peers, clients := make(map[string]int), make(map[string]int)
	for _, node := range g.Nodes() {
		n := node.(*Node)
		if n.IsClient() {
			clients[n.ID()]++
		} else {
			peers[n.ID()]++
		}
	}

	fmt.Println("Graph stats:")
	fmt.Printf("Nodes - %d (%d clients, %d servers)\n", len(g.Nodes()), len(clients), len(peers))
	fmt.Println("Links:", len(g.Links()))

	findLinks := func(idx int) []*graph.Link {
		var ret []*graph.Link
		for _, link := range g.Links() {
			if link.From == idx || link.To == idx {
				ret = append(ret, link)
			}
		}
		return ret
	}
	for i, node := range g.Nodes() {
		n := node.(*Node)

		var peers, clients int
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
			} else {
				peers++
			}
		}
		if n.IsClient() {
			fmt.Printf(" Client ")
		} else {
			fmt.Printf(" Peer ")
		}
		fmt.Printf("%s - %d (%d clients, %d servers)\n", node.ID(), len(links), clients, peers)
	}

}
