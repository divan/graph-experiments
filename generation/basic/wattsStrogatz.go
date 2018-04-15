package basic

import (
	"math"
	"math/rand"

	"github.com/divan/graph-experiments/graph"
)

// WattsStrogatzGenerator implements generator for Watts-Strogatz graph.
type WattsStrogatzGenerator struct {
	nodes              int // number of nodes
	conns              int // number of neigbours
	rewritePropability float64
}

// NewWattsStrogatzGenerator creates new Watts-Strogatz generator for N nodes graph.
func NewWattsStrogatzGenerator(n, conns int) *WattsStrogatzGenerator {
	if conns > n {
		panic("conns should be less then number of nodes")
	}
	return &WattsStrogatzGenerator{
		nodes:              n,
		conns:              conns,
		rewritePropability: 0.01,
	}
}

// Generate generates the data for graph. Implements Generator interface.
func (l *WattsStrogatzGenerator) Generate() *graph.Data {
	data := graph.NewData()

	for i := 0; i < l.nodes; i++ {
		addNode(data, i)

	}

	// connect each node conns/2 neigbors
	neigbors := int(math.Floor(float64(l.conns/2 + 1)))
	for i := 1; i < neigbors; i++ {
		for j := 0; j < l.nodes; j++ {
			to := int(math.Mod(float64(i+j), float64(l.nodes)))
			addLink(data, j, to)
		}
	}

	// rewire edges from each node
	neigbors = int(math.Floor(float64(l.conns/2 + 1)))
	for j := 1; j < neigbors; j++ {
		for i := 0; i < l.nodes; i++ {
			if rand.Float64() > l.rewritePropability {
				continue
			}

			from := i
			to := int(math.Mod(float64(i+j), float64(l.nodes)))
			newTo := rand.Intn(l.nodes)

			// TODO: switch to link indexes
			needsRewire := (newTo == i) || data.LinkExists(data.Nodes[from].ID, data.Nodes[newTo].ID)
			if needsRewire && (data.NodeLinks(data.Nodes[from].ID) == l.nodes-1) {
				continue
			}

			for needsRewire {
				newTo = rand.Intn(l.nodes)
				needsRewire = (newTo == i) || data.LinkExists(data.Nodes[from].ID, data.Nodes[newTo].ID)
			}

			rewireLink(data, from, to, newTo)
		}
	}

	return data
}

// TODO: move it into graph package
func rewireLink(d *graph.Data, from, to, newTo int) {
	source := d.Nodes[from].ID
	target := d.Nodes[to].ID
	newTarget := d.Nodes[newTo].ID
	for i := range d.Links {
		if d.Links[i].Source == source && d.Links[i].Target == target {
			d.Links[i].Target = newTarget
		} else if d.Links[i].Target == source && d.Links[i].Source == target {
			d.Links[i].Source = newTarget
		}
	}
}
