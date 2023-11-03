// message_test.go
package entity

import (
	"github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMessage(t *testing.T) {
	chatRoomID := entity.NewID()
	userID := entity.NewID()

	// Test case 1: Valid message
	message, err := NewMessage(chatRoomID, userID, "Hello, World!")
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.NotEmpty(t, message.ID)
	assert.Equal(t, chatRoomID, message.ChatroomID)
	assert.Equal(t, userID, message.UserID)
	assert.Equal(t, "Hello, World!", message.Content)
	assert.NotNil(t, message.CreatedAt)
}

func TestMessage_Validate(t *testing.T) {
	chatRoomID := entity.NewID()
	userID := entity.NewID()

	// Test case 1: Valid message
	message := &Message{
		ID:         entity.NewID(),
		ChatroomID: chatRoomID,
		UserID:     userID,
		Content:    "Hello, World!",
	}
	err := message.Validate()
	assert.Nil(t, err)
}
