package database

import (
	"errors"
	"time"

	"github.com/maisto1/WasaText/service/models"
)

func (db *appdbimpl) CheckUserConversation(user_id int64, conversation_id int64) (bool, error) {
	var exists bool
	err := db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Conversations WHERE conversation_id = ?)", conversation_id).Scan(&exists)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, errors.New("conversation not found")
	}

	var isParticipant bool
	err = db.c.QueryRow("SELECT EXISTS(SELECT 1 FROM Partecipants WHERE user_id = ? AND conversation_id = ?)", user_id, conversation_id).Scan(&isParticipant)
	if err != nil {
		return false, err
	}
	if !isParticipant {
		return false, errors.New("user is not a partecipant")
	}

	return true, nil
}

func (db *appdbimpl) GetMessages(user_id int64, conversation_id int64) ([]models.Message, error) {
	messages := make([]models.Message, 0)

	isValid, err := db.CheckUserConversation(user_id, conversation_id)
	if err != nil {
		return messages, err
	}
	if !isValid {
		return messages, errors.New("user is not a partecipant")
	}

	rows, err := db.c.Query(`
		SELECT message_id, timestamp, user_id, type, content, media, status, isForwarded
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

func (db *appdbimpl) CreateMessage(user_id int64, conversation_id int64, typeMessage string, content string, media []byte) (models.Message, error) {
	var message_id int64
	var message models.Message
	var user models.User

	current_time := time.Now().Unix()

	isValid, err := db.CheckUserConversation(user_id, conversation_id)
	if err != nil {
		return message, err
	}
	if !isValid {
		return message, errors.New("user is not a partecipant")
	}

	err = db.c.QueryRow(`
		INSERT INTO Messages (conversation_id,user_id,content,media,type,timestamp,status,isForwarded) 
		VALUES (?,?,?,?,?,?,?,?) RETURNING message_id;`,
		conversation_id,
		user_id,
		content,
		media,
		typeMessage,
		current_time,
		"sent",
		false,
	).Scan(&message_id)
	if err != nil {
		return message, err
	}

	err = db.c.QueryRow(`SELECT * FROM Users WHERE user_id = ?`, user_id).Scan(&user.User_id, &user.Username, &user.Photo)
	if err != nil {
		return message, err
	}

	message.Message_id = message_id
	message.Timestamp = current_time
	message.Sender = user
	message.Type = typeMessage
	message.Content = content
	message.Media = media
	message.Status = "sent"
	message.Forwarded = false

	return message, nil
}

func (db *appdbimpl) DeleteMessage(user_id int64, conversation_id int64, message_id int64) error {
	isValid, err := db.CheckUserConversation(user_id, conversation_id)
	if err != nil {
		return err
	}
	if !isValid {
		return errors.New("user is not a partecipant")
	}

	var exists bool
	err = db.c.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM Messages
			WHERE message_id = ? AND user_id = ?
		);`,
		message_id, user_id).Scan(&exists)

	if err != nil {
		return err
	}
	if !exists {
		return errors.New("this message doesn't belogs from this user")
	}

	_, err = db.c.Exec("DELETE FROM Messages WHERE user_id = ? AND message_id = ?;", user_id, message_id)
	if err != nil {
		return err
	}

	return nil
}
