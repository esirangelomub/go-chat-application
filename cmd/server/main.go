package main

import (
	"github.com/esirangelomub/go-chat-application/configs"
	dbutils "github.com/esirangelomub/go-chat-application/database"
	"github.com/esirangelomub/go-chat-application/internal/entity"
	"github.com/esirangelomub/go-chat-application/internal/infra/repository"
	"github.com/esirangelomub/go-chat-application/internal/infra/webserver/handlers"
	"github.com/esirangelomub/go-chat-application/internal/infra/webserver/websockets"
	middlewarePkg "github.com/esirangelomub/go-chat-application/pkg/middlewares"
	"github.com/esirangelomub/go-chat-application/pkg/services/rabbitmq"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth"
	"net/http"
	"time"
)

func main() {
	// load configs
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// database connection
	db, err := dbutils.InitializeDB(config)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{}, &entity.Chatroom{}, &entity.ChatroomUser{}, &entity.Message{})

	// rabbitmq connection
	rabbitMQConn, rabbitMQCH := rabbitmq.SetupRabbitMQ(config)
	defer rabbitMQConn.Close()
	defer rabbitMQCH.Close()

	userDB := repository.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	chatRoomDB := repository.NewChatroom(db)
	chatRoomHandler := handlers.NewChatRoomHandler(chatRoomDB)

	chatRoomUserDB := repository.NewChatroomUser(db)
	messageDB := repository.NewMessage(db)

	chatWebsocket := websockets.NewChatWebsocket(userDB, chatRoomUserDB, messageDB, rabbitMQCH)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", config.JwtExpiresIn))
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	r.Use(middlewarePkg.Handler)
	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/generate_token", userHandler.GetJWT)

	r.With(jwtauth.Verifier(config.TokenAuth), jwtauth.Authenticator).Get("/users/me", userHandler.Logged)

	r.Route("/chats/rooms", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", chatRoomHandler.Create)
		r.Get("/", chatRoomHandler.List)
		r.Get("/{id}", chatRoomHandler.Fetch)
		r.Put("/{id}", chatRoomHandler.Update)
		r.Delete("/{id}", chatRoomHandler.Delete)
	})

	r.Route("/ws", func(r chi.Router) {
		r.Use(middlewarePkg.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/{chatroomID}", chatWebsocket.HandleConnections)
		r.Post("/bot", chatWebsocket.HandleBotMessages)
	})

	http.ListenAndServe(":8000", r)
}
