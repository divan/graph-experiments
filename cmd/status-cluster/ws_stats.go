package main

import "github.com/gorilla/websocket"

func (ws *WSServer) sendStats(c *websocket.Conn) {
	msg := &WSResponse{
		Type:  RespStats,
		Stats: *ws.stats.Stats(),
	}

	ws.sendMsg(c, msg)
}

func (ws *WSServer) broadcastStats() {
	for i := 0; i < len(ws.hub); i++ {
		ws.sendStats(ws.hub[i])
	}
}
