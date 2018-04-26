package main

import (
	"flag"
	"log"
	"time"

	"github.com/divan/graph-experiments/graph"
	"github.com/ethereum/go-ethereum/p2p"
)

func main() {
	port := flag.String("port", "20002", "Port to bind server to")
	flag.Parse()

	hosts := []string{
		//"163.172.176.22:30503",
		//"163.172.176.22:30403",
		"51.15.85.243:30403",
		"51.15.35.110:30303",
		"51.15.85.243:30503",
	}

	ws := NewWSServer(hosts, 10*time.Second)
	ws.refresh()

	log.Printf("Starting web server...")
	startWeb(ws, *port)
	select {}
}

func AddPeer(g *graph.Graph, fromID string, to *p2p.PeerInfo) {
	toID := to.ID
	addNode(g, fromID, false)
	addNode(g, toID, isClient(to.Name))

	if g.LinkExistsByID(fromID, toID) {
		return
	}
	if to.Network.Inbound == false {
		g.AddLinkByIDs(fromID, toID)
	} else {
		g.AddLinkByIDs(toID, fromID)
	}
}

func addNode(g *graph.Graph, id string, client bool) {
	if _, err := g.NodeByID(id); err == nil {
		// already exists
		return
	}
	node := NewNode(id, client)
	g.AddNode(node)
}
