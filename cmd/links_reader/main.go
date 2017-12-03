package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

type Link struct {
	X, Y, Z int32
}

func main() {
	file := "./links.bin"
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	fd, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	var links int64
	for {
		var id int32
		err = binary.Read(fd, binary.LittleEndian, &id)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if id < 0 {
			fmt.Println("Link from:", -id)
		} else {
			fmt.Println("     to:", id)
			links++
		}
	}

	fmt.Printf("Found %d links\n", links)
}
