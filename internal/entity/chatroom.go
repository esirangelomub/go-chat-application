package entity

import (
	"errors"
	"github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"time"
)

var (
	ErrIDIsRequired   = errors.New("id is required")
	ErrInvalidID      = errors.New("invalid id")
	ErrNameIsRequired = errors.New("name is required")
)

type Chatroom struct {
	ID          entity.ID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewChatRoom(name string, description string) (*Chatroom, error) {
	c := &Chatroom{
		ID:          entity.NewID(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}

func (p *Chatroom) Validate() error {
	if p.ID.String() == "" {
		return ErrIDIsRequired
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalidID
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	return nil
}
