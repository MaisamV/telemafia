package command

import (
	"context"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioPort "telemafia/internal/domain/scenario/port"
)

// AddRoleCommand represents the command to add a role to a scenario
type AddRoleCommand struct {
	ScenarioID string
	Role       scenarioEntity.Role // Use imported Role type
}

// RemoveRoleCommand represents the command to remove a role from a scenario
type RemoveRoleCommand struct {
	ScenarioID string
	RoleName   string
}

// ManageRolesHandler handles adding and removing roles from scenarios
type ManageRolesHandler struct {
	scenarioRepo scenarioPort.ScenarioRepository // Use imported Repository interface
}

// NewManageRolesHandler creates a new ManageRolesHandler
func NewManageRolesHandler(repo scenarioPort.ScenarioRepository) *ManageRolesHandler {
	return &ManageRolesHandler{scenarioRepo: repo}
}

// HandleAddRole adds a role to a scenario
func (h *ManageRolesHandler) HandleAddRole(ctx context.Context, cmd AddRoleCommand) error {
	return h.scenarioRepo.AddRoleToScenario(cmd.ScenarioID, cmd.Role) // Propagates errors
}

// HandleRemoveRole removes a role from a scenario
func (h *ManageRolesHandler) HandleRemoveRole(ctx context.Context, cmd RemoveRoleCommand) error {
	return h.scenarioRepo.RemoveRoleFromScenario(cmd.ScenarioID, cmd.RoleName) // Propagates errors
}
