package models

type Admin struct {
	User
	Username string `json:"username"`
}
