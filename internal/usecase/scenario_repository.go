package usecase

import (
	"telemafia/internal/entity"
)

// ScenarioReader defines the interface for reading scenario data
type ScenarioReader interface {
	GetScenarioByID(id string) (*entity.Scenario, error)
	GetAllScenarios() ([]*entity.Scenario, error)
}

// ScenarioWriter defines the interface for writing scenario data
type ScenarioWriter interface {
	CreateScenario(scenario *entity.Scenario) error
	DeleteScenario(id string) error
	AddRoleToScenario(scenarioID string, role entity.Role) error
	RemoveRoleFromScenario(scenarioID string, roleName string) error
}

// ScenarioRepository defines the interface for scenario persistence
type ScenarioRepository interface {
	ScenarioReader
	ScenarioWriter
}
