package database

import (
	"github.com/maisto1/WasaText/service/components/schemas"
)

// Create a new conversation
func (db *appdbimpl) SetConversations(userId int64, conversationType string, groupName string, groupPhoto string, partecipantId int64) (*schemas.Conversation, error) {
	var conversation schemas.Conversation

	err := db.c.QueryRow("INSERT INTO Conversations (Name, ConversationType, ConversationPhoto) VALUES (?, ?, ?) RETURNING ID;", groupName, conversationType, groupPhoto).Scan(&conversation.ID)
	if err != nil {
		return nil, err
	}

	if conversationType == "private" {
		_, err = db.c.Exec("INSERT INTO UserConversations (ConversationID, UserID) VALUES (?, ?);", conversation.ID, partecipantId)
		if err != nil {
			return nil, err
		}

		_, err = db.c.Exec("INSERT INTO UserConversations (ConversationID, UserID) VALUES (?, ?);", conversation.ID, userId)
		if err != nil {
			return nil, err
		}

		_, err = db.c.Exec("INSERT INTO ConversationParticipants (ConversationID, UserID) VALUES (?, ?);", conversation.ID, partecipantId)
		if err != nil {
			return nil, err
		}

		conversation.Partecipants = append(conversation.Partecipants, partecipantId)

	}

	_, err = db.c.Exec("INSERT INTO ConversationParticipants (ConversationID, UserID) VALUES (?, ?);", conversation.ID, userId)
	if err != nil {
		return nil, err
	}

	conversation.Name = groupName
	conversation.ConversationType = conversationType
	conversation.ConversationPhoto = groupPhoto
	conversation.Partecipants = append(conversation.Partecipants, userId)

	return &conversation, nil
}
