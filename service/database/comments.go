package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/maisto1/WasaText/service/models"
)

func (db *appdbimpl) GetComments(user_id int64, conversation_id int64, message_id int64) ([]models.Comment, error) {
	comments := make([]models.Comment, 0)
	var exists bool

	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM Messages WHERE message_id = ?)`, message_id).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("message not found")
	}

	isValid, err := db.CheckUserConversation(user_id, conversation_id)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, errors.New("user is not a partecipant")
	}

	rows, err := db.c.Query(`
	SELECT user_id,comment_id, content, timestamp
	FROM Comments
	WHERE message_id = ?`,
		message_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment
		var user_id int64
		var comment_id int64
		var content string
		var timestamp int64
		var sender models.User

		err = rows.Scan(
			&user_id,
			&comment_id,
			&content,
			&timestamp,
		)
		if err != nil {
			return nil, err
		}
		if rows.Err() != nil {
			return nil, rows.Err()
		}

		comment.Comment_id = comment_id
		comment.Content = content
		comment.Message_id = message_id
		comment.Timestamp = timestamp

		err = db.c.QueryRow(`SELECT * FROM Users WHERE user_id = ?`, user_id).Scan(&sender.User_id, &sender.Username, &sender.Photo)
		if err != nil {
			return nil, err
		}

		comment.Sender = sender

		comments = append(comments, comment)
	}

	return comments, nil
}

func (db *appdbimpl) CreateComment(user_id int64, conversation_id int64, message_id int64, content string) (models.Comment, error) {
	var comment models.Comment
	var exists bool
	current_time := time.Now().Unix()
	var user models.User
	var comment_id int64

	err := db.c.QueryRow(`SELECT EXISTS(SELECT 1 FROM Messages WHERE message_id = ?)`, message_id).Scan(&exists)
	if err != nil {
		return comment, err
	}
	if !exists {
		return comment, errors.New("message not found")
	}

	isValid, err := db.CheckUserConversation(user_id, conversation_id)
	if err != nil {
		return comment, err
	}
	if !isValid {
		return comment, errors.New("user is not a partecipant")
	}

	err = db.c.QueryRow(`
		INSERT INTO Comments (message_id,user_id,content,timestamp)
		VALUES (?,?,?,?)
		RETURNING comment_id`,
		message_id,
		user_id,
		content,
		current_time,
	).Scan(&comment_id)
	if err != nil {
		return comment, err
	}

	err = db.c.QueryRow(`SELECT * FROM Users WHERE user_id = ?`, user_id).Scan(&user.User_id, &user.Username, &user.Photo)
	if err != nil {
		return comment, err
	}
	fmt.Print(content)
	comment.Timestamp = current_time
	comment.Message_id = message_id
	comment.Comment_id = comment_id
	comment.Sender = user
	comment.Content = content

	return comment, nil
}

func (db *appdbimpl) DeleteComment(user_id int64, conversation_id int64, comment_id int64) error {
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
			FROM Comments
			WHERE comment_id = ? AND user_id = ?
		);`,
		comment_id, user_id).Scan(&exists)

	if err != nil {
		return err
	}
	if !exists {
		return errors.New("this comment doesn't belogs from this user")
	}

	_, err = db.c.Exec("DELETE FROM Comments WHERE user_id = ? AND comment_id = ?;", user_id, comment_id)
	if err != nil {
		return err
	}

	return nil
}
