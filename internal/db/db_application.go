package db

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func IsApplicationInDatabase(applicationName string) bool {
	dbFilePath := filepath.Join(os.Getenv("HOME"), ".teu", "db.json")

	var database Database
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		database = Database{Applications: make(map[string]Application)}
	} else {
		dbFile, err := os.ReadFile(dbFilePath)
		if err != nil {
			return false
		}

		if err := json.Unmarshal(dbFile, &database); err != nil {
			return false
		}
	}

	_, exists := database.Applications[applicationName]
	return exists
}

func AddApplicationToDatabase(applicationName, hash string) error {
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

	database.Applications[applicationName] = Application{Hash: hash}

	dbFile, err := json.Marshal(database)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dbFilePath, dbFile, 0644); err != nil {
		return err
	}

	return nil
}

func RemoveApplicationFromDatabase(applicationName string) error {
	dbFilePath := filepath.Join(os.Getenv("HOME"), ".teu", "db.json")

	var database Database
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		return nil
	}

	dbFile, err := os.ReadFile(dbFilePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(dbFile, &database); err != nil {
		return err
	}

	delete(database.Applications, applicationName)

	dbFile, err = json.Marshal(database)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dbFilePath, dbFile, 0644); err != nil {
		return err
	}
	
	return nil
}

func GetApplicationHashFromDatabase(applicationName string) (string, error) {
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

	application, exists := database.Applications[applicationName]
	if !exists {
		return "", nil
	}

	return application.Hash, nil
}

func GetAllApplicationsFromDatabase() ([]string, error) {
	dbFilePath := filepath.Join(os.Getenv("HOME"), ".teu", "db.json")

	var database Database
	if _, err := os.Stat(dbFilePath); os.IsNotExist(err) {
		return nil, nil
	}

	dbFile, err := os.ReadFile(dbFilePath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dbFile, &database); err != nil {
		return nil, err
	}

	var applications []string
	for app := range database.Applications {
		applications = append(applications, app)
	}

	return applications, nil
}
