package db

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Initiate the database
func DatabaseInit() error {
	var database Database
	dbFilePath := filepath.Join(os.Getenv("HOME"), ".teu", "db.json")

	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(os.Getenv("HOME"), ".teu"), 0755); err != nil {
			return err
		}
		dbFile, err := json.Marshal(database)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dbFilePath, dbFile, 0644); err != nil {
			return err
		}
	}

	return nil
}