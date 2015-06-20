package tippler

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type connection struct {
	ws      *websocket.Conn
	send    chan []byte
	channel string
}

type hub struct {
	register    chan *connection
	unregister  chan *connection
	connections map[*connection]bool
	broadcast   chan *broadcastMsg
	channels    map[string]map[*connection]bool
}

var h = hub{
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	broadcast:   make(chan *broadcastMsg),
	connections: make(map[*connection]bool),
	channels:    make(map[string]map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			if _, ok := h.channels[c.channel]; !ok {
				h.channels[c.channel] = make(map[*connection]bool)
			}
			h.channels[c.channel][c] = true
			log.Println("Client connected")
		case bm := <-h.broadcast:
			h.broadcastTweet(bm)
		case c := <-h.unregister:
			delete(h.connections, c)
			delete(h.channels[c.channel], c)
			close(c.send)
			log.Println("Client disconnected")
		}
	}
}

func (h *hub) broadcastTweet(msg *broadcastMsg) {
	for conn, _ := range h.channels[msg.channel] {
		conn.send <- []byte(msg.msg)
	}
	for conn, _ := range h.channels["all"] {
		conn.send <- []byte(msg.msg)
	}
}

// serverWS handles websocket requests from the peer.
func serveWS(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	channel := r.URL.Path[1:]
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws, channel: channel}
	h.register <- c
	go c.ReadLoop()
	c.WriteLoop()
}

func (c *connection) ReadLoop() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	for {
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			return
		}
	}
}

func (c *connection) WriteLoop() {
	defer func() {
		c.ws.Close()
	}()
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func RunServer() {
	go h.run() // needs to be a goroutine so we can receive messages from clients
	http.HandleFunc("/", serveWS)
	log.Println("listning on port 9292")
	err := http.ListenAndServe(":9292", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
