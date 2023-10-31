package repository

import (
	"github.com/esirangelomub/go-chat-application/internal/entity"
	entityPkg "github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestChatroomUserRepository(t *testing.T) {
	// Initialize in-memory database
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	db.AutoMigrate(&entity.ChatroomUser{})

	// Create a new ChatroomUser
	chatroomUser, _ := entity.NewChatroomUser(entityPkg.NewID(), entityPkg.NewID())
	chatroomUserDB := NewChatroomUser(db)

	// Test Create method
	createdChatroomUser, err := chatroomUserDB.Create(chatroomUser)
	assert.Nil(t, err)
	assert.Equal(t, chatroomUser.ID, createdChatroomUser.ID)
	assert.Equal(t, chatroomUser.UserID, createdChatroomUser.UserID)
	assert.Equal(t, chatroomUser.ChatroomID, createdChatroomUser.ChatroomID)

	// Test FindByID method
	foundChatroomUser, err := chatroomUserDB.FindByID(chatroomUser.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, chatroomUser.ID, foundChatroomUser.ID)

	// Test Delete method
	err = chatroomUserDB.Delete(chatroomUser)
	assert.Nil(t, err)

	// Verify deletion
	_, err = chatroomUserDB.FindByID(chatroomUser.ID.String())
	assert.Error(t, err)

	// Cleanup
	t.Cleanup(func() {
		db.Migrator().DropTable(&entity.ChatroomUser{})
	})
}
