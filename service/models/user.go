package models

type User struct {
	User_id  int64  `json:"id"`
	Username string `json:"username"`
	Photo    []byte `json:"profilePhoto"`
}
