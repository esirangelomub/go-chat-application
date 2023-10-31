package websockets

import (
	"encoding/json"
	"github.com/esirangelomub/go-chat-application/internal/entity"
	"github.com/esirangelomub/go-chat-application/internal/infra/repository"
	entityPkg "github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"github.com/esirangelomub/go-chat-application/pkg/services/rabbitmq"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type ChatHandler struct {
	Upgrader        websocket.Upgrader
	Broadcast       chan entityPkg.ChatMessage
	Mutex           sync.Mutex
	Chatrooms       map[string]map[*websocket.Conn]bool
	UserDB          repository.UserInterface
	ChatroomUserDB  repository.ChatRoomUserInterface
	MessageDB       repository.MessageInterface
	RabbitMQQueueCH *amqp.Channel
}

func NewChatWebsocket(userDB repository.UserInterface, chatroomUserDB repository.ChatRoomUserInterface,
	messageDB repository.MessageInterface, rabbitMQQueueCH *amqp.Channel) *ChatHandler {
	ch := &ChatHandler{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Broadcast:       make(chan entityPkg.ChatMessage),
		Chatrooms:       make(map[string]map[*websocket.Conn]bool),
		UserDB:          userDB,
		ChatroomUserDB:  chatroomUserDB,
		MessageDB:       messageDB,
		RabbitMQQueueCH: rabbitMQQueueCH,
	}

	go ch.HandleMessages()

	return ch
}

func (c *ChatHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := c.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	chatroomID := chi.URLParam(r, "chatroomID")
	if chatroomID == "" {
		log.Printf("ChatroomID not found")
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	userID := claims["userID"]
	if userID == "" {
		log.Printf("User not found or invalid")
		return
	}
	log.Printf("New WebSocket connection established for ChatroomID: %s, UserID: %s", chatroomID, userID)

	// Lock once and make all necessary changes
	c.Mutex.Lock()
	if _, ok := c.Chatrooms[chatroomID]; !ok {
		c.Chatrooms[chatroomID] = make(map[*websocket.Conn]bool)
	}
	c.Chatrooms[chatroomID][ws] = true
	c.Mutex.Unlock()

	for {
		var msg entityPkg.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			c.removeClientFromChatroom(chatroomID, ws)
			break
		}

		msg.UserID = userID.(string)
		msg.ChatroomID = chatroomID

		user, err := c.UserDB.FindByID(userID.(string))
		if err != nil {
			log.Printf("error fetching user: %v", err)
			continue
		}
		msg.Username = user.Name
		msg.Timestamp = time.Now().Unix()

		// Check if the message content starts with the stock command prefix
		if strings.HasPrefix(msg.Content, "/stock=") {
			log.Printf("Queueing: ChatroomID: %s, UserID: %s, Content: %s", msg.ChatroomID, msg.UserID, msg.Content)
			go c.HandleStockCommand(msg)
		} else {
			log.Printf("Broadcast: ChatroomID: %s, UserID: %s, Content: %s", msg.ChatroomID, msg.UserID, msg.Content)
			c.Broadcast <- msg

			log.Printf("Persistence: ChatroomID: %s, UserID: %s, Content: %s", msg.ChatroomID, msg.UserID, msg.Content)
			go c.PersistMessage(msg)
		}
	}
}

func (c *ChatHandler) removeClientFromChatroom(chatroomID string, ws *websocket.Conn) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	delete(c.Chatrooms[chatroomID], ws)
}

func (c *ChatHandler) HandleBotMessages(w http.ResponseWriter, r *http.Request) {
	var msg entityPkg.ChatMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	c.PostMessageToChatroom(msg)
}

func (c *ChatHandler) HandleMessages() {
	for {
		msg := <-c.Broadcast
		log.Printf("Handling message for ChatroomID: %s, Content: %s", msg.ChatroomID, msg.Content)

		c.Mutex.Lock()
		clients := c.Chatrooms[msg.ChatroomID]
		c.Mutex.Unlock()

		for client := range clients {
			log.Printf("Sending message to client in ChatroomID: %s", msg.ChatroomID)
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				c.removeClientFromChatroom(msg.ChatroomID, client)
			}
		}
	}
}

func (c *ChatHandler) PersistMessage(msg entityPkg.ChatMessage) {
	chatRoomID, err := entityPkg.ParseID(msg.ChatroomID)
	if err != nil {
		log.Printf("error parsing chatroom id: %v", err)
		return
	}

	userID, err := entityPkg.ParseID(msg.UserID)
	if err != nil {
		log.Printf("error parsing user id: %v", err)
		return
	}

	chu, err := entity.NewChatroomUser(chatRoomID, userID)
	if err != nil {
		log.Printf("error creating chatroom user: %v", err)
		return
	}

	chatRoomUser, err := c.ChatroomUserDB.Create(chu)
	if err != nil {
		log.Printf("error creating chatroom user: %v", err)
		return
	}

	m, err := entity.NewMessage(chatRoomUser.ID, msg.Content)
	if err != nil {
		log.Printf("error creating message: %v", err)
		return
	}

	err = c.MessageDB.Create(m)
	if err != nil {
		if delErr := c.ChatroomUserDB.Delete(chatRoomUser); delErr != nil {
			log.Printf("Failed to delete chatUser after message creation error: %v", delErr)
		}
		return
	}
}

func (c *ChatHandler) HandleStockCommand(msg entityPkg.ChatMessage) {
	rabbitmq.Publish(c.RabbitMQQueueCH, msg, "amq.direct", "bot")
}

func (c *ChatHandler) PostMessageToChatroom(msg entityPkg.ChatMessage) {
	c.Broadcast <- msg
}
