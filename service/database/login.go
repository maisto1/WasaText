package database

import (
	"database/sql"
	"errors"
)

// If the user exist return userId, otherwise create a new user and return new userID
func (db *appdbimpl) Login(username string) (int64, error) {
	var user_id int64

	//search user in database
	err := db.c.QueryRow("SELECT user_id FROM Users WHERE username = ?;", username).Scan(&user_id)
	if errors.Is(err, sql.ErrNoRows) {
		err := db.c.QueryRow(`INSERT INTO users (username) VALUES (?) RETURNING user_id;`, username).Scan(&user_id)
		if err != nil {
			return 0, err
		}
		return user_id, nil
	} else if err == nil {
		return user_id, nil
	}
	return 0, err
}
