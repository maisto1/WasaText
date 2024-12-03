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
	WITH OtherUser AS (
		SELECT 
			p1.conversation_id,
			u.user_id AS other_user_id,
			u.username AS other_username,
			u.profile_photo AS other_photo
		FROM Partecipants p1
		JOIN Partecipants p2 ON p1.conversation_id = p2.conversation_id
		JOIN Users u ON u.user_id = p2.user_id
		WHERE p1.user_id = ? AND p2.user_id != ?
	)
	SELECT 
		c.conversation_id,
		COALESCE(c.group_name, o.other_username, '') AS group_name,
		COALESCE(c.group_photo, o.other_photo, '') AS group_photo
	FROM Conversations c
	LEFT JOIN OtherUser o ON c.conversation_id = o.conversation_id
	WHERE c.conversation_id IN (
		SELECT conversation_id 
		FROM Partecipants 
		WHERE user_id = ?
	);

	`, user_id, user_id, user_id)

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
func (db *appdbimpl) CreateConversation(user_id int64, group_name string, typeConv string, partecipant string) error {
	var conversation_id int64
	var user_partecipant []models.User

	if typeConv == "private" {

		user_partecipant = db.GetUsers(partecipant)
		if len(user_partecipant) == 0 {
			return errors.New("partecipant not found")
		}
		partecipant_id := user_partecipant[0].User_id

		isVAlid, err := db.CheckPrivateConversation(user_id, partecipant_id)
		if err != nil {
			return err
		}

		if isVAlid {
			return errors.New("already have a conversation")
		}
	}

	err := db.c.QueryRow(`INSERT INTO Conversations (type) VALUES (?) RETURNING conversation_id;`, typeConv).Scan(&conversation_id)
	if err != nil {
		return err
	}

	_, err = db.c.Exec(`INSERT INTO Partecipants (user_id,conversation_id) VALUES (?,?);`, user_id, conversation_id)
	if err != nil {
		return err
	}

	if typeConv == "private" {
		_, err = db.c.Exec(`INSERT INTO Partecipants (user_id,conversation_id) VALUES (?,?);`, user_partecipant[0].User_id, conversation_id)
		if err != nil {
			return err
		}
	} else {
		_, err = db.c.Exec(`UPDATE Conversations SET group_name = ? WHERE conversation_id = ?;`, group_name, conversation_id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *appdbimpl) CheckPrivateConversation(user_id1, user_id2 int64) (bool, error) {
	var count int

	err := db.c.QueryRow(`
		SELECT COUNT(*)
		FROM Conversations c
		JOIN Partecipants p1 ON c.conversation_id = p1.conversation_id
		JOIN Partecipants p2 ON c.conversation_id = p2.conversation_id
		WHERE c.type = 'private'
		AND ((p1.user_id = ? AND p2.user_id = ?)
		OR (p1.user_id = ? AND p2.user_id = ?))
	`, user_id1, user_id2, user_id2, user_id1).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
