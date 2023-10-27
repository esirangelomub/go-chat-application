package entity

import "github.com/esirangelomub/go-chat-application/pkg/models/entity"

type Message struct {
	ID             entity.ID    `json:"id"`
	ChatroomUserID entity.ID    `json:"chatroom_user_id"`
	ChatroomUser   ChatroomUser `json:"chatroom_user"`
	Content        string       `json:"content"`
	Timestamp      int64        `json:"timestamp"`
}
