package database

import (
	"strings"

	"github.com/maisto1/WasaText/service/models"
)

func (db *appdbimpl) GetUsers(names string) []models.User {
	users := make([]models.User, 0)

	usernameList := strings.Split(names, ",")

	for _, name := range usernameList {
		var user models.User
		name = strings.TrimSpace(name)
		err := db.c.QueryRow("SELECT * FROM Users WHERE username = ?;", name).Scan(&user.User_id, &user.Username, &user.Photo)
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
