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
	UniqueConfirmAssignments       = "confirm_assignments"
	UniqueShowMyRole               = "show_my_role"
	UniqueGetInviteLink            = "get_invite_link"
	UniqueCreateGameSelectRoom     = "cg_room"
	UniqueCreateGameSelectScenario = "cg_scen"
	UniqueStartGame                = "cg_start"
	UniqueCancelGame               = "cancel_cg"

	// Kick User Flow
	UniqueKickUserSelect  = "kick_user_select"  // Shows the list of users to kick
	UniqueKickUserConfirm = "kick_user_confirm" // Confirms kicking the selected user

	// Generic Cancel (might need context)
	UniqueCancel = "cancel"
)
