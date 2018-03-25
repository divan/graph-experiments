package layout

import (
	"fmt"

	"github.com/divan/graph-experiments/graph"
)

// ForceRule define algorithm/rules to apply force on a graph. Force can be applied an a variety of different ways and this abstraction should ideally catch and encapsulate all these differences.
//
// vectors and debugInfo are passed for optimization purposes, to avoid allocating new memory.
type ForceRule func(
	force Force,
	nodes []*Node,
	links []*graph.LinkData,
	vectors map[int]*ForceVector,
	debugInfo ForcesDebugData)

// ForEachLink applies force to both ends of each link in the graph, with positive and negative signs respectively.
var ForEachLink = func(
	force Force,
	nodes []*Node,
	links []*graph.LinkData,
	vectors map[int]*ForceVector,
	debugInfo ForcesDebugData) {
	for _, link := range links {
		from := nodes[link.FromIdx]
		to := nodes[link.ToIdx]
		f := force.Apply(from.Point, to.Point)

		// Update force vectors
		ff := vectors[link.FromIdx]
		vectors[link.FromIdx] = ff.Add(f)
		ft := vectors[link.ToIdx]
		vectors[link.ToIdx] = ft.Sub(f)

		// Update debug information
		name := force.Name()
		debugInfo.Append(link.FromIdx, name, *f)
		debugInfo.Append(link.ToIdx, name, f.Negative())
	}
}

// BarneHutMethod applies force for each node agains each node,
// using Barne-Hut optimization method.
var BarneHutMethod = func(
	force Force,
	nodes []*Node,
	links []*graph.LinkData,
	vectors map[int]*ForceVector,
	debugInfo ForcesDebugData) {

	otree := NewOctreeFromNodes(nodes, force)

	for i, node := range nodes {
		f, err := otree.CalcForce(i)
		if err != nil {
			fmt.Println("[ERROR] Force calc failed:", i, err)
			continue
		}

		// Update force vectors
		f1 := vectors[node.Idx]
		vectors[node.Idx] = f1.Add(f)

		// Update debug information
		name := force.Name()
		debugInfo.Append(node.Idx, name, *f)
	}
}

// ForEachNode applies force to every node in the graph.
var ForEachNode = func(
	force Force,
	nodes []*Node,
	links []*graph.LinkData,
	vectors map[int]*ForceVector,
	debugInfo ForcesDebugData) {
	for i, node := range nodes {
		f := force.Apply(node.Point, nil)

		// Update force vectors
		ff := vectors[i]
		vectors[i] = ff.Add(f)

		// Update debug information
		name := force.Name()
		debugInfo.Append(i, name, *f)
	}
}
