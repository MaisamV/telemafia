package entity

// Role represents a single assignable role with its name and side affiliation.
// This is used *after* extracting roles from the Scenario structure for assignment.
type Role struct {
	Name string
	Side string // e.g., "Mafia", "Civilian", "Neutral"
}

// Side represents a group of roles within a scenario.
type Side struct {
	Name        string   `json:"name"`
	DefaultRole string   `json:"default_role,omitempty"` // Optional: A representative role for the side
	Roles       []string `json:"roles"`                  // List of role names belonging to this side
}

// Scenario represents a game scenario containing sides and their roles.
type Scenario struct {
	ID    string `json:"-"` // Internal ID, not usually part of the input JSON
	Name  string `json:"name"`
	Sides []Side `json:"sides"`
}
