package entity

import (
	"sort"
	"telemafia/internal/shared/common"
)

// Role represents a single assignable role with its name and side affiliation.
// This is used *after* extracting roles from the Scenario structure for assignment.
type Role struct {
	Name    string `json:"name"`
	AddedAt int    `json:"added_at,omitempty"`
	ImageID string `json:"image_id,omitempty"`
	Side    string `json:"side,omitempty"` // e.g., "Mafia", "Civilian", "Neutral"
}

// Side represents a group of roles within a scenario.
type Side struct {
	Name           string   `json:"name"`
	PopulationRate *float32 `json:"population_rate,omitempty"`
	DefaultRole    *Role    `json:"default_role,omitempty"`
	Roles          []Role   `json:"roles,omitempty"` // List of role names belonging to this side
}

// Scenario represents a game scenario containing sides and their roles.
type Scenario struct {
	ID    string `json:"-"` // Internal ID, not usually part of the input JSON
	Name  string `json:"name"`
	Sides []Side `json:"sides"`
}

func (s *Scenario) FlatRoles(playerNum int) []Role {
	flatRoles := make([]Role, 0)
	for _, side := range s.Sides {
		for _, role := range side.Roles {
			if playerNum >= role.AddedAt {
				role.Side = side.Name
				flatRoles = append(flatRoles, role)
			}
		}
	}
	return flatRoles
}

func (s *Scenario) GetRoles(playerNum int) []Role {
	flatRoles := s.FlatRoles(playerNum)
	for _, side := range s.Sides {
		currentRoleNum := len(flatRoles)
		if currentRoleNum < playerNum && side.DefaultRole != nil {
			side.DefaultRole.Side = side.Name
			if side.PopulationRate != nil && *side.PopulationRate >= 0.01 && *side.PopulationRate <= 0.99 {
				sidePopulationNum := int(float32(playerNum) * (*side.PopulationRate))
				defaultRoleNum := sidePopulationNum - len(side.Roles)
				for i := 0; i < defaultRoleNum; i++ {
					flatRoles = append(flatRoles, *side.DefaultRole)
				}
			} else {
				defaultRoleNum := playerNum - currentRoleNum
				for i := 0; i < defaultRoleNum; i++ {
					flatRoles = append(flatRoles, *side.DefaultRole)
				}
			}
		}
	}
	sort.Slice(flatRoles, func(i, j int) bool {
		return common.Hash(flatRoles[i].Name) < common.Hash(flatRoles[j].Name)
	})
	return flatRoles
}

func (s *Scenario) GetShuffledRoles(playerNum int) []Role {
	flatRoles := s.GetRoles(playerNum)

	common.Shuffle(len(flatRoles), func(i, j int) {
		flatRoles[i], flatRoles[j] = flatRoles[j], flatRoles[i]
	})

	return flatRoles
}
