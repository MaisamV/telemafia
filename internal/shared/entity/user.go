package entity

import "fmt"

// UserID represents a unique identifier for a user
type UserID int64

// User represents a user entity
type User struct {
	ID        UserID
	FirstName string
	LastName  string
	Username  string
	Admin     bool

	// Adding TelegramID based on previous observations in room.go placeholder
	TelegramID int64 // Assuming this was the intended ID field
}

// CanCreateRoom checks if the user has permission to create rooms
func (u *User) CanCreateRoom() bool {
	return u.Admin
}

func (u *User) GetProfileLink() string {
	profileLink := ""
	if u.Username != "" {
		profileLink = fmt.Sprintf("[%s](https://t.me/%s)", u.FirstName, u.Username)
	} else {
		profileLink = fmt.Sprintf("[%s](tg://user?id=%d)", u.FirstName, u.ID)
	}
	return profileLink
}
