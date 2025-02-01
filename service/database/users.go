package database

import (
	"strings"

	"github.com/maisto1/WasaText/service/models"
)

func (db *appdbimpl) GetUsers(names string) []models.User {
	users := make([]models.User, 0)
	name := strings.TrimSpace(names)

	rows, err := db.c.Query(`
        SELECT * FROM Users 
        WHERE username LIKE ?`,
		"%"+name+"%",
	)
	if err != nil {
		return users
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.User_id, &user.Username, &user.Photo)
		if err == nil {
			users = append(users, user)
		}
	}

	return users
}

func (db *appdbimpl) EditProfileName(user_id int64, username string) error {
	_, err := db.c.Exec("UPDATE Users SET username = ? WHERE user_id = ?", username, user_id)
	if err != nil {
		return err
	}

	return nil

}

func (db *appdbimpl) EditProfilePhoto(user_id int64, photo []byte) error {
	_, err := db.c.Exec("UPDATE Users SET profile_photo = ? WHERE user_id = ?", photo, user_id)
	if err != nil {
		return err
	}

	return nil
}
