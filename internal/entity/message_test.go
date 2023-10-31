// message_test.go
package entity

import (
	"github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMessage(t *testing.T) {
	chatroomUserID := entity.NewID()

	// Test case 1: Valid message
	message, err := NewMessage(chatroomUserID, "Hello, World!")
	assert.Nil(t, err)
	assert.NotNil(t, message)
	assert.NotEmpty(t, message.ID)
	assert.Equal(t, chatroomUserID, message.ChatroomUserID)
	assert.Equal(t, "Hello, World!", message.Content)
	assert.NotZero(t, message.Timestamp)
}

func TestMessage_Validate(t *testing.T) {
	chatroomUserID := entity.NewID()

	// Test case 1: Valid message
	message := &Message{
		ID:             entity.NewID(),
		ChatroomUserID: chatroomUserID,
		Content:        "Hello, World!",
	}
	err := message.Validate()
	assert.Nil(t, err)
}
