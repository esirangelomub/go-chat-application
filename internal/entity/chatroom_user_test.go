// chatroom_user_test.go
package entity

import (
	"github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewChatroomUser(t *testing.T) {
	chatroomID := entity.NewID()
	userID := entity.NewID()

	// Test case 1: Valid chatroom user
	chatroomUser, err := NewChatroomUser(chatroomID, userID)
	assert.Nil(t, err)
	assert.NotNil(t, chatroomUser)
	assert.NotEmpty(t, chatroomUser.ID)
	assert.Equal(t, chatroomID, chatroomUser.ChatroomID)
	assert.Equal(t, userID, chatroomUser.UserID)
}

func TestChatroomUser_Validate(t *testing.T) {
	chatroomID := entity.NewID()
	userID := entity.NewID()

	// Test case 1: Valid chatroom user
	chatroomUser := &ChatroomUser{
		ID:         entity.NewID(),
		ChatroomID: chatroomID,
		UserID:     userID,
	}
	err := chatroomUser.Validate()
	assert.Nil(t, err)
}
