package models

import "net/mail"

// User is a user of the application
type User struct {
	ID           int64
	Email        string
	DisplayName  string
	Name         string
	IsAdmin      bool
	PasswordHash string
	IsActive     bool
}

// NewUser represents the data needed by the storage layer to create a new user
type NewUser struct {
	Email        string
	DisplayName  string
	Name         string
	PasswordHash string
}

type NewUserRequest struct {
	Email           string
	DisplayName     string
	Name            string
	Password        string
	ConfirmPassword string
}

// Validate validates a new user request
func (u NewUserRequest) Validate() error {
	ve := ValidationErrors{}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		ve.Add("email", "invalid email address")
	}
	if len(u.DisplayName) < 3 || len(u.DisplayName) > 50 {
		ve.Add("display_name", "display name must be between 3 and 50 characters long")
	}
	if len(u.Name) < 3 || len(u.Name) > 50 {
		ve.Add("name", "name must be between 3 and 50 characters long")
	}
	if len(u.Password) < 8 || len(u.Password) > 64 {
		ve.Add("password", "password must be between 8 and 64 characters long")
	}
	if u.Password != u.ConfirmPassword {
		ve.Add("confirm_password", "passwords do not match")
	}

	if ve.Any() {
		return ve
	}
	return nil
}

// UserFieldsTaken represents fields that are taken when creating a user
type UserFieldsTaken struct {
	EmailTaken       bool
	DisplayNameTaken bool
}
