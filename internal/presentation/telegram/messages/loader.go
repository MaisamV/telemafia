package messages

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadMessages reads the specified JSON file and unmarshals it into a Messages struct.
func LoadMessages(filename string) (*Messages, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read messages file '%s': %w", filename, err)
	}

	var msgs Messages
	if err := json.Unmarshal(data, &msgs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages file '%s': %w", filename, err)
	}

	// Optional: Add validation here to ensure critical messages are not empty
	if msgs.Common.Help == "" {
		fmt.Println("Warning: Common.Welcome message is empty in", filename)
	}
	// ... add other checks as needed ...

	fmt.Println("✅ Loaded messages from", filename)
	return &msgs, nil
}
