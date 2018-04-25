package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/divan/graph-experiments/graph"
)

func processSSH(g *graph.Graph, hosts []string) {
	for _, host := range hosts {
		GetAdminPeersOutput(g, host)
	}
}

// GetAdminPeersOutput runs ipc `admin.peers` JSON-RPC request against geth.ipc
// sockets found in /var (beloning to docker).
//
// host comes in the form of "<IP>:docker-exposed-port", and this port will be used
// to match/identify node. Actual connection happens on default SSH port (22).
func GetAdminPeersOutput(g *graph.Graph, hostport string) {
	host, _ := parseHostPort(hostport)
	ssh, err := NewSSH(host, "root")
	if err != nil {
		log.Println("[ERROR] Connect to SSH failed:", err)
		return
	}

	adminNodeInfoCmd := `{"jsonrpc":"2.0","method":"admin_nodeInfo","params":[],"id":1}`
	adminPeersCmd := `{"jsonrpc":"2.0","method":"admin_peers","params":[],"id":2}`
	out, err := ssh.Exec(`for i in $(find /var -name geth.ipc); do
		echo "{"
		echo -n '"nodeID": "'
		echo -n '` + adminNodeInfoCmd + `' | nc -U $i | json_pp | grep '"id" : "' | cut -d'"' -f4 | tr -d '\n';
		echo -n '", "adminPeers": '
		echo -n '` + adminPeersCmd + `' | nc -U $i;
		echo -n "},"
	done
	`)
	if err != nil {
		log.Fatal(err)
	}

	out = strings.TrimSuffix(strings.TrimSpace(out), ",")
	out = "[" + out + "]"

	var peersOutput []*PeerSSHOutput
	err = json.Unmarshal([]byte(out), &peersOutput)
	if err != nil {
		log.Fatal("Failed to unmarshal JSON:", err, len(out), out)
	}
	for _, p := range peersOutput {
		fmt.Println("Adding peers for", p.ID[:6], " (", len(p.AdminPeers.Result), "peers )")
		for _, peer := range p.AdminPeers.Result {
			AddPeer(g, p.ID, peer)
		}
	}
}

// PeerSSHOutput represents output from SSH command executed on status cluster machine.
type PeerSSHOutput struct {
	ID         string           `json:"nodeID"`
	AdminPeers *JSONRPCResponse `json:"AdminPeers"`
}

func parseHostPort(hostport string) (string, string) {
	host, port, _ := net.SplitHostPort(hostport)
	return host, port
}
