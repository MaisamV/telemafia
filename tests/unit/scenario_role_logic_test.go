package tests

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	memrepo "telemafia/internal/adapters/repository/memory"
	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	"telemafia/internal/shared/common"
	sharedEntity "telemafia/internal/shared/entity"
)

func TestScenarioRoleGenerationAndMafiaCount(t *testing.T) {
	common.InitSeed()

	scenarioDir := "../../resources/scenario/" // Adjusted path from tests/unit/
	scenarioFiles := []string{
		"classic.json",
		"godfather.json",
		"sherlock.json",
		"bazpors.json",
		"royabin.json",
	}

	scenarioRepo := memrepo.NewInMemoryScenarioRepository()
	scenarioHandler := scenarioCommand.NewAddScenarioJSONHandler(scenarioRepo)

	for _, fileName := range scenarioFiles {
		t.Run(fileName, func(t *testing.T) {
			scenarioFilePath := filepath.Join(scenarioDir, fileName)
			file, err := os.Open(scenarioFilePath)
			if err != nil {
				t.Fatalf("Failed to open scenario file '%s': %v", scenarioFilePath, err)
			}
			defer file.Close()

			bytes, err := io.ReadAll(file)
			if err != nil {
				t.Fatalf("Failed to read scenario file: %v", err)
			}
			jsonData := strings.TrimSpace(string(bytes))

			scenario, err := scenarioHandler.Handle(context.Background(), scenarioCommand.AddScenarioJSONCommand{
				Requester: sharedEntity.User{Admin: true},
				JSONData:  jsonData,
			})
			if err != nil {
				t.Fatalf("Failed to load scenario from '%s': %v", fileName, err)
			}
			if scenario == nil {
				t.Fatalf("Loaded scenario from '%s' is nil", fileName)
			}

			t.Logf("Scenario: %s (%s)", scenario.Name, fileName)

			for playerNum := 3; playerNum <= 19; playerNum++ {
				t.Run(fmt.Sprintf("PlayerNum_%d", playerNum), func(t *testing.T) {
					shuffledRoles := scenario.GetShuffledRoles(playerNum)

					if len(shuffledRoles) != playerNum {
						t.Errorf("Expected %d roles, got %d", playerNum, len(shuffledRoles))
						// Log the roles if count mismatch for debugging
						// var roleInfo []string
						// for _, r := range shuffledRoles {
						// 	roleInfo = append(roleInfo, fmt.Sprintf("%s (%s)", r.Name, r.Side))
						// }
						// t.Logf("Shuffled Roles (%d): %s", len(shuffledRoles), strings.Join(roleInfo, ", "))
						return // Stop this sub-test if role count is wrong
					}

					// Log the shuffled roles (can be verbose)
					var roleInfo []string
					for _, r := range shuffledRoles {
						roleInfo = append(roleInfo, fmt.Sprintf("%s (%s)", r.Name, r.Side))
					}
					t.Logf("PlayerNum: %d, Shuffled Roles: [%s]", playerNum, strings.Join(roleInfo, "; "))

					mafiaCount := 0
					for _, role := range shuffledRoles {
						if role.Side == "مافیا" { // Assuming this is the consistent side name for Mafia
							mafiaCount++
						}
					}

					expectedMafiaNum := playerNum / 3 // Integer division
					if mafiaCount != expectedMafiaNum {
						t.Errorf("PlayerNum: %d - Expected %d Mafia roles, got %d", playerNum, expectedMafiaNum, mafiaCount)
					}
				})
			}
		})
	}
}
