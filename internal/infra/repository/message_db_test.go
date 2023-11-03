package repository

import (
	"github.com/esirangelomub/go-chat-application/internal/entity"
	entityPkg "github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestMessageRepository(t *testing.T) {
	// Initialize in-memory database
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.Message{})

	// Create a new Message
	message, _ := entity.NewMessage(entityPkg.NewID(), entityPkg.NewID(), "Hello World!")
	messageDB := NewMessage(db)

	// Test Create method
	err = messageDB.Create(message)
	assert.Nil(t, err)

	// Verify creation
	var createdMessage entity.Message
	err = db.First(&createdMessage, "id = ?", message.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, message.ID, createdMessage.ID)
	assert.Equal(t, message.ChatroomID, createdMessage.ChatroomID)
	assert.Equal(t, message.UserID, createdMessage.UserID)
	assert.Equal(t, message.Content, createdMessage.Content)

	// Cleanup
	t.Cleanup(func() {
		db.Migrator().DropTable(&entity.Message{})
	})
}
