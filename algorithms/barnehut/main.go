package main

import (
	"fmt"
	"log"
)

func main() {
	g, err := NewGraphDataFromJSON("data_small.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded graph: %d nodes, %d links\n", len(g.Nodes), len(g.Links))

	l := &Layout3D{}
	nodes := l.InitCoordinates(g.Nodes)

	fmt.Println("Nodes:", nodes)
}
