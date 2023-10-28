package repository

import (
	"github.com/esirangelomub/go-chat-application/internal/entity"
	"gorm.io/gorm"
)

type Chatroom struct {
	DB *gorm.DB
}

func NewChatroom(db *gorm.DB) *Chatroom {
	return &Chatroom{DB: db}
}

func (c *Chatroom) Create(chatroom *entity.Chatroom) error {
	return c.DB.Create(chatroom).Error
}

func (c *Chatroom) FindAll(page, limit int, sort string) ([]*entity.Chatroom, error) {
	var chatRooms []*entity.Chatroom
	var err error
	if sort == "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		err = c.DB.Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&chatRooms).Error
	} else {
		err = c.DB.Order("created_at " + sort).Find(&chatRooms).Error
	}
	return chatRooms, err
}

func (c *Chatroom) FindByID(id string) (*entity.Chatroom, error) {
	var chatRoom entity.Chatroom
	err := c.DB.First(&chatRoom, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &chatRoom, nil
}

func (c *Chatroom) Update(chatRoom *entity.Chatroom) error {
	_, err := c.FindByID(chatRoom.ID.String())
	if err != nil {
		return err
	}
	return c.DB.Save(chatRoom).Error
}

func (c *Chatroom) Delete(chatRoom *entity.Chatroom) error {
	_, err := c.FindByID(chatRoom.ID.String())
	if err != nil {
		return err
	}
	return c.DB.Delete(chatRoom).Error
}
