package client

import "github.com/gorilla/websocket"

type Client struct {
	Username string
	Conn     *websocket.Conn
}
