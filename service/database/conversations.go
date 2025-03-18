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
        COALESCE(c.name, o.other_username, '') AS conversation_name,
        COALESCE(c.conversation_photo, o.other_photo, '') AS conversation_photo,
        c.conversation_type
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
		var conversationType string
		var preview models.Preview

		err = rows.Scan(&conversation_id, &name, &photo, &conversationType)
		if err != nil {
			return previews, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}

		preview.Conversation_id = conversation_id
		preview.Name = name
		preview.Photo = photo
		preview.ConversationType = conversationType

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

// Create a new conversation
func (db *appdbimpl) CreateConversation(user_id int64, name string, typeConv string, partecipant string) (int, error) {
	var conversation_id int64
	var user_partecipant []models.User

	// Validazione del tipo di conversazione
	if typeConv != "private" && typeConv != "group" {
		return 0, errors.New("invalid conversation type")
	}

	if typeConv == "private" {
		user_partecipant = db.GetUsers(partecipant)
		if len(user_partecipant) == 0 {
			return 0, errors.New("participant not found")
		}
		participant_id := user_partecipant[0].User_id

		isValid, err := db.CheckPrivateConversation(user_id, participant_id)
		if err != nil {
			return 0, err
		}
		if isValid {
			return 0, errors.New("conversation already exists")
		}
	}

	// Per le conversazioni di gruppo, usa il nome passato
	// Per le conversazioni private, usa il nome del partecipante
	if typeConv == "private" && name == "" {
		name = user_partecipant[0].Username
	}

	// Inserisci la conversazione con nome e tipo
	err := db.c.QueryRow(
		`INSERT INTO Conversations (name, conversation_type) VALUES (?, ?) RETURNING conversation_id;`,
		name,
		typeConv,
	).Scan(&conversation_id)
	if err != nil {
		return 0, err
	}

	// Inserisci il primo partecipante (creatore)
	_, err = db.c.Exec(
		`INSERT INTO Partecipants (user_id, conversation_id) VALUES (?, ?);`,
		user_id,
		conversation_id,
	)
	if err != nil {
		return 0, err
	}

	// Per conversazioni private, aggiungi l'altro partecipante
	if typeConv == "private" {
		_, err = db.c.Exec(
			`INSERT INTO Partecipants (user_id, conversation_id) VALUES (?, ?);`,
			user_partecipant[0].User_id,
			conversation_id,
		)
		if err != nil {
			return 0, err
		}
	}

	return int(conversation_id), nil
}
func (db *appdbimpl) CheckPrivateConversation(user_id1, user_id2 int64) (bool, error) {
	var count int

	err := db.c.QueryRow(`
        SELECT COUNT(*)
        FROM Conversations c
        JOIN Partecipants p1 ON c.conversation_id = p1.conversation_id
        JOIN Partecipants p2 ON c.conversation_id = p2.conversation_id
        WHERE c.conversation_type = 'private'
        AND ((p1.user_id = ? AND p2.user_id = ?)
        OR (p1.user_id = ? AND p2.user_id = ?))
    `, user_id1, user_id2, user_id2, user_id1).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
