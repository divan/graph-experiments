package main

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/divan/graph-experiments/export"
	"github.com/divan/graph-experiments/graph"
	"github.com/divan/graph-experiments/layout"
	"github.com/gorilla/websocket"
)

func (ws *WSServer) sendGraphData(c *websocket.Conn) {
	var buf bytes.Buffer
	err := export.NewJSON(&buf, false).ExportGraph(ws.graph)
	if err != nil {
		log.Fatal("Can't marshal graph to JSON")
	}
	msg := &WSResponse{
		Type:  RespGraph,
		Graph: json.RawMessage(buf.Bytes()),
	}

	ws.sendMsg(c, msg)
}

func (ws *WSServer) updateGraph(g *graph.Graph, l layout.Layout) {
	ws.graph = g
	ws.layout = l

	ws.broadcastGraphData()
}

func (ws *WSServer) broadcastGraphData() {
	for i := 0; i < len(ws.hub); i++ {
		ws.sendGraphData(ws.hub[i])
	}
}
