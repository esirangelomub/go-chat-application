package repository

import (
	"github.com/esirangelomub/go-chat-application/internal/entity"
	"gorm.io/gorm"
	"log"
)

type Message struct {
	DB *gorm.DB
}

func NewMessage(db *gorm.DB) *Message {
	return &Message{DB: db}
}

func (c *Message) Create(message *entity.Message) error {
	return c.DB.Create(message).Error
}

func (c *Message) FindAllByChatRoomID(chatRoomId string, page, limit int, sort string) ([]*entity.Message, error) {
	var messages []*entity.Message
	var err error

	if sort == "" {
		sort = "asc"
	}

	if page != 0 && limit != 0 {
		err = c.DB.Limit(limit).Offset((page-1)*limit).Order("timestamp "+sort).Preload("User").Find(&messages, "chatroom_id = ?", chatRoomId).Error
	} else {
		err = c.DB.Order("timestamp "+sort).Preload("User").Find(&messages, "chatroom_id = ?", chatRoomId).Error
	}
	log.Printf("chatRoomUserId: %v: page: %v: limit: %v: sort: %v", chatRoomId, page, limit, sort)
	return messages, err
}
