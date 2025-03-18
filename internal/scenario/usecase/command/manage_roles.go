package command

import (
	"context"
	"telemafia/internal/scenario/entity"
	"telemafia/internal/scenario/repo"
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

// ManageRolesHandler handles role management in scenarios
type ManageRolesHandler struct {
	scenarioRepo repo.Repository
}

// NewManageRolesHandler creates a new ManageRolesHandler
func NewManageRolesHandler(repo repo.Repository) *ManageRolesHandler {
	return &ManageRolesHandler{
		scenarioRepo: repo,
	}
}

// HandleAddRole processes the add role command
func (h *ManageRolesHandler) HandleAddRole(ctx context.Context, cmd AddRoleCommand) error {
	return h.scenarioRepo.AddRoleToScenario(cmd.ScenarioID, cmd.Role)
}

// HandleRemoveRole processes the remove role command
func (h *ManageRolesHandler) HandleRemoveRole(ctx context.Context, cmd RemoveRoleCommand) error {
	return h.scenarioRepo.RemoveRoleFromScenario(cmd.ScenarioID, cmd.RoleName)
}
