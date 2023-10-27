package websockets

import (
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

// Message represents the structure of a chat message.
type Message struct {
	ChatroomID string `json:"chatroom_id"`
	UserID     string `json:"user_id"`
	Content    string `json:"content"`
}

type ChatHandler struct {
	Upgrader  websocket.Upgrader
	Broadcast chan Message
	Mutex     sync.Mutex
	Chatrooms map[string]map[*websocket.Conn]bool
}

func NewChatWebsocket() *ChatHandler {
	ch := &ChatHandler{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Broadcast: make(chan Message),
		Chatrooms: make(map[string]map[*websocket.Conn]bool),
	}

	go ch.HandleMessages() // Start the goroutine here

	return ch
}

func (c *ChatHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	chatroomID := chi.URLParam(r, "chatroomID")
	log.Printf("New WebSocket connection established for ChatroomID: %s", chatroomID)

	c.Mutex.Lock()
	if _, ok := c.Chatrooms[chatroomID]; !ok {
		c.Chatrooms[chatroomID] = make(map[*websocket.Conn]bool)
	}
	c.Chatrooms[chatroomID][ws] = true
	c.Mutex.Unlock()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			c.Mutex.Lock()
			delete(c.Chatrooms[chatroomID], ws)
			c.Mutex.Unlock()
			break
		}

		log.Printf("Received WebSocket message from ChatroomID: %s, UserID: %s, Content: %s", msg.ChatroomID, msg.UserID, msg.Content)
		c.Broadcast <- msg
	}
}

func (c *ChatHandler) HandleMessages() {
	for {
		msg := <-c.Broadcast
		log.Printf("Handling message for ChatroomID: %s, Content: %s", msg.ChatroomID, msg.Content)

		c.Mutex.Lock()
		for client := range c.Chatrooms[msg.ChatroomID] {
			log.Printf("Sending message to client in ChatroomID: %s", msg.ChatroomID)
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(c.Chatrooms[msg.ChatroomID], client)
			}
		}
		c.Mutex.Unlock()
	}
}
