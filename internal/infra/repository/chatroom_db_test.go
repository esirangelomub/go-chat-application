package repository

import (
	"github.com/esirangelomub/go-chat-application/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
	"testing"
)

func TestCreateChatroom(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	assert.Nil(t, err)
	db.AutoMigrate(&entity.Chatroom{})

	chatroomDB := NewChatroom(db)
	chatroom := &entity.Chatroom{Name: "Chatroom1", Description: "Chatroom1 desc"}

	err = chatroomDB.Create(chatroom)
	assert.Nil(t, err)

	var chatroomFound entity.Chatroom
	err = db.First(&chatroomFound, "id = ?", chatroom.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, chatroom.ID, chatroomFound.ID)
	assert.Equal(t, chatroom.Name, chatroomFound.Name)

	t.Cleanup(func() {
		db.Migrator().DropTable(&entity.Chatroom{})
	})
}

func TestFindAllChatRooms(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	assert.Nil(t, err)
	db.AutoMigrate(&entity.Chatroom{})

	chatroomDB := NewChatroom(db)

	for i := 0; i < 11; i++ {
		name := "Chatroom " + strconv.Itoa(i)
		description := "Chatroom " + strconv.Itoa(i) + " desc"

		chatroom, _ := entity.NewChatRoom(name, description)
		err = chatroomDB.Create(chatroom)
		assert.Nil(t, err)
	}

	chatrooms, err := chatroomDB.FindAll(1, 10, "asc")
	assert.Nil(t, err)
	assert.Len(t, chatrooms, 10)

	t.Cleanup(func() {
		db.Migrator().DropTable(&entity.Chatroom{})
	})
}

func TestFindChatroomByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	assert.Nil(t, err)
	db.AutoMigrate(&entity.Chatroom{})

	chatroom, _ := entity.NewChatRoom("Chatroom1", "Chatroom1 desc")
	chatroomDB := NewChatroom(db)

	err = chatroomDB.Create(chatroom)
	assert.Nil(t, err)

	chatroomFound, err := chatroomDB.FindByID(chatroom.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, chatroom.ID, chatroomFound.ID)
	assert.Equal(t, chatroom.Name, chatroomFound.Name)

	t.Cleanup(func() {
		db.Migrator().DropTable(&entity.Chatroom{})
	})
}

func TestUpdateChatroom(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	assert.Nil(t, err)
	db.AutoMigrate(&entity.Chatroom{})

	chatroom, _ := entity.NewChatRoom("Chatroom1", "Chatroom1 desc")
	chatroomDB := NewChatroom(db)

	err = chatroomDB.Create(chatroom)
	assert.Nil(t, err)

	chatroom.Name = "UpdatedChatroom"
	err = chatroomDB.Update(chatroom)
	assert.Nil(t, err)

	chatroomFound, err := chatroomDB.FindByID(chatroom.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, "UpdatedChatroom", chatroomFound.Name)

	t.Cleanup(func() {
		db.Migrator().DropTable(&entity.Chatroom{})
	})
}

func TestDeleteChatroom(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:memory:"), &gorm.Config{})
	assert.Nil(t, err)
	db.AutoMigrate(&entity.Chatroom{})

	chatroom, _ := entity.NewChatRoom("Chatroom1", "Chatroom1 desc")
	chatroomDB := NewChatroom(db)

	err = chatroomDB.Create(chatroom)
	assert.Nil(t, err)

	err = chatroomDB.Delete(chatroom)
	assert.Nil(t, err)

	chatroomFound, err := chatroomDB.FindByID(chatroom.ID.String())
	assert.NotNil(t, err)
	assert.Nil(t, chatroomFound)

	t.Cleanup(func() {
		db.Migrator().DropTable(&entity.Chatroom{})
	})
}
