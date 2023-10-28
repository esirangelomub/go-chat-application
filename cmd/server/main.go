package main

import (
	"github.com/esirangelomub/go-chat-application/configs"
	dbutils "github.com/esirangelomub/go-chat-application/database"
	"github.com/esirangelomub/go-chat-application/internal/entity"
	"github.com/esirangelomub/go-chat-application/internal/infra/repository"
	"github.com/esirangelomub/go-chat-application/internal/infra/webserver/handlers"
	"github.com/esirangelomub/go-chat-application/internal/infra/webserver/websockets"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"net/http"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := dbutils.InitializeDB(config)
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{}, &entity.Chatroom{}, &entity.ChatroomUser{}, &entity.Message{})

	userDB := repository.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	chatRoomDB := repository.NewChatroom(db)
	chatRoomHandler := handlers.NewChatRoomHandler(chatRoomDB)

	//chatRoomUserDB := repository.NewChatroomUser(db)
	//messageDB := repository.NewMessage(db)
	//chatService := service.NewChatService(chatRoomUserDB, messageDB)

	chatWebsocket := websockets.NewChatWebsocket()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", config.JwtExpiresIn))

	r.Post("/users", userHandler.CreateUser)
	r.Post("/users/generate_token", userHandler.GetJWT)

	r.Route("/chats/rooms", func(r chi.Router) {
		// middleware to verify and authenticate JWT
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", chatRoomHandler.Create)
		r.Get("/", chatRoomHandler.List)
		r.Get("/{id}", chatRoomHandler.Fetch)
		r.Put("/{id}", chatRoomHandler.Update)
		r.Delete("/{id}", chatRoomHandler.Delete)
	})

	r.Route("/ws", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/{chatroomID}", chatWebsocket.HandleConnections)
	})

	http.ListenAndServe(":8000", r)
}
