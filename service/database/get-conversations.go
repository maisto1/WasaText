package database

import (
	"database/sql"
	"errors"

	"github.com/maisto1/WasaText/service/components/schemas"
)

// Get user preview conversations
func (db *appdbimpl) GetConversations(userId int64) (*schemas.User, error) {
	var user schemas.User

	// search user in the database
	err := db.c.QueryRow("SELECT ID, Username, COALESCE(ProfilePhoto, '') FROM Users WHERE ID = ?;", userId).Scan(&user.ID, &user.Username, &user.ProfilePhoto)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	// Find associated User Conversation
	rows, err := db.c.Query(`
	SELECT 
			c.ID, 
			c.Name, 
			COALESCE(c.ConversationPhoto, ''), 
			uc.LastMessageContent, 
			uc.LastMessageTimestamp
		FROM UserConversations uc
		JOIN Conversations c ON uc.ConversationID = c.ID
		WHERE uc.UserID = ?;
	`, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//Build list user conversation
	user.Conversations = &[]schemas.UserConv{}
	for rows.Next() {
		var conv schemas.UserConv

		var lastMessageContent sql.NullString
		var lastMessageTimestamp sql.NullString

		err = rows.Scan(&conv.ID, &conv.Name, &conv.ConversationPhoto, &lastMessageContent, &lastMessageTimestamp)
		if err != nil {
			return nil, err
		}

		//Build lastmessage only if is presente
		if lastMessageContent.Valid && lastMessageContent.String != "" && lastMessageTimestamp.Valid {
			conv.LastMessage = &schemas.LastMessage{
				Content:   lastMessageContent.String, // Usa il valore della stringa se valido
				Timestamp: lastMessageTimestamp.String,
			}
		} else {
			conv.LastMessage = &schemas.LastMessage{
				Content:   "",
				Timestamp: "",
			}
		}

		*user.Conversations = append(*user.Conversations, conv)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil

}
