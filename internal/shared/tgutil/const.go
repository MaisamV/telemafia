package tgutil

// Unique identifiers for inline buttons
const (
	// Join/Leave related
	UniqueJoinRoom            = "join_room"
	UniqueLeaveRoomSelectRoom = "leave_room_select"
	UniqueLeaveRoomConfirm    = "leave_room_confirm"

	// Delete Room related
	UniqueDeleteRoomSelectRoom = "delete_room_select"
	UniqueDeleteRoomConfirm    = "delete_room_confirm"

	// Game/Assignment related
	UniqueConfirmAssignments = "confirm_assignments"
	// UniqueShowMyRole          = "show_my_role" // Placeholder if needed later

	// Generic Cancel (might need context)
	UniqueCancel = "cancel"
)
