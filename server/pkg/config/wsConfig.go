package config

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func WsUpgrade() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}
