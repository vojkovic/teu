package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

// Common function to get a key from db.
func SetKeyStringInDatabase(key string) (string, error) {
	dbFilePath := filepath.Join(os.Getenv("HOME"), ".teu", "db.json")

	var database Database
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		return "", nil
	}

	dbFile, err := os.ReadFile(dbFilePath)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(dbFile, &database); err != nil {
		return "", err
	}

	val := reflect.ValueOf(database)
	field := val.FieldByName(key)
	if !field.IsValid() {
		return "", fmt.Errorf("key not found")
	}

	return field.String(), nil
}

//Common function to set a key in db.
func SetValueStringInDatabase(key string, value string) error {
	dbFilePath := filepath.Join(os.Getenv("HOME"), ".teu", "db.json")

	var database Database
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		database = Database{Applications: make(map[string]Application)}
	} else {
		dbFile, err := os.ReadFile(dbFilePath)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(dbFile, &database); err != nil {
			return err
		}
	}

	val := reflect.ValueOf(&database).Elem()
	field := val.FieldByName(key)
	if !field.IsValid() {
		return fmt.Errorf("key not found")
	}

	field.SetString(value)

	dbFile, err := json.Marshal(database)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dbFilePath, dbFile, 0644); err != nil {
		return err
	}

	return nil
}