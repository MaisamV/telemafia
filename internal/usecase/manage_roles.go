package usecase

import (
	"context"
	"telemafia/internal/entity"
)

// AddRoleCommand represents the command to add a role to a scenario
type AddRoleCommand struct {
	ScenarioID string
	Role       entity.Role
}

// RemoveRoleCommand represents the command to remove a role from a scenario
type RemoveRoleCommand struct {
	ScenarioID string
	RoleName   string
}

// ManageRolesHandler handles adding and removing roles from scenarios
type ManageRolesHandler struct {
	scenarioRepo ScenarioRepository
}

// NewManageRolesHandler creates a new ManageRolesHandler
func NewManageRolesHandler(repo ScenarioRepository) *ManageRolesHandler {
	return &ManageRolesHandler{scenarioRepo: repo}
}

// HandleAddRole adds a role to a scenario
func (h *ManageRolesHandler) HandleAddRole(ctx context.Context, cmd AddRoleCommand) error {
	return h.scenarioRepo.AddRoleToScenario(cmd.ScenarioID, cmd.Role)
}

// HandleRemoveRole removes a role from a scenario
func (h *ManageRolesHandler) HandleRemoveRole(ctx context.Context, cmd RemoveRoleCommand) error {
	return h.scenarioRepo.RemoveRoleFromScenario(cmd.ScenarioID, cmd.RoleName)
}
