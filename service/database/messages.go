package database

import (
	"errors"

	"github.com/maisto1/WasaText/service/models"
)

func (db *appdbimpl) GetMessages(conversation_id int64) ([]models.Message, error) {
	messages := make([]models.Message, 0)

	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Conversations WHERE conversation_id = ?)", conversation_id).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("conversation not found")
	}

	rows, err := db.c.Query(`
		SELECT message_id, timestamp, user_id, type, content, media, status, forwarded
		FROM Messages
		WHERE conversation_id = ?
	`, conversation_id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message_id int64
		var timestamp int64
		var user_id int64
		var typeMedia string
		var content string
		var media []byte
		var status string
		var forwarded bool
		var message models.Message
		var sender models.User

		err = rows.Scan(
			&message_id,
			&timestamp,
			&user_id,
			&typeMedia,
			&content,
			&media,
			&status,
			&forwarded,
		)
		if err != nil {
			return messages, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}

		message.Message_id = message_id
		message.Timestamp = timestamp
		message.Type = typeMedia
		message.Content = content
		message.Media = media
		message.Status = status
		message.Forwarded = forwarded

		err = db.c.QueryRow(`SELECT * FROM Users WHERE user_id = ?`, user_id).Scan(&sender.User_id, &sender.Username, &sender.Photo)
		if err != nil {
			return messages, err
		}

		message.Sender = sender

		messages = append(messages, message)
	}

	return messages, nil
}
