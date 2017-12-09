package webmods

import (
  "html/template"
  "net/http"
  "log"
  "fmt"
  "time"
  "github.com/gorilla/websocket"
)

const (
	writeWait = 1 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Msg struct {
  Type string `json:"type"`
}

type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
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
		_, message, err := c.conn.ReadMessage()
    if err != nil {
      if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
        log.Printf("error: %v", err)
      }
      break
    }
    fmt.Println(string(message))
		c.hub.broadcast <- message
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
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
      if err != nil {
        return
      }
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
        w.Write([]byte{'\n'})
				w.Write(<-c.send)
        fmt.Println(string(message))
			}
      if err := w.Close(); err != nil {
        return
      }
    case <-ticker.C:
      c.conn.SetWriteDeadline(time.Now().Add(writeWait))
      if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
        return
      }
    }
	}
}

func CLR(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("./web/templates/clr.html")
  t.Execute(w, nil)
}

func SendCLR(hub *Hub, w http.ResponseWriter, r *http.Request) {
  conn, _ := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
  client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
  client.hub.register <- client
  
  go client.writePump()
  go client.readPump()
}