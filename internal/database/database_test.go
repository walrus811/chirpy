package database

import (
	"os"
	"testing"
)

const temp_db_path = "temp_db.json"

func TestCreateDB(t *testing.T) {
	db, newDBErr := NewDB(temp_db_path)
	if newDBErr != nil {
		t.Errorf("Error creating DB: %v", newDBErr)
	}
	if db == nil {
		t.Errorf("DB is nil")
	}

	// Cleanup

	removeErr := os.Remove(temp_db_path)
	if removeErr != nil {
		t.Errorf("Error cleaning up: %v", removeErr)
	}
}

func TestGetChirps(t *testing.T) {
	db, newDBErr := NewDB(temp_db_path)
	if newDBErr != nil {
		t.Errorf("Error creating DB: %v", newDBErr)
	}
	if db == nil {
		t.Errorf("DB is nil")
	}

	_, createChirp := db.CreateChirp("test chirp1")

	if createChirp != nil {
		t.Errorf("Error creating chirp: %v", createChirp)
	}

	chirps, getErr := db.GetChirps()
	if getErr != nil {
		t.Errorf("Error getting chirps: %v", getErr)
	}

	if len(chirps) != 1 {
		t.Errorf("Expected 1 chirps, got %v", len(chirps))
	}

	// Cleanup

	removeErr := os.Remove(temp_db_path)
	if removeErr != nil {
		t.Errorf("Error cleaning up: %v", removeErr)
	}
}
