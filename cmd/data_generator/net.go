package main

import (
	//"github.com/divan/graph-experiments/graph"

	"encoding/json"
	"io"
	"math/rand"
)

type NetworkData struct {
	Nodes []*NetNode `json:"nodes"`
	Links []*NetLink `json:"links"`
}

type NetNode struct {
	IP string `json:"id"`
}

// NewNetNode generates new NetNode with givan IP address.
func NewNetNode(ip string) *NetNode {
	return &NetNode{
		IP: ip,
	}
}

type NetLink struct {
	From string `json:"source"`
	To   string `json:"target"`
}

func (data *NetworkData) linkExists(fromIP, toIP string) bool {
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
func GenerateNetwork(hosts, conn int) *NetworkData {
	data := &NetworkData{}
	gen := NewIPGenerator("192.168.1.1")
	for i := 0; i < hosts; i++ {
		ip := gen.NextAddress()
		node := NewNetNode(ip)
		data.Nodes = append(data.Nodes, node)
	}

	for i := 0; i < hosts; i++ {
		for j := 0; j < conn; j++ {
			link := &NetLink{
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

// DumpNetworkData serializes and dumps network graph data into the
// given writer.
func DumpNetworkData(w io.Writer, data *NetworkData) error {
	return json.NewEncoder(w).Encode(data)
}
