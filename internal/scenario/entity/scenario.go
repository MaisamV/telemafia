package entity

// Role represents a role within a scenario
type Role struct {
	Name string
}

// Scenario represents a scenario containing roles
type Scenario struct {
	ID    string
	Name  string
	Roles []Role
}
