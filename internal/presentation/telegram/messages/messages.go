package messages

// Messages holds all user-facing strings, loaded from messages.json
type Messages struct {
	Common   CommonMessages   `json:"common"`
	Room     RoomMessages     `json:"room"`
	Scenario ScenarioMessages `json:"scenario"`
	Game     GameMessages     `json:"game"`
	Refresh  RefreshMessages  `json:"refresh"`
}

type CommonMessages struct {
	Help                   string `json:"help"`
	ErrorGeneric           string `json:"error_generic"`
	ErrorIdentifyUser      string `json:"error_identify_user"`
	ErrorIdentifyRequester string `json:"error_identify_requester"`
	ErrorPreparingContent  string `json:"error_preparing_content"`
	ErrorCommandUsage      string `json:"error_command_usage"`
	ErrorPermissionDenied  string `json:"error_permission_denied"`
	CallbackErrorGeneric   string `json:"callback_error_generic"`
	CallbackCancelled      string `json:"callback_cancelled"`
	CallbackFailedEdit     string `json:"callback_failed_edit"`
	CallbackFailedRespond  string `json:"callback_failed_respond"`
}

type RoomMessages struct {
	CreatePrompt                string `json:"create_prompt"`
	CreateSuccess               string `json:"create_success"`
	CreateError                 string `json:"create_error"`
	JoinPrompt                  string `json:"join_prompt"`
	JoinSuccess                 string `json:"join_success"`
	JoinError                   string `json:"join_error"`
	JoinButtonText              string `json:"join_button_text"`
	RoomDetail                  string `json:"room_detail"`
	RoomDetailWithScenario      string `json:"room_detail_with_scenario"`
	LeavePrompt                 string `json:"leave_prompt"`
	LeaveSuccess                string `json:"leave_success"`
	LeaveError                  string `json:"leave_error"`
	LeaveConfirmPrompt          string `json:"leave_confirm_prompt"`
	LeaveConfirmButton          string `json:"leave_confirm_button"`
	LeaveCancelButton           string `json:"leave_cancel_button"`
	LeaveButton                 string `json:"leave_button"`
	LeaveCallbackSuccess        string `json:"leave_callback_success"`
	LeaveCallbackEditSuccess    string `json:"leave_callback_edit_success"`
	LeaveCallbackEditFail       string `json:"leave_callback_edit_fail"`
	LeaveCallbackMultiroomError string `json:"leave_callback_multiroom_error"`
	LeaveCallbackMultiroomEdit  string `json:"leave_callback_multiroom_edit"`
	LeaveCallbackNoroomError    string `json:"leave_callback_noroom_error"`
	LeaveCallbackNoroomEdit     string `json:"leave_callback_noroom_edit"`
	KickPrompt                  string `json:"kick_prompt"`
	KickInvalidUserID           string `json:"kick_invalid_user_id"`
	KickSuccess                 string `json:"kick_success"`
	KickError                   string `json:"kick_error"`
	DeletePromptSelect          string `json:"delete_prompt_select"`
	DeletePromptConfirm         string `json:"delete_prompt_confirm"`
	DeleteNoRooms               string `json:"delete_no_rooms"`
	DeleteErrorFetch            string `json:"delete_error_fetch"`
	DeleteConfirmButton         string `json:"delete_confirm_button"`
	DeleteCancelButton          string `json:"delete_cancel_button"`
	DeleteCallbackSuccess       string `json:"delete_callback_success"`
	DeleteCallbackEditSuccess   string `json:"delete_callback_edit_success"`
	DeleteCallbackEditFail      string `json:"delete_callback_edit_fail"`
	DeleteCallbackError         string `json:"delete_callback_error"`
	ListTitle                   string `json:"list_title"`
	ListNoRooms                 string `json:"list_no_rooms"`
	ListError                   string `json:"list_error"`
	ListErrorPrepare            string `json:"list_error_prepare"`
	MyRoomsTitle                string `json:"my_rooms_title"`
	MyRoomsNone                 string `json:"my_rooms_none"`
	MyRoomsError                string `json:"my_rooms_error"`
}

type ScenarioMessages struct {
	CreatePrompt      string `json:"create_prompt"`
	CreateSuccess     string `json:"create_success"`
	CreateError       string `json:"create_error"`
	DeletePrompt      string `json:"delete_prompt"`
	DeleteSuccess     string `json:"delete_success"`
	DeleteError       string `json:"delete_error"`
	AddRolePrompt     string `json:"add_role_prompt"`
	AddRoleSuccess    string `json:"add_role_success"`
	AddRoleError      string `json:"add_role_error"`
	RemoveRolePrompt  string `json:"remove_role_prompt"`
	RemoveRoleSuccess string `json:"remove_role_success"`
	RemoveRoleError   string `json:"remove_role_error"`
}

type GameMessages struct {
	AssignScenarioPrompt                string `json:"assign_scenario_prompt"`
	AssignScenarioSuccess               string `json:"assign_scenario_success"`
	AssignScenarioErrorRoomFind         string `json:"assign_scenario_error_room_find"`
	AssignScenarioErrorRoomNotfound     string `json:"assign_scenario_error_room_notfound"`
	AssignScenarioErrorScenarioFind     string `json:"assign_scenario_error_scenario_find"`
	AssignScenarioErrorScenarioNotfound string `json:"assign_scenario_error_scenario_notfound"`
	AssignScenarioErrorUpdateRoom       string `json:"assign_scenario_error_update_room"`
	AssignScenarioErrorGameCreate       string `json:"assign_scenario_error_game_create"`
	AssignRolesPrompt                   string `json:"assign_roles_prompt"`
	AssignRolesSuccessPublic            string `json:"assign_roles_success_public"`
	AssignRolesSuccessPrivate           string `json:"assign_roles_success_private"`
	AssignRolesError                    string `json:"assign_roles_error"`
	AssignRolesErrorGameFind            string `json:"assign_roles_error_game_find"`
	AssignRolesErrorNoScenario          string `json:"assign_roles_error_no_scenario"`
	AssignRolesErrorNoRoom              string `json:"assign_roles_error_no_room"`
	AssignRolesErrorPlayerMismatch      string `json:"assign_roles_error_player_mismatch"`
	AssignRolesErrorFetchingPlayers     string `json:"assign_roles_error_fetching_players"`
	AssignRolesErrorUpdatingGame        string `json:"assign_roles_error_updating_game"`
	AssignRolesErrorSendingPrivate      string `json:"assign_roles_error_sending_private"`
	ListGamesTitle                      string `json:"list_games_title"`
	ListGamesEntry                      string `json:"list_games_entry"`
	ListGamesNoGames                    string `json:"list_games_no_games"`
	ListGamesError                      string `json:"list_games_error"`
	AssignmentsConfirmButton            string `json:"assignments_confirm_button"`
	AssignmentsConfirmedResponse        string `json:"assignments_confirmed_response"`
}

type RefreshMessages struct {
	ErrorPrepare          string `json:"error_prepare"`
	ErrorEdit             string `json:"error_edit"`
	ErrorEditRemoving     string `json:"error_edit_removing"`
	ErrorSendNew          string `json:"error_send_new"`
	LogUpdateSuccess      string `json:"log_update_success"`
	LogEditFailSendingNew string `json:"log_edit_fail_sending_new"`
	LogSendNewSuccess     string `json:"log_send_new_success"`
	LogRemovedUser        string `json:"log_removed_user"`
}
