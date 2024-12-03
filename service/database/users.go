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
