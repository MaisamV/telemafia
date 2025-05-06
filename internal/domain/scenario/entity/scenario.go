package entity

// Role represents a single assignable role with its name and side affiliation.
// This is used *after* extracting roles from the Scenario structure for assignment.
type Role struct {
	Name    string `json:"name"`
	ImageID string `json:"image_id,omitempty"`
	Side    string `json:"side,omitempty"` // e.g., "Mafia", "Civilian", "Neutral"
}

// Side represents a group of roles within a scenario.
type Side struct {
	Name        string `json:"name"`
	DefaultRole Role   `json:"default_role,omitempty"`
	Roles       []Role `json:"roles,omitempty"` // List of role names belonging to this side
}

// Scenario represents a game scenario containing sides and their roles.
type Scenario struct {
	ID    string `json:"-"` // Internal ID, not usually part of the input JSON
	Name  string `json:"name"`
	Sides []Side `json:"sides"`
}

func (s *Scenario) FlatRoles() []Role {
	flatRoles := make([]Role, 0)
	for _, side := range s.Sides {
		for _, role := range side.Roles {
			role.Side = side.Name
			flatRoles = append(flatRoles, role)
		}
	}
	return flatRoles
}
