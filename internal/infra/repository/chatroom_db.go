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

func (c *Chatroom) Create(product *entity.Chatroom) error {
	return c.DB.Create(product).Error
}

func (c *Chatroom) FindAll(page, limit int, sort string) ([]*entity.Chatroom, error) {
	var products []*entity.Chatroom
	var err error
	if sort == "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		err = c.DB.Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&products).Error
	} else {
		err = c.DB.Order("created_at " + sort).Find(&products).Error
	}
	return products, err
}

func (c *Chatroom) FindByID(id string) (*entity.Chatroom, error) {
	var product entity.Chatroom
	err := c.DB.First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (c *Chatroom) Update(product *entity.Chatroom) error {
	_, err := c.FindByID(product.ID.String())
	if err != nil {
		return err
	}
	return c.DB.Save(product).Error
}

func (c *Chatroom) Delete(product *entity.Chatroom) error {
	_, err := c.FindByID(product.ID.String())
	if err != nil {
		return err
	}
	return c.DB.Delete(product).Error
}
