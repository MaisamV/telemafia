package entity

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
