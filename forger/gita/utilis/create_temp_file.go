package utilis

import (
	"encoding/json"
	"os"
)

func CreateTempJSONFile(data interface{}) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "active_users_*.json")
	if err != nil {
		return nil, err
	}

	// Encode the data to JSON and write to the temporary file
	if err := json.NewEncoder(tempFile).Encode(data); err != nil {
		return nil, err
	}

	// Return the temporary file path and the file object
	return tempFile, nil
}
