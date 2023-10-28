package service

import (
	"github.com/esirangelomub/go-chat-application/internal/entity"
	"github.com/esirangelomub/go-chat-application/internal/infra/repository"
	entityPkg "github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"log"
)

type ChatService struct {
	ChatroomUserDB  repository.ChatRoomUserInterface
	MessageDB       repository.MessageInterface
	saveMessageChan chan *entityPkg.ChatMessage
}

func NewChatService(cu repository.ChatRoomUserInterface, m repository.MessageInterface) *ChatService {
	service := &ChatService{
		ChatroomUserDB:  cu,
		MessageDB:       m,
		saveMessageChan: make(chan *entityPkg.ChatMessage),
	}

	//go service.listenAndSaveMessages()

	return service
}

func (s *ChatService) listenAndSaveMessages() {
	for msg := range s.saveMessageChan {
		chatRoomUser, err := entity.NewChatroomUser(msg.ChatroomID, msg.UserID)
		if err != nil {
			continue
		}
		log.Printf("ChatroomUser: %v", chatRoomUser)
		err = s.ChatroomUserDB.Create(chatRoomUser)
		if err != nil {
			continue
		}

		message, err := entity.NewMessage(chatRoomUser.ID, msg.Content)
		if err != nil {
			if delErr := s.ChatroomUserDB.Delete(chatRoomUser); delErr != nil {
				log.Printf("Failed to delete chatUser after message creation error: %v", delErr)
			}
			continue
		}

		err = s.MessageDB.Create(message)
		if err != nil {
			if delErr := s.ChatroomUserDB.Delete(chatRoomUser); delErr != nil {
				log.Printf("Failed to delete chatUser after message creation error: %v", delErr)
			}
			continue
		}
	}
}

func (s *ChatService) QueueMessageForPersistence(msg *entityPkg.ChatMessage) {
	s.saveMessageChan <- msg
}
