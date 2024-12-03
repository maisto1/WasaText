package database

import (
	"database/sql"
	"errors"

	"github.com/maisto1/WasaText/service/models"
)

// Get preview Conversations
func (db *appdbimpl) GetPreviewConversations(user_id int64) ([]models.Preview, error) {
	previews := make([]models.Preview, 0)

	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE user_id = ?)", user_id).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("user not found in database")
	}

	rows, err := db.c.Query(`
		SELECT Conversations.conversation_id as Conversations_id,  COALESCE(p.conversation_name, Conversations.name) AS name, 
        COALESCE(p.conversation_photo, Conversations.photo) AS photo
		FROM Partecipants p
		JOIN Conversations ON Conversations.conversation_id = p.conversation_id
		WHERE p.user_id = ?
	`, user_id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var conversation_id int64
		var name string
		var photo []byte
		var preview models.Preview

		err = rows.Scan(&conversation_id, &name, &photo)
		if err != nil {
			return previews, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}

		preview.Conversation_id = conversation_id
		preview.Name = name
		preview.Photo = photo

		latestMessage, err := db.GetLatestMessage(conversation_id)
		if err != nil {
			preview.LatestMessage = nil
		} else {
			preview.LatestMessage = &latestMessage
		}

		previews = append(previews, preview)
	}

	return previews, nil
}

// Get latest message in a specific conversation
func (db *appdbimpl) GetLatestMessage(conversation_id int64) (models.Message, error) {
	var message models.Message
	var user models.User
	var user_id int64

	err := db.c.QueryRow(`
		SELECT 
			m.message_id,
			m.content,
			m.media,
			m.type,
			m.timestamp,
			m.status,
			m.isForwarded,
			m.user_id
		FROM 
			Messages m
		WHERE 
			m.conversation_id = ?
		ORDER BY 
			m.timestamp DESC
		LIMIT 1;
	`, conversation_id).Scan(
		&message.Message_id,
		&message.Content,
		&message.Media,
		&message.Type,
		&message.Timestamp,
		&message.Status,
		&message.Forwarded,
		&user_id,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return message, err
		} else {
			return message, err
		}
	}

	err = db.c.QueryRow(`SELECT * FROM Users WHERE user_id = ?`, user_id).Scan(&user.User_id, &user.Username, &user.Photo)
	if err != nil {
		return message, err
	}

	message.Sender = user

	return message, nil
}

// Create a new conevrsation
// func (db *appdbimpl) CreateConversation(user_id int64, name string, photo []byte, typeConv string) (bool, error) {

// 	err := db.c.QueryRow(`INSERT INTO Conversations (name,photo,type) VALUES (?) RETURNING conversation_id;`, user_id, name, photo).Scan(&user_id)

// }
