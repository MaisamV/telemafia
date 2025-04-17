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
}

// CanCreateRoom checks if the user has permission to create rooms
func (u *User) CanCreateRoom() bool {
	return u.Admin
}
