package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/divan/graph-experiments/graph"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Msg struct {
	Type      MsgType     `json:"type"`
	Positions []*position `json:"positions,omitempty"`
	Graph     *graph.Data `json:"graph,omitempty"`
}

type WSRequest struct {
	Cmd WSCommand `json:"cmd"`
}

type position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	Z int32 `json:"z"`
}

var positionsData = []*position{}

type MsgType string
type WSCommand string

const Positions MsgType = "positions"
const Graph MsgType = "graph"
const CmdInit WSCommand = "init"

func handleWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", mt, err)
			break
		}
		processWSCommand(c, mt, message)
	}
}

func processWSCommand(c *websocket.Conn, mtype int, data []byte) {
	log.Printf("recv: %s\n", data)

	var cmd WSRequest
	err := json.Unmarshal(data, &cmd)
	if err != nil {
		log.Fatal("unmarshal command", err)
		return
	}

	switch cmd.Cmd {
	case CmdInit:
		sendGraphData(c)
		sendPositions(c)
	}
}

func sendGraphData(c *websocket.Conn) {
	msg := &Msg{
		Type:  Graph,
		Graph: graphData,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("write:", err)
		return
	}

	fmt.Println("Sending", string(data))

	err = c.WriteMessage(1, data)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func sendPositions(c *websocket.Conn) {
	msg := &Msg{
		Type:      Positions,
		Positions: positionsData,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("write:", err)
		return
	}

	fmt.Println("Sending", string(data))

	err = c.WriteMessage(1, data)
	if err != nil {
		log.Println("write:", err)
		return
	}
}
