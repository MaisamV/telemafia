package error

var (
	ErrInvalidRoomName     = NewDomainError("ROOM_INVALID_NAME", "room name must be between 3 and 50 characters")
	ErrRoomNotFound        = NewDomainError("ROOM_NOT_FOUND", "room not found")
	ErrRoomAlreadyExists   = NewDomainError("ROOM_ALREADY_EXISTS", "room already exists")
	ErrPlayerNotInRoom     = NewDomainError("ROOM_PLAYER_NOT_IN_ROOM", "player not in room")
	ErrPlayerAlreadyInRoom = NewDomainError("ROOM_PLAYER_ALREADY_IN_ROOM", "player already in room")
)
