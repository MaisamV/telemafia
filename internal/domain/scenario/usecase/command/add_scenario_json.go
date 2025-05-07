package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioPort "telemafia/internal/domain/scenario/port"
	sharedEntity "telemafia/internal/shared/entity"
)

// AddScenarioJSONCommand represents the command to add a scenario from JSON data.
type AddScenarioJSONCommand struct {
	Requester sharedEntity.User
	JSONData  string
}

// AddScenarioJSONHandler handles the AddScenarioJSONCommand.
type AddScenarioJSONHandler struct {
	scenarioRepo scenarioPort.ScenarioWriter
}

// NewAddScenarioJSONHandler creates a new AddScenarioJSONHandler.
func NewAddScenarioJSONHandler(repo scenarioPort.ScenarioWriter) *AddScenarioJSONHandler {
	return &AddScenarioJSONHandler{scenarioRepo: repo}
}

// Handle executes the command to add a scenario from JSON.
func (h *AddScenarioJSONHandler) Handle(ctx context.Context, cmd AddScenarioJSONCommand) (*scenarioEntity.Scenario, error) {
	// 1. Authorization Check
	if !cmd.Requester.Admin {
		return nil, fmt.Errorf("permission denied: user is not an admin")
	}

	// 2. Unmarshal JSON directly into the domain entity
	var scenario scenarioEntity.Scenario
	if err := json.Unmarshal([]byte(cmd.JSONData), &scenario); err != nil {
		return nil, fmt.Errorf("invalid JSON format: %w", err)
	}

	// 3. Validate Input
	if scenario.Name == "" {
		return nil, fmt.Errorf("scenario name cannot be empty")
	}
	if len(scenario.Sides) == 0 {
		return nil, fmt.Errorf("scenario must have at least one side defined")
	}

	for sideIdx, side := range scenario.Sides {
		if side.Name == "" {
			return nil, fmt.Errorf("side name cannot be empty (side index %d)", sideIdx)
		}
		if len(side.Roles) == 0 && side.DefaultRole == nil {
			return nil, fmt.Errorf("side '%s' must have at least one role", side.Name)
		}

		for roleIdx, role := range side.Roles {
			if role.Name == "" {
				return nil, fmt.Errorf("role name cannot be empty (side '%s', role index %d)", side.Name, roleIdx)
			}
		}
	}

	// 4. Assign Internal ID (Input JSON doesn't contain it)
	scenario.ID = fmt.Sprintf("scen_%d", time.Now().UnixNano())

	// 5. Persist using Repository (No transformation needed now)
	if err := h.scenarioRepo.CreateScenario(&scenario); err != nil {
		return nil, fmt.Errorf("failed to create scenario in repository: %w", err)
	}

	// 6. Return created entity
	return &scenario, nil
}
