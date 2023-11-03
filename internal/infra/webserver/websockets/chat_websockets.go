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
	MessageDB       repository.MessageInterface
	RabbitMQQueueCH *amqp.Channel
}

func NewChatWebsocket(userDB repository.UserInterface, messageDB repository.MessageInterface, rabbitMQQueueCH *amqp.Channel) *ChatHandler {
	ch := &ChatHandler{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Broadcast:       make(chan entityPkg.ChatMessage),
		Chatrooms:       make(map[string]map[*websocket.Conn]bool),
		UserDB:          userDB,
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

	c.loadMessagesDB(chatroomID, ws)

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
		msg.CreatedAt = time.Now()

		if strings.HasPrefix(msg.Content, "/stock=") {
			msg.Content = strings.TrimPrefix(msg.Content, "/stock=")

			// Get bot user
			botUser, err := c.UserDB.FindByEmail("bot@example.com")
			if err != nil {
				log.Printf("error fetching bot user: %v", err)
				continue
			}
			msg.UserID = botUser.ID.String()
			log.Printf("Queueing: ChatroomID: %s, UserID: %s, Content: %s", msg.ChatroomID, msg.UserID, msg.Content)
			go c.HandleStockCommand(msg)
		} else {
			log.Printf("Broadcast: ChatroomID: %s, UserID: %s, Content: %s", msg.ChatroomID, msg.UserID, msg.Content)
			c.Broadcast <- msg

			log.Printf("Persistence: ChatroomID: %s, UserID: %s, Content: %s", msg.ChatroomID, msg.UserID, msg.Content)
			go c.storeMessageDB(msg)
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

	go c.PostMessageToChatroom(msg)
	go c.storeMessageDB(msg)
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

func (c *ChatHandler) loadMessagesDB(chatroomID string, ws *websocket.Conn) {
	page := 1
	limit := 50
	sort := "asc"
	messages, _ := c.MessageDB.FindAllByChatRoomID(chatroomID, page, limit, sort)
	log.Printf("Loading %d messages from DB", len(messages))
	for _, msg := range messages {
		chatMsg := entityPkg.ChatMessage{
			ChatroomID: msg.ChatroomID.String(),
			UserID:     msg.UserID.String(),
			Content:    msg.Content,
			Username:   msg.User.Name,
			CreatedAt:  msg.CreatedAt,
		}

		if err := ws.WriteJSON(chatMsg); err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}
	}
}

func (c *ChatHandler) storeMessageDB(msg entityPkg.ChatMessage) {
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

	m, err := entity.NewMessage(chatRoomID, userID, msg.Content)
	if err != nil {
		log.Printf("error creating message: %v", err)
		return
	}

	err = c.MessageDB.Create(m)
	log.Printf("message created: %v", m)
	if err != nil {
		log.Printf("error persist message: %v", err)
		return
	}
}

func (c *ChatHandler) HandleStockCommand(msg entityPkg.ChatMessage) {
	rabbitmq.Publish(c.RabbitMQQueueCH, msg, "amq.direct", "bot")
}

func (c *ChatHandler) PostMessageToChatroom(msg entityPkg.ChatMessage) {
	c.Broadcast <- msg
}
