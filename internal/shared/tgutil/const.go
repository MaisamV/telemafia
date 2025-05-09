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
	UniqueCreateGameSelectRoom     = "cg_room"     // -> Select Scenario
	UniqueCreateGameSelectScenario = "cg_scen"     // -> Show Confirmation (Start/Choose/Cancel)
	UniqueStartGame                = "cg_start"    // -> Assign roles directly
	UniqueChooseCardStart          = "cg_choose"   // -> Start interactive card selection
	UniquePlayerSelectsCard        = "cg_sel_card" // Player clicks a numbered card
	UniqueCancelGame               = "cancel_cg"

	// Kick User Flow
	UniqueKickUserSelect  = "kick_user_select"  // Shows the list of users to kick
	UniqueKickUserConfirm = "kick_user_confirm" // Confirms kicking the selected user

	// Change Moderator Flow
	UniqueChangeModeratorSelect  = "mod_user_select"  // Shows the list of users to make moderator
	UniqueChangeModeratorConfirm = "mod_user_confirm" // Confirms setting the selected user as moderator

	// Common
	UniqueCancel = "cancel"
)
