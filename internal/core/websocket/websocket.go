package websocket

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"time"

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

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	name string
	hub  *Hub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
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
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// Parse and handle different events
		if bytes.HasPrefix(message, []byte("join_room:")) {
			roomName := string(message[len("join_room:"):])
			c.hub.HandleJoinRoom(c, roomName)
			// Send confirmation back to client
			confirmationMessage := "join_room_success:" + roomName
			c.send <- []byte(confirmationMessage)
		} else if bytes.HasPrefix(message, []byte("room_message:")) {
			// Handle room message event
			roomNameAndMessage := bytes.SplitN(message[len("room_message:"):], []byte(" "), 2)
			roomName := string(roomNameAndMessage[0])
			roomMessage := roomNameAndMessage[1]
			c.hub.handleRoomBroadcast(roomName, roomMessage)
		} else if bytes.HasPrefix(message, []byte("leave_room:")) {
			// Handle leave room event
			roomName := string(message[len("leave_room:"):])
			c.hub.handleLeaveRoom(c, roomName)
			// Send confirmation back to client
			confirmationMessage := "leave_room_success:" + roomName
			c.send <- []byte(confirmationMessage)
		} else if bytes.HasPrefix(message, []byte("private_message:")) {
			// Handle private message event
			receiverAndMessage := bytes.SplitN(message[len("private_message:"):], []byte(" "), 2)
			receiver := string(receiverAndMessage[0])
			privateMessage := receiverAndMessage[1]
			c.hub.handlePrivateMessage(receiver, privateMessage)
		} else {
			// Handle global messages or unhandled events
			c.hub.broadcast <- message
		}
	}
}

func (h *Hub) handleRoomBroadcast(roomName string, message []byte) {
	if room, ok := h.rooms[roomName]; ok {
		for client := range room {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(room, client)
			}
		}
	}
}

func (h *Hub) handleLeaveRoom(client *Client, roomName string) {
	if room, ok := h.rooms[roomName]; ok {
		if _, exists := room[client]; exists {
			delete(room, client)
			if len(room) == 0 {
				delete(h.rooms, roomName)
			}
		}
	}
}

func (h *Hub) handlePrivateMessage(receiverName string, message []byte) {
	for client := range h.clients {
		if client.name == receiverName {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
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
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
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

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), name: "root"}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Maps room names to clients
	rooms map[string]map[*Client]bool
}

func (h *Hub) createRoom(name string) {
	if _, exists := h.rooms[name]; !exists {
		h.rooms[name] = make(map[*Client]bool)
	}
}

func (h *Hub) joinRoom(client *Client, room string) {
	if _, exists := h.rooms[room]; exists {
		h.rooms[room][client] = true
	}
}

func (h *Hub) LeaveRoom(client *Client, room string) {
	if _, exists := h.rooms[room]; exists {
		delete(h.rooms[room], client)
	}
}

func (h *Hub) BroadcastToRoom(room string, message []byte) {
	if clients, exists := h.rooms[room]; exists {
		for client := range clients {
			client.send <- message
		}
	}
}

func (c *Client) Emit(event string, data []byte) {
	message := append([]byte(event+": "), data...)
	c.send <- message
}

func (h *Hub) HandleJoinRoom(client *Client, roomName string) {
	h.createRoom(roomName)
	h.joinRoom(client, roomName)
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
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

type WebSocketServer struct{}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{}
}

func (wss *WebSocketServer) NewWsServer(addr string) {
	var _addr = flag.String("addr", ":8080", "http service address")
	flag.Parse()
	hub := newHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	err := http.ListenAndServe(*_addr, nil)
	if err != nil {
		log.Fatal("WSS ListenAndServe: ", err)
	}
}
