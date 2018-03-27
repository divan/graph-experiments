package main

import (
	"flag"
	"io"
	"log"
	"os"
)

func main() {
	var (
		dataKind = flag.String("type", "net", "Example random IPs network")
		netHosts = flag.Int("net.hosts", 20, "Number of hosts for net generator")
		netConns = flag.Int("net.connections", 4, "Number of connections between hosts for net generator")
		output   = flag.String("o", "data.json", "Output filename (use - for stdout)")
	)
	flag.Parse()

	var w io.Writer
	if *output == "-" {
		w = os.Stdout
	} else {
		fd, err := os.Create(*output)
		if err != nil {
			log.Fatal("Open file for writing failed:", err)
		}
		w = fd
		defer fd.Close()
	}

	if *dataKind == "net" {
		GenerateNetwork(w, *netHosts, *netConns)
	}
}
