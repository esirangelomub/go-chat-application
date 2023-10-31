package repository

import (
	"github.com/esirangelomub/go-chat-application/internal/entity"
	"gorm.io/gorm"
)

type ChatroomUser struct {
	DB *gorm.DB
}

func NewChatroomUser(db *gorm.DB) *ChatroomUser {
	return &ChatroomUser{DB: db}
}

func (c *ChatroomUser) Create(chatroomUser *entity.ChatroomUser) (*entity.ChatroomUser, error) {
	if err := c.DB.Create(chatroomUser).Error; err != nil {
		return nil, err
	}
	return chatroomUser, nil
}

func (c *ChatroomUser) FindByID(id string) (*entity.ChatroomUser, error) {
	var chatroomUser entity.ChatroomUser
	err := c.DB.First(&chatroomUser, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &chatroomUser, nil
}

func (c *ChatroomUser) Delete(chatroomUser *entity.ChatroomUser) error {
	_, err := c.FindByID(chatroomUser.ID.String())
	if err != nil {
		return err
	}
	return c.DB.Delete(chatroomUser).Error
}
