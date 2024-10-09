package websocket

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
	
	// Maximum undelivered messages
	maxUndeliveredMsg = 100
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		return true
	},
}

// Client represents a connection.
type Client struct {
	name string
	id   string          // Unique client ID for reconnection
	hub  *Hub            // Reference to the Hub
	conn *websocket.Conn // WebSocket connection
	send chan []byte     // Buffered channel for outbound messages
	undeliveredMsg [][]byte        // Queue for undelivered messages
}

func (c *Client) addUndeliveredMsg(message []byte) {
    if len(c.undeliveredMsg) >= maxUndeliveredMsg {
        // Deleting the oldest message to free up space
        c.undeliveredMsg = c.undeliveredMsg[1:]
    }
    c.undeliveredMsg = append(c.undeliveredMsg, message)
}

// readPump listens for incoming messages.
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }

        message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

        switch {
        case bytes.HasPrefix(message, []byte("join_room:")):
            roomName := string(message[len("join_room:"):])
            c.hub.HandleJoinRoom(c, roomName)
            c.send <- []byte("join_room_success:" + roomName)

        case bytes.HasPrefix(message, []byte("room_message:")):
            roomNameAndMessage := bytes.SplitN(message[len("room_message:"):], []byte(" "), 2)
            roomName := string(roomNameAndMessage[0])
            roomMessage := roomNameAndMessage[1]
            c.hub.handleRoomBroadcast(roomName, roomMessage)

        case bytes.HasPrefix(message, []byte("leave_room:")):
            roomName := string(message[len("leave_room:"):])
            c.hub.handleLeaveRoom(c, roomName)
            c.send <- []byte("leave_room_success:" + roomName)

        case bytes.HasPrefix(message, []byte("private_message:")):
            receiverAndMessage := bytes.SplitN(message[len("private_message:"):], []byte(" "), 2)
            receiver := string(receiverAndMessage[0])
            privateMessage := receiverAndMessage[1]
            c.hub.handlePrivateMessage(receiver, privateMessage)

        default:
            c.hub.broadcast <- message
        }
    }
}

// writePump sends messages to the client.
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
				// If the connection is broken, add the message to the unread queue
				c.undeliveredMsg = c.addUndeliveredMsg(message)
				return
			}
			w.Write(message)

			// Add queued messages to current WebSocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
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

// Hub manages clients and rooms.
type Hub struct {
	clients    map[string]*Client // Track clients by ID for reconnection
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	rooms      map[string]map[*Client]bool
}

// Create a new room.
func (h *Hub) createRoom(name string) {
	if _, exists := h.rooms[name]; !exists {
		h.rooms[name] = make(map[*Client]bool)
	}
}

// Join a room.
func (h *Hub) joinRoom(client *Client, room string) {
	if _, exists := h.rooms[room]; exists {
		h.rooms[room][client] = true
	}
}

// Leave a room.
func (h *Hub) handleLeaveRoom(client *Client, room string) {
	if roomClients, ok := h.rooms[room]; ok {
		delete(roomClients, client)
		if len(roomClients) == 0 {
			delete(h.rooms, room)
		}
	}
}

// Broadcast message to a room.
func (h *Hub) handleRoomBroadcast(roomName string, message []byte) {
	if clients, ok := h.rooms[roomName]; ok {
		for client := range clients {
			client.send <- message
		}
	}
}

// Handle private message.
func (h *Hub) handlePrivateMessage(receiverName string, message []byte) {
	for _, client := range h.clients {
		if client.name == receiverName {
			client.send <- message
		}
	}
}

// Handle join room.
func (h *Hub) HandleJoinRoom(client *Client, roomName string) {
	h.createRoom(roomName)
	h.joinRoom(client, roomName)
}

// Run starts the Hub.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		case message := <-h.broadcast:
			for _, client := range h.clients {
				client.send <- message
			}
		}
	}
}

// Serve WebSocket connection and handle reconnections.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clientID := r.URL.Query().Get("client_id")
	var client *Client
	if clientID != "" && hub.clients[clientID] != nil {
		// Reconnect existing client
		client = hub.clients[clientID]
		client.conn = conn
		client.sendUndeliveredMsg() // function that sends unread messages
	} else {
		// New client connection
		clientID = uuid.NewString()
		client = &Client{
			hub:  hub,
			conn: conn,
			send: make(chan []byte, 256),
			id:   clientID,
			name: "root",
			undeliveredMsg: [][]byte{},
		}
	}

	client.hub.register <- client
	go client.writePump()
	go client.readPump()
}

// Send all unread messages to the client after reconnection.
func (c *Client) sendUndeliveredMsg() {
	for _, msg := range c.undeliveredMsg {
		c.send <- msg
	}
	// Clearing the queue of unread messages after sending
	c.undeliveredMsg = [][]byte{}
}

// WebSocketServer manages the WebSocket server.
type WebSocketServer struct{}

// NewWebSocketServer creates a new server.
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{}
}

// NewWsServer starts the WebSocket server.
func (wss *WebSocketServer) NewWsServer(addr string) {
	var _addr = flag.String("addr", addr, "http service address")
	flag.Parse()
	hub := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
		rooms:      make(map[string]map[*Client]bool),
	}
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*_addr, nil)
	if err != nil {
		log.Fatal("WebSocket server error:", err)
	}
}
