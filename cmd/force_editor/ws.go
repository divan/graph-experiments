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

	layout layout.Layout
	graph  *graph.Data

	positionHistory [][]*position
	currentIdx      int
}

func NewWSServer(layout layout.Layout) *WSServer {
	ws := &WSServer{
		upgrader: websocket.Upgrader{},
		layout:   layout,
	}
	ws.updatePositions()
	return ws
}

type WSResponse struct {
	Type      MsgType     `json:"type"`
	Idx       int         `json:"idx"`
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

// WebSocket response types
const (
	RespPositions MsgType = "positions"
	RespGraph     MsgType = "graph"
)

// WebSocket commands
const (
	CmdInit WSCommand = "init"
	CmdPrev           = "prev"
	CmdNext           = "next"
)

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
	case CmdPrev:
		ws.prev()
	case CmdNext:
		ws.next()
	}
}

func (ws *WSServer) sendGraphData(c *websocket.Conn) {
	msg := &WSResponse{
		Type:  RespGraph,
		Idx:   ws.currentIdx,
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
		Type:      RespPositions,
		Idx:       ws.currentIdx,
		Positions: ws.positionHistory[ws.currentIdx],
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

func (ws *WSServer) updatePositions() {
	nodes := ws.layout.Nodes()
	positions := []*position{}
	for i := 0; i < len(nodes); i++ {
		pos := &position{
			X: nodes[i].X,
			Y: nodes[i].Y,
			Z: nodes[i].Z,
		}
		positions = append(positions, pos)
	}
	ws.positionHistory = append(ws.positionHistory, positions)
	ws.currentIdx = len(ws.positionHistory) - 1
	ws.broadcastPositions()
}

func (ws *WSServer) broadcastPositions() {
	for i := 0; i < len(ws.hub); i++ {
		ws.sendPositions(ws.hub[i])
	}
}

func (ws *WSServer) updateGraph(data *graph.Data) {
	ws.graph = data

	ws.broadcastGraphData()
}

func (ws *WSServer) broadcastGraphData() {
	for i := 0; i < len(ws.hub); i++ {
		ws.sendGraphData(ws.hub[i])
	}
}

func (ws *WSServer) prev() {
	if ws.currentIdx > 0 {
		ws.currentIdx--
		ws.broadcastPositions()
	}
}

func (ws *WSServer) next() {
	if ws.currentIdx == len(ws.positionHistory)-1 {
		ws.layout.Calculate(1)
		ws.updatePositions()
	} else {
		ws.currentIdx++
		ws.broadcastPositions()
	}
}
