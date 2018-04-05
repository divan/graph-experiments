package p2p

import "github.com/divan/graph-experiments/cmd/data_generator/net"

// LinkIndex stores link information in form of indexes, rather than nodes IP.
type LinkIndex struct {
	From int
	To   int
}

// PrecalculateLinkIndexes prepares slice of LinkIndex for faster lookup.
// TODO: move this indexes stuff into Data structure itself, so this can be removed.
func PrecalculateLinkIndexes(data *net.Data) []LinkIndex {
	m := make(map[string]int)
	for idx := range data.Nodes {
		m[data.Nodes[idx].IP] = idx
	}

	ret := make([]LinkIndex, len(data.Links))
	for i, link := range data.Links {
		ret[i] = LinkIndex{
			From: m[link.From],
			To:   m[link.To],
		}
	}
	return ret
}

// PrecalculatePeers creates map with peers indexes for faster lookup.
func PrecalculatePeers(data *net.Data) map[int][]int {
	links := PrecalculateLinkIndexes(data)

	ret := make(map[int][]int)
	for _, link := range links {
		if _, ok := ret[link.From]; !ok {
			ret[link.From] = make([]int, 0)
		}

		peers := ret[link.From]
		peers = append(peers, link.To)
		ret[link.From] = peers
	}
	return ret
}
