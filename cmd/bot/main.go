package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/esirangelomub/go-chat-application/configs"
	dbutils "github.com/esirangelomub/go-chat-application/database"
	"github.com/esirangelomub/go-chat-application/internal/entity"
	"github.com/esirangelomub/go-chat-application/internal/infra/repository"
	entityPkg "github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"github.com/esirangelomub/go-chat-application/pkg/services/rabbitmq"
	"github.com/go-chi/jwtauth"
	amqp "github.com/rabbitmq/amqp091-go"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	StoOqURLTemplate = "https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv"
)

var tokenAuth *jwtauth.JWTAuth

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := dbutils.InitializeDB(config)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{}, &entity.Chatroom{}, &entity.Message{})

	userDB := repository.NewUser(db)

	rabbitMQConn, rabbitMQCH := rabbitmq.SetupRabbitMQ(config)
	defer rabbitMQConn.Close()
	defer rabbitMQCH.Close()

	tokenAuth = jwtauth.New("HS256", []byte(config.JwtSecret), nil)

	botMsgs := make(chan amqp.Delivery)
	go rabbitmq.Consume(rabbitMQCH, botMsgs, config.RabbitMQQueueBot)

	go func() {
		for msg := range botMsgs {
			log.Printf("Received message from: %s: %v", config.RabbitMQQueueBot, string(msg.Body))
			handleStockMessage(msg, rabbitMQCH)
		}
	}()

	wsMsgs := make(chan amqp.Delivery)
	go rabbitmq.Consume(rabbitMQCH, wsMsgs, config.RabbitMQQueueWebSocket)

	go func() {
		for msg := range wsMsgs {
			log.Printf("Received message from: %s: %v", config.RabbitMQQueueWebSocket, string(msg.Body))
			sendMessageToChatroom(msg, config, userDB)
		}
	}()

	select {}
}

func handleStockMessage(msg amqp.Delivery, ch *amqp.Channel) {
	var parsedMsg entityPkg.ChatMessage
	err := json.Unmarshal(msg.Body, &parsedMsg)
	if err != nil {
		log.Printf("Error parsing message body: %v", err)
		return
	}

	stockCode := parsedMsg.Content

	// Fetch and process stock data
	messageContent, err := fetchStockData(stockCode)
	if err != nil {
		log.Printf("Error fetching data for stock code %s: %v", stockCode, err)
		return
	}

	parsedMsg.Content = messageContent
	parsedMsg.Username = "Bot"
	parsedMsg.CreatedAt = time.Now()

	// Publish the message to RabbitMQ
	err = rabbitmq.Publish(ch, parsedMsg, "amq.direct", "websocket")
	if err != nil {
		log.Printf("Error publishing message to RabbitMQ: %v", err)
		return
	}

	msg.Ack(true)
}

func fetchStockData(stockCode string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(StoOqURLTemplate, stockCode))
	log.Printf("url %s", fmt.Sprintf(StoOqURLTemplate, stockCode))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Response body: %s", string(body))
	r := csv.NewReader(strings.NewReader(string(body)))
	records, err := r.ReadAll()
	if err != nil || len(records) < 2 {
		return "", err
	}

	// Extracting the Symbol and Close price from the CSV records
	symbol := records[1][0]
	closePrice := records[1][6]

	return fmt.Sprintf("%s quote is $%s per share", symbol, closePrice), nil
}

func sendMessageToChatroom(msg amqp.Delivery, config *configs.Conf, userDB repository.UserInterface) {
	var message entityPkg.ChatMessage
	err := json.Unmarshal(msg.Body, &message)
	if err != nil {
		log.Printf("Error parsing websocket message body: %v", err)
		return
	}

	// Serialize the message to JSON
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to serialize message: %v", err)
	}

	// Get bot user
	botUser, err := userDB.FindByEmail("bot@example.com")
	if err != nil {
		log.Printf("error fetching bot user: %v", err)
		return
	}

	jwt := tokenAuth
	jwtExpiresIn := config.JwtExpiresIn
	_, token, _ := jwt.Encode(map[string]interface{}{
		"userID": botUser.ID.String(),
		"exp":    time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	req, err := http.NewRequest("POST", config.BotURL, bytes.NewBuffer(msgBytes))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	defer resp.Body.Close()
	msg.Ack(true)
}
