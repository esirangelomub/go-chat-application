package repository

import "github.com/esirangelomub/go-chat-application/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
}

type ChatRoomInterface interface {
	Create(chatRoom *entity.Chatroom) error
	FindAll(page, limit int, sort string) ([]*entity.Chatroom, error)
	FindByID(id string) (*entity.Chatroom, error)
	Update(chatRoom *entity.Chatroom) error
	Delete(chatRoom *entity.Chatroom) error
}

type MessageInterface interface {
	Create(message *entity.Message) error
	FindAllByChatRoomID(chatRoomId string, page, limit int, sort string) ([]*entity.Message, error)
}
