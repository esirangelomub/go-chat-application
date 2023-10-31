package entity

import (
	"github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewChatRoom(t *testing.T) {
	// Test case 1: Successful creation of a chatroom
	chatroom, err := NewChatRoom("General", "A general chatroom for everyone.")
	assert.Nil(t, err)              // No error should be returned
	assert.NotNil(t, chatroom)      // Chatroom object should not be nil
	assert.NotEmpty(t, chatroom.ID) // ID should not be empty
	assert.Equal(t, "General", chatroom.Name)
	assert.Equal(t, "A general chatroom for everyone.", chatroom.Description)

	// Test case 2: Error when name is empty
	chatroom, err = NewChatRoom("", "Description but no name.")
	assert.Error(t, err)                    // Error should be returned
	assert.Equal(t, ErrNameIsRequired, err) // Error should be specific
	assert.Nil(t, chatroom)                 // Chatroom object should be nil
}

func TestChatroom_Validate(t *testing.T) {
	// Test case 1: Valid chatroom
	chatroom := &Chatroom{
		ID:          entity.NewID(),
		Name:        "General",
		Description: "A general chatroom for everyone.",
	}
	err := chatroom.Validate()
	assert.Nil(t, err)

	chatroom.ID = entity.NewID()
	chatroom.Name = ""
	err = chatroom.Validate()
	assert.Equal(t, ErrNameIsRequired, err) // Error should be specific
}
