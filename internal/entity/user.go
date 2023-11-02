package entity

import (
	"github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"golang.org/x/crypto/bcrypt"
)

type UserType string

const (
	BOT  UserType = "BOT"
	USER UserType = "USER"
)

type User struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Type     UserType  `json:"type" gorm:"type:varchar(4);default:'USER'"` // USER or BOT
}

func NewUser(name, email, password string, tp UserType) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       entity.NewID(),
		Name:     name,
		Email:    email,
		Password: string(hash),
		Type:     tp,
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
