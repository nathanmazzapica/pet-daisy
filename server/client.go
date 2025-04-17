package server

import (
	"github.com/gorilla/websocket"
	"github.com/nathanmazzapica/pet-daisy/db"
	"github.com/nathanmazzapica/pet-daisy/logger"
	"log"
	"sync"
	"time"
)

// Client represents a WebSocket connection
type Client struct {
	conn        *websocket.Conn
	user        db.User
	lastPetTime time.Time
	susPets     int
	petTimes    [PET_WINDOW]time.Time
	sessionPets int
	mutex       sync.Mutex

	hub  *Server
	send chan ServerMessage
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

func (s *Server) newClient(conn *websocket.Conn, user db.User) *Client {
	client := &Client{
		conn: conn,
		user: user,
		hub:  s,
		send: make(chan ServerMessage, 256),
	}

	return client
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.LogError(err)
			}
			break
		}

		message, err := buildClientMessage(msg, c)

		if err != nil {
			log.Printf("Error processing raw message: %v", err)
			continue
		}

		c.hub.in <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				//Hub closed channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) DisplayName() string {
	return c.user.DisplayName
}
