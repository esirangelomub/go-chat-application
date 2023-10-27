package repository

import "github.com/esirangelomub/go-chat-application/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

type ChatRoomInterface interface {
	Create(product *entity.Chatroom) error
	FindAll(page, limit int, sort string) ([]*entity.Chatroom, error)
	FindByID(id string) (*entity.Chatroom, error)
	Update(product *entity.Chatroom) error
	Delete(product *entity.Chatroom) error
}
