package main

import (
	"github.com/divan/graph-experiments/graph"
	"github.com/gorilla/websocket"
)

func (ws *WSServer) sendGraphData(c *websocket.Conn) {
	msg := &WSResponse{
		Type:  RespGraph,
		Graph: ws.graph,
	}

	ws.sendMsg(c, msg)
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
