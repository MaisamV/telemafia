package tgutil

import (
	"sync"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	sharedEntity "telemafia/internal/shared/entity"
)

// InteractiveSelectionState holds the temporary state for the "Choose Card" flow
type InteractiveSelectionState struct {
	Mutex         sync.Mutex
	ShuffledRoles []scenarioEntity.Role
	Selections    map[sharedEntity.UserID]PlayerSelection // playerID -> chosenIndex (1-based)
	TakenIndices  map[int]bool                            // chosenIndex (1-based) -> true
}

type PlayerSelection struct {
	ChosenIndex int
	Player      sharedEntity.User
}
