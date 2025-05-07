package tests // Assuming tests are in their own package, adjust if needed

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	// Adjust these imports based on your actual project structure
	// These are based on the imports in m2.go and common Go project layouts
	memrepo "telemafia/internal/adapters/repository/memory"
	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	"telemafia/internal/shared/common"
	sharedEntity "telemafia/internal/shared/entity"
)

func TestRoleShuffleDistribution(t *testing.T) {
	common.InitSeed() // Important for reproducibility if not using t.Setenv for seed

	// Test parameters (can be adjusted)
	playerNum := 13
	repeatIterations := 100000        // Number of times to shuffle and check distribution
	expectedMafiaNum := playerNum / 3 // Based on the logic in m2.go

	// 1. Load Scenario
	// IMPORTANT: Replace with the correct relative path from this test file
	// Path updated assuming the file is now in tests/unit/
	scenarioFilePath := "../../resources/godfather.json"

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

	scenarioRepo := memrepo.NewInMemoryScenarioRepository()
	scenarioHandler := scenarioCommand.NewAddScenarioJSONHandler(scenarioRepo)
	scenario, err := scenarioHandler.Handle(context.Background(), scenarioCommand.AddScenarioJSONCommand{
		Requester: sharedEntity.User{Admin: true},
		JSONData:  jsonData,
	})
	if err != nil {
		t.Fatalf("Failed to load scenario: %v", err)
	}
	if scenario == nil {
		t.Fatal("Loaded scenario is nil")
	}

	// 2. Create mock users
	users := make([]sharedEntity.User, playerNum)
	for i := 0; i < playerNum; i++ {
		users[i] = sharedEntity.User{ID: sharedEntity.UserID(i)}
	}

	// 3. Record role distribution
	mafiaSideCounts := make([]int, playerNum)

	for i := 0; i < repeatIterations; i++ {
		shuffledRoles := scenario.GetShuffledRoles(len(users))

		if len(shuffledRoles) != playerNum {
			t.Fatalf("Iteration %d: Expected %d roles, got %d", i, playerNum, len(shuffledRoles))
		}

		for slotIndex, role := range shuffledRoles {
			// Using the side name "مافیا" as seen in m2.go
			if role.Side == "مافیا" {
				mafiaSideCounts[slotIndex]++
			}
		}
	}

	// 4. Analyze and report distribution for Mafia side
	t.Logf("Mafia side distribution across %d player slots after %d iterations:", playerNum, repeatIterations)
	expectedMafiaChance := float64(expectedMafiaNum) / float64(playerNum)
	t.Logf("Expected Mafia chance per slot: %.4f (based on %d mafia / %d players)", expectedMafiaChance, expectedMafiaNum, playerNum)

	tolerance := 0.01
	distributionOK := true

	for slot := 0; slot < playerNum; slot++ {
		actualMafiaChanceInSlot := float64(mafiaSideCounts[slot]) / float64(repeatIterations)
		t.Logf("Slot %d: Mafia assignments = %d (%.4f chance)", slot, mafiaSideCounts[slot], actualMafiaChanceInSlot)

		if actualMafiaChanceInSlot > expectedMafiaChance+tolerance || actualMafiaChanceInSlot < expectedMafiaChance-tolerance {
			distributionOK = false
			t.Errorf("Slot %d: Mafia chance %.4f is outside the expected range (%.4f +/- %.4f)",
				slot, actualMafiaChanceInSlot, expectedMafiaChance, tolerance)
		}
	}

	if distributionOK {
		t.Log("Mafia distribution appears within tolerance.")
	} else {
		t.Log("Mafia distribution variance detected. See ERRORS above.")
	}
}
