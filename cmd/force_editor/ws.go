package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/divan/graph-experiments/graph"
	"github.com/divan/graph-experiments/layout"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	upgrader websocket.Upgrader
	hub      []*websocket.Conn

	positions []*position
	graph     *graph.Data
}

func NewWSServer() *WSServer {
	return &WSServer{
		upgrader: websocket.Upgrader{},
	}
}

type WSResponse struct {
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

type MsgType string
type WSCommand string

const Positions MsgType = "positions"
const Graph MsgType = "graph"
const CmdInit WSCommand = "init"

func (ws *WSServer) Handle(w http.ResponseWriter, r *http.Request) {
	c, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	ws.hub = append(ws.hub, c)

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", mt, err)
			break
		}
		ws.processRequest(c, mt, message)
	}
}

func (ws *WSServer) processRequest(c *websocket.Conn, mtype int, data []byte) {
	log.Printf("recv: %s\n", data)

	var cmd WSRequest
	err := json.Unmarshal(data, &cmd)
	if err != nil {
		log.Fatal("unmarshal command", err)
		return
	}

	switch cmd.Cmd {
	case CmdInit:
		ws.sendGraphData(c)
		ws.sendPositions(c)
	}
}

func (ws *WSServer) sendGraphData(c *websocket.Conn) {
	msg := &WSResponse{
		Type:  Graph,
		Graph: ws.graph,
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

func (ws *WSServer) sendPositions(c *websocket.Conn) {
	msg := &WSResponse{
		Type:      Positions,
		Positions: ws.positions,
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

func (ws *WSServer) updatePositions(nodes []*layout.Node) {
	positions := []*position{}
	for i := 0; i < len(nodes); i++ {
		pos := &position{
			X: nodes[i].X,
			Y: nodes[i].Y,
			Z: nodes[i].Z,
		}
		positions = append(positions, pos)
	}
	ws.positions = positions

	ws.broadcastPositions()
}

func (ws *WSServer) broadcastPositions() {
	for i := 0; i < len(ws.hub); i++ {
		ws.sendPositions(ws.hub[i])
	}
}

func (ws *WSServer) broadcastGraphData() {
	for i := 0; i < len(ws.hub); i++ {
		ws.sendGraphData(ws.hub[i])
	}
}

func (ws *WSServer) updateGraph(data *graph.Data) {
	ws.graph = data

	ws.broadcastGraphData()
}
