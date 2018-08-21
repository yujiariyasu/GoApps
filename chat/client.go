package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// clientはチャットを行っている1人のユーザーを表します。
type client struct {
	// socketはこのクライアントのためのWebSocketです。
	socket *websocket.Conn
	// sendはメッセージが送られるチャネルです。
	send chan *message
	// roomはこのクライアントが参加しているチャットルームです。
	room *room
	// ユーザに関する情報を保持
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			if avatarURL, ok := c.userData["avatar_url"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}
func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
