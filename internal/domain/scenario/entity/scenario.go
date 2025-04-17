package entity

// Role represents a role within a scenario
type Role struct {
	Name string
}

// Scenario represents a scenario containing roles
type Scenario struct {
	ID    string // Consider using a more specific type like ScenarioID if needed elsewhere
	Name  string
	Roles []Role
}
