package memory

import (
	// "errors"
	"fmt"
	"sync"

	// scenarioEntity "telemafia/internal/scenario/entity"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	// scenarioPort "telemafia/internal/scenario/port"
	scenarioPort "telemafia/internal/domain/scenario/port"
)

// Ensure InMemoryScenarioRepository implements the scenarioPort.ScenarioRepository interface.
var _ scenarioPort.ScenarioRepository = (*InMemoryScenarioRepository)(nil)

// InMemoryScenarioRepository provides an in-memory implementation of the scenario repository
type InMemoryScenarioRepository struct {
	data  map[string]*scenarioEntity.Scenario
	mutex sync.RWMutex
}

// NewInMemoryScenarioRepository creates a new in-memory scenario repository
func NewInMemoryScenarioRepository() scenarioPort.ScenarioRepository {
	return &InMemoryScenarioRepository{
		data: make(map[string]*scenarioEntity.Scenario),
	}
}

// GetScenarioByID retrieves a scenario by its ID
func (r *InMemoryScenarioRepository) GetScenarioByID(id string) (*scenarioEntity.Scenario, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	scenario, exists := r.data[id]
	if !exists {
		return nil, fmt.Errorf("scenario with ID %s not found", id)
	}
	return scenario, nil
}

// GetAllScenarios retrieves all scenarios
func (r *InMemoryScenarioRepository) GetAllScenarios() ([]*scenarioEntity.Scenario, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var scenarios []*scenarioEntity.Scenario
	for _, scenario := range r.data {
		scenarios = append(scenarios, scenario)
	}
	return scenarios, nil
}

// CreateScenario adds a new scenario
func (r *InMemoryScenarioRepository) CreateScenario(scenario *scenarioEntity.Scenario) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.data[scenario.ID]; exists {
		return fmt.Errorf("scenario with ID %s already exists", scenario.ID)
	}
	r.data[scenario.ID] = scenario
	return nil
}

// DeleteScenario removes a scenario by its ID
func (r *InMemoryScenarioRepository) DeleteScenario(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.data[id]; !exists {
		return fmt.Errorf("scenario with ID %s not found for deletion", id)
	}
	delete(r.data, id)
	return nil
}
