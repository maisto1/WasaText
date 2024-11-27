package database

import (
	"database/sql"
	"errors"
)

// If the user exist return userId, otherwise create a new user and return new userID
func (db *appdbimpl) Login(username string) (int64, error) {
	var user_id int64
	// search user in the database
	err := db.c.QueryRow("SELECT ID FROM Users WHERE Username = ?;", username).Scan(&user_id)
	if errors.Is(err, sql.ErrNoRows) {
		// if the user does not exist, create a new user
		_, err = db.c.Exec(`INSERT INTO Users (Username) VALUES (?);`, username)
		if err != nil {
			return 0, err
		}
		// return the id of the new user
		err := db.c.QueryRow("SELECT ID FROM Users WHERE Username = ?;", username).Scan(&user_id)
		if err != nil {
			return 0, err
		}
		return user_id, nil
	} else if err == nil {
		// if the user exists, return the id of the user
		return user_id, nil
	}
	// return an error if the query fails
	return 0, err
}
