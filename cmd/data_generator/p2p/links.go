package p2p

import (
	"github.com/divan/graph-experiments/graph"
)

// LinkIndex stores link information in form of indexes, rather than nodes IP.
type LinkIndex struct {
	From int
	To   int
}

// PrecalculateLinkIndexes prepares slice of LinkIndex for faster lookup.
// TODO: move this indexes stuff into graph.Data structure itself, so this can be removed.
func PrecalculateLinkIndexes(data *graph.Data) []LinkIndex {
	m := make(map[string]int)
	for idx := range data.Nodes {
		m[data.Nodes[idx].ID] = idx
	}

	ret := make([]LinkIndex, len(data.Links))
	for i, link := range data.Links {
		ret[i] = LinkIndex{
			From: m[link.Source],
			To:   m[link.Target],
		}
	}
	return ret
}

// PrecalculatePeers creates map with peers indexes for faster lookup.
func PrecalculatePeers(data *graph.Data) map[int][]int {
	links := PrecalculateLinkIndexes(data)

	ret := make(map[int][]int)
	for _, link := range links {
		if link.From == link.To {
			continue
		}
		if _, ok := ret[link.From]; !ok {
			ret[link.From] = make([]int, 0)
		}

		peers := ret[link.From]
		peers = append(peers, link.To)
		ret[link.From] = peers
	}
	return ret
}
