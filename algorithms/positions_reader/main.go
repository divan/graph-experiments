package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Point struct {
	X, Y, Z int32
}

func main() {
	fd, err := os.Open("./positions.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	var points []Point
	for {
		var p Point
		err = binary.Read(fd, binary.LittleEndian, &p)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		points = append(points, p)
	}

	fmt.Printf("Found %d positions\n", len(points))
	data, err := json.MarshalIndent(points, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(data))
}
