package net

import (
	"encoding/json"
	"io"
	"math/rand"
)

type Data struct {
	Nodes []*Node `json:"nodes"`
	Links []*Link `json:"links"`
}

type Node struct {
	IP string `json:"id"`
}

// NewNode generates new Node with givan IP address.
func NewNode(ip string) *Node {
	return &Node{
		IP: ip,
	}
}

// Link represents link between two nodes.
type Link struct {
	From string `json:"source"`
	To   string `json:"target"`
}

func (data *Data) linkExists(fromIP, toIP string) bool {
	for _, link := range data.Links {
		if link.From == fromIP && link.To == toIP ||
			link.To == fromIP && link.From == toIP {
			return true
		}
	}
	return false
}

// GenerateNetwork generates dummy network with known number of
// hosts with known number of connections each.
func GenerateNetwork(hosts, conn int) *Data {
	data := &Data{}
	gen := NewIPGenerator("192.168.1.1")
	for i := 0; i < hosts; i++ {
		ip := gen.NextAddress()
		node := NewNode(ip)
		data.Nodes = append(data.Nodes, node)
	}

	for i := 0; i < hosts; i++ {
		for j := 0; j < conn; j++ {
			link := &Link{
				From: data.Nodes[i].IP,
			}
			// use uniform distribution for now
			idx := rand.Intn(len(data.Nodes) - 1)
			if data.linkExists(data.Nodes[idx].IP, data.Nodes[i].IP) {
				idx = rand.Intn(len(data.Nodes) - 1)
				if data.linkExists(data.Nodes[idx].IP, data.Nodes[i].IP) {
					idx = rand.Intn(len(data.Nodes) - 1)
				}
			}

			link.To = data.Nodes[idx].IP
			data.Links = append(data.Links, link)
		}
	}

	return data
}

// DumpData serializes and dumps network graph data into the
// given writer.
func DumpData(w io.Writer, data *Data) error {
	return json.NewEncoder(w).Encode(data)
}
