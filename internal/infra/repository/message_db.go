package repository

import (
	"github.com/esirangelomub/go-chat-application/internal/entity"
	"gorm.io/gorm"
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
