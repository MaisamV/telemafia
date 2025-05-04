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
	CreatePrompt                   string `json:"create_prompt"`
	CreateSuccess                  string `json:"create_success"`
	CreateError                    string `json:"create_error"`
	RoomDetail                     string `json:"room_detail"`
	RoomDetailWithScenario         string `json:"room_detail_with_scenario"`
	JoinPrompt                     string `json:"join_prompt"`
	JoinSuccess                    string `json:"join_success"`
	JoinError                      string `json:"join_error"`
	JoinButtonText                 string `json:"join_button_text"`
	LeavePrompt                    string `json:"leave_prompt"`
	LeaveSuccess                   string `json:"leave_success"`
	LeaveError                     string `json:"leave_error"`
	LeaveConfirmPrompt             string `json:"leave_confirm_prompt"`
	LeaveConfirmButton             string `json:"leave_confirm_button"`
	LeaveCancelButton              string `json:"leave_cancel_button"`
	LeaveButton                    string `json:"leave_button"`
	LeaveCallbackSuccess           string `json:"leave_callback_success"`
	LeaveCallbackEditSuccess       string `json:"leave_callback_edit_success"`
	LeaveCallbackEditFail          string `json:"leave_callback_edit_fail"`
	LeaveCallbackMultiroomError    string `json:"leave_callback_multiroom_error"`
	LeaveCallbackMultiroomEdit     string `json:"leave_callback_multiroom_edit"`
	LeaveCallbackNoroomError       string `json:"leave_callback_noroom_error"`
	LeaveCallbackNoroomEdit        string `json:"leave_callback_noroom_edit"`
	RoomNotFound                   string `json:"RoomNotFound"`
	InviteLinkButton               string `json:"InviteLinkButton"`
	InviteLinkResponse             string `json:"InviteLinkResponse"`
	KickPrompt                     string `json:"KickPrompt"`
	KickInvalidUserID              string `json:"kick_invalid_user_id"`
	KickSuccess                    string `json:"kick_success"`
	KickError                      string `json:"kick_error"`
	DeletePromptSelect             string `json:"delete_prompt_select"`
	DeletePromptConfirm            string `json:"delete_prompt_confirm"`
	DeleteNoRooms                  string `json:"delete_no_rooms"`
	DeleteErrorFetch               string `json:"delete_error_fetch"`
	DeleteConfirmButton            string `json:"delete_confirm_button"`
	DeleteCancelButton             string `json:"delete_cancel_button"`
	DeleteCallbackSuccess          string `json:"delete_callback_success"`
	DeleteCallbackEditSuccess      string `json:"delete_callback_edit_success"`
	DeleteCallbackEditFail         string `json:"delete_callback_edit_fail"`
	DeleteCallbackError            string `json:"delete_callback_error"`
	ListTitle                      string `json:"list_title"`
	ListNoRooms                    string `json:"list_no_rooms"`
	ListError                      string `json:"list_error"`
	ListErrorPrepare               string `json:"list_error_prepare"`
	MyRoomsTitle                   string `json:"my_rooms_title"`
	MyRoomsNone                    string `json:"my_rooms_none"`
	MyRoomsError                   string `json:"my_rooms_error"`
	KickUserButton                 string `json:"KickUserButton"`
	KickUserSelectPrompt           string `json:"KickUserSelectPrompt"`
	KickUserConfirmPrompt          string `json:"KickUserConfirmPrompt"`
	KickUserCallbackSuccess        string `json:"KickUserCallbackSuccess"`
	KickUserCallbackError          string `json:"KickUserCallbackError"`
	KickUserNoPlayers              string `json:"KickUserNoPlayers"`
	ChangeModeratorButton          string `json:"ChangeModeratorButton"`
	ChangeModeratorSelectPrompt    string `json:"ChangeModeratorSelectPrompt"`
	ChangeModeratorCallbackSuccess string `json:"ChangeModeratorCallbackSuccess"`
	ChangeModeratorCallbackError   string `json:"ChangeModeratorCallbackError"`
	ChangeModeratorNoCandidates    string `json:"ChangeModeratorNoCandidates"`
}

type ScenarioMessages struct {
	CreatePrompt                   string `json:"create_prompt"`
	CreateSuccess                  string `json:"create_success"`
	CreateError                    string `json:"create_error"`
	DeletePrompt                   string `json:"delete_prompt"`
	DeleteSuccess                  string `json:"delete_success"`
	DeleteError                    string `json:"delete_error"`
	AddScenarioJSONPrompt          string `json:"add_scenario_json_prompt"`
	AddScenarioJSONSuccess         string `json:"add_scenario_json_success"`
	AddScenarioJSONInvalidJSON     string `json:"add_scenario_json_invalid_json"`
	AddScenarioJSONValidationError string `json:"add_scenario_json_validation_error"`
	AddScenarioJSONErrorGeneric    string `json:"add_scenario_json_error_generic"`
}

type GameMessages struct {
	AssignScenarioSuccess               string `json:"assign_scenario_success"`
	AssignScenarioErrorRoomFind         string `json:"assign_scenario_error_room_find"`
	AssignScenarioErrorRoomNotFound     string `json:"assign_scenario_error_room_notfound"`
	AssignScenarioErrorScenarioFind     string `json:"assign_scenario_error_scenario_find"`
	AssignScenarioErrorScenarioNotFound string `json:"assign_scenario_error_scenario_notfound"`
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
	CreateGameSelectRoomPrompt          string `json:"create_game_select_room_prompt"`
	CreateGameSelectScenarioPrompt      string `json:"create_game_select_scenario_prompt"`
	CreateGameConfirmPrompt             string `json:"create_game_confirm_prompt"`
	CreateGameStartedSuccess            string `json:"create_game_started_success"`
	CreateGameErrorFetchRooms           string `json:"create_game_error_fetch_rooms"`
	CreateGameErrorFetchScenarios       string `json:"create_game_error_fetch_scenarios"`
	CreateGameErrorFetchPlayers         string `json:"create_game_error_fetch_players"`
	CreateGameErrorFetchScenarioDetails string `json:"create_game_error_fetch_scenario_details"`
	CreateGameErrorCreatingGame         string `json:"create_game_error_creating_game"`
	CreateGameErrorAssigningRoles       string `json:"create_game_error_assigning_roles"`
	CreateGameStartButton               string `json:"create_game_start_button"`
	CreateGameCancelButton              string `json:"create_game_cancel_button"`
	SelectRoomPrompt                    string `json:"SelectRoomPrompt"`
	SelectRoomButton                    string `json:"SelectRoomButton"`
	GameCreatedSuccess                  string `json:"GameCreatedSuccess"`
	GameAlreadyExists                   string `json:"GameAlreadyExists"`
	AssignRolesButton                   string `json:"AssignRolesButton"`
	RolesAssignedSuccess                string `json:"RolesAssignedSuccess"`
	RoleAssignmentPM                    string `json:"RoleAssignmentPM"`
	ErrorAssignRolesPlayerCount         string `json:"ErrorAssignRolesPlayerCount"`
	ErrorAssignRolesNoScenario          string `json:"ErrorAssignRolesNoScenario"`
	ErrorAssignRolesGameNotFound        string `json:"ErrorAssignRolesGameNotFound"`
	StartButton                         string `json:"StartButton"`
	ListGames                           string `json:"ListGames"`
	NoActiveGames                       string `json:"NoActiveGames"`
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
