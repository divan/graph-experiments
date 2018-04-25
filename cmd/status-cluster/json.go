package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/divan/graph-experiments/graph"
	"github.com/ethereum/go-ethereum/p2p"
)

type JSONRPCResponse struct {
	Version string          `json:"jsonrpc"`
	Id      interface{}     `json:"id,omitempty"`
	Result  []*p2p.PeerInfo `json:"result"`
}

func processFile(path string) ([]*p2p.PeerInfo, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	var resp JSONRPCResponse
	err = json.NewDecoder(fd).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

func processFiles(dir string, g *graph.Graph) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(path, "whisper-") && strings.HasSuffix(path, ".json") {
			peers, err := processFile(path)
			id := idFromPath(path)
			if id == "" {
				return nil
			}
			for _, peer := range peers {
				AddPeer(g, id, peer)
			}
			return err
		}
		return nil
	})
	return err
}

func idFromPath(path string) string {
	path = strings.TrimPrefix(path, "whisper-")
	path = strings.TrimSuffix(path, ".json")
	return path
}

func printPeer(peer *p2p.PeerInfo) {
	direction := "→"
	if peer.Network.Inbound {
		direction = "←"
	}
	typ := isClient(peer.Name)
	fmt.Printf("'%s' (%v) %v %s\n", peer.ID[:6], peer.Network.RemoteAddress, typ, direction)
}

func isClient(name string) bool {
	parts := strings.Split(name, "/")
	arch := parts[2]
	client := false
	if strings.Contains(arch, "darwin") || strings.Contains(arch, "android") {
		client = true
	}
	return client
}
