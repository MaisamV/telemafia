package memory

import (
	"errors"
	"sync"
	"telemafia/internal/entity"
	"telemafia/internal/usecase"
)

// Ensure InMemoryScenarioRepository implements the usecase.ScenarioRepository interface.
var _ usecase.ScenarioRepository = (*InMemoryScenarioRepository)(nil)

// InMemoryScenarioRepository provides an in-memory implementation of the scenario repository
type InMemoryScenarioRepository struct {
	data  map[string]*entity.Scenario
	mutex sync.RWMutex
}

// NewInMemoryScenarioRepository creates a new in-memory scenario repository
func NewInMemoryScenarioRepository() usecase.ScenarioRepository {
	return &InMemoryScenarioRepository{
		data: make(map[string]*entity.Scenario),
	}
}

// GetScenarioByID retrieves a scenario by its ID
func (r *InMemoryScenarioRepository) GetScenarioByID(id string) (*entity.Scenario, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	scenario, exists := r.data[id]
	if !exists {
		return nil, errors.New("scenario not found")
	}
	return scenario, nil
}

// GetAllScenarios retrieves all scenarios
func (r *InMemoryScenarioRepository) GetAllScenarios() ([]*entity.Scenario, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var scenarios []*entity.Scenario
	for _, scenario := range r.data {
		scenarios = append(scenarios, scenario)
	}
	return scenarios, nil
}

// CreateScenario adds a new scenario
func (r *InMemoryScenarioRepository) CreateScenario(scenario *entity.Scenario) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.data[scenario.ID]; exists {
		return errors.New("scenario already exists")
	}
	r.data[scenario.ID] = scenario
	return nil
}

// DeleteScenario removes a scenario by its ID
func (r *InMemoryScenarioRepository) DeleteScenario(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.data[id]; !exists {
		return errors.New("scenario not found")
	}
	delete(r.data, id)
	return nil
}

// AddRoleToScenario adds a role to a scenario
func (r *InMemoryScenarioRepository) AddRoleToScenario(scenarioID string, role entity.Role) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	scenario, exists := r.data[scenarioID]
	if !exists {
		return errors.New("scenario not found")
	}
	// Allow duplicate role names
	scenario.Roles = append(scenario.Roles, role)
	return nil
}

// RemoveRoleFromScenario removes the first instance of a role with the given name from a scenario
func (r *InMemoryScenarioRepository) RemoveRoleFromScenario(scenarioID string, roleName string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	scenario, exists := r.data[scenarioID]
	if !exists {
		return errors.New("scenario not found")
	}

	for i, role := range scenario.Roles {
		if role.Name == roleName {
			// Remove only the first matching instance
			scenario.Roles = append(scenario.Roles[:i], scenario.Roles[i+1:]...)
			return nil
		}
	}
	return errors.New("role not found")
}
