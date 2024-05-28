package database

import (
	"fmt"
	"os"
	"testing"
)

func TestCreateDB(t *testing.T) {
	const dbPath = "TestCreateDB.json"
	db, newDBErr := NewDB(dbPath)
	if newDBErr != nil {
		t.Errorf("Error creating DB: %v", newDBErr)
	}
	if db == nil {
		t.Errorf("DB is nil")
	}

	// Cleanup

	removeErr := os.Remove(dbPath)
	if removeErr != nil {
		t.Errorf("Error cleaning up: %v", removeErr)
	}
}

func TestGetChirps(t *testing.T) {
	dbPath := "TestGetChirps.json"
	db, newDBErr := NewDB(dbPath)
	if newDBErr != nil {
		t.Errorf("Error creating DB: %v", newDBErr)
	}
	if db == nil {
		t.Errorf("DB is nil")
	}

	testData := []string{
		"t1",
		"t2",
		"t3",
		"t4",
	}

	for _, data := range testData {
		_, createErr := db.CreateChirp(data)
		if createErr != nil {
			t.Errorf("Error creating chirp: %v", createErr)
		}
	}

	chirps, getErr := db.GetChirps()
	if getErr != nil {
		t.Errorf("Error getting chirps: %v", getErr)
	}

	if len(chirps) != 4 {
		t.Errorf("Expected 1 chirps, got %v", len(chirps))
	}

	// Cleanup

	removeErr := os.Remove(dbPath)
	if removeErr != nil {
		t.Errorf("Error cleaning up: %v", removeErr)
	}
}

func TestGetChirp(t *testing.T) {
	dbPath := "TestGetChirp.json"
	db, newDBErr := NewDB(dbPath)
	if newDBErr != nil {
		t.Errorf("Error creating DB: %v", newDBErr)
	}
	if db == nil {
		t.Errorf("DB is nil")
	}

	testData := []string{
		"t1",
		"t2",
		"t3",
		"t4",
	}

	for _, data := range testData {
		_, createErr := db.CreateChirp(data)
		if createErr != nil {
			t.Errorf("Error creating chirp: %v", createErr)
		}
	}
	c, e := db.GetChirp(1)
	fmt.Println(c, e)
	chirp, getErr := db.GetChirp(len(testData))
	if getErr != nil {
		t.Errorf("Error getting chirps: %v", getErr)
	}

	if chirp.Body != testData[len(testData)-1] {
		t.Errorf("Expected %v, got %v", testData[len(testData)-1], chirp.Body)
	}

	// Cleanup

	removeErr := os.Remove(dbPath)
	if removeErr != nil {
		t.Errorf("Error cleaning up: %v", removeErr)
	}
}

func TestCreateUsers(t *testing.T) {
	dbPath := "TestCreateUsers.json"
	db, newDBErr := NewDB(dbPath)
	if newDBErr != nil {
		t.Errorf("Error creating DB: %v", newDBErr)
	}
	if db == nil {
		t.Errorf("DB is nil")
	}
	testData := []User{
		{Email: "t1@naver.com", Password: "1234"},
		{Email: "t2@naver.com", Password: "1234"},
		{Email: "t3@naver.com", Password: "1234"},
		{Email: "t4@naver.com", Password: "1234"},
		{Email: "t5@naver.com", Password: "1234"},
	}

	for _, data := range testData {
		_, createErr := db.CreateUser(data.Email, data.Password)
		if createErr != nil {
			t.Errorf("Error creating user: %v", createErr)
		}
	}

	users, getErr := db.GetUsers()
	if getErr != nil {
		t.Errorf("Error getting users: %v", getErr)
	}

	if len(users) != len(testData) {
		t.Errorf("Expected %d users, got %v", len(testData), len(users))
	}

	// Cleanup

	removeErr := os.Remove(dbPath)
	if removeErr != nil {
		t.Errorf("Error cleaning up: %v", removeErr)
	}
}

func TestLogin(t *testing.T) {
	dbPath := "TestLogin.json"
	db, newDBErr := NewDB(dbPath)
	if newDBErr != nil {
		t.Errorf("Error creating DB: %v", newDBErr)
	}
	if db == nil {
		t.Errorf("DB is nil")
	}
	testData := []User{
		{Email: "t1@naver.com", Password: "1234"},
	}

	for _, data := range testData {
		_, createErr := db.CreateUser(data.Email, data.Password)
		if createErr != nil {
			t.Errorf("Error creating user: %v", createErr)
		}
	}

	_, loginErr := db.LoginUser(testData[0].Email, testData[0].Password)
	if loginErr != nil {
		t.Errorf("Error logging in: %v", loginErr)
	}

	// Cleanup

	removeErr := os.Remove(dbPath)
	if removeErr != nil {
		t.Errorf("Error cleaning up: %v", removeErr)
	}
}

func TestUpdateUsers(t *testing.T) {
	dbPath := "TestUpdateUsers.json"
	db, newDBErr := NewDB(dbPath)
	if newDBErr != nil {
		t.Errorf("Error creating DB: %v", newDBErr)
	}
	if db == nil {
		t.Errorf("DB is nil")
	}
	testData := User{
		Email: "t1@naver.com", Password: "1234",
	}

	newEmail := "t2@gmail.com"
	newPassword := "12345"

	user, createErr := db.CreateUser(testData.Email, testData.Password)
	if createErr != nil {
		t.Errorf("Error creating user: %v", createErr)
	}

	_, updateErr := db.UpdateUser(user.Id, newEmail, newPassword)

	if updateErr != nil {
		t.Errorf("Error updating user: %v", updateErr)
	}

	updatedUser, getErr := db.GetUser(user.Id)

	if getErr != nil {
		t.Errorf("Error getting user: %v", getErr)
	}

	if updatedUser.Email != newEmail {
		t.Errorf("Expected %v, got %v", newEmail, updatedUser.Email)
	}

	// Cleanup

	removeErr := os.Remove(dbPath)
	if removeErr != nil {
		t.Errorf("Error cleaning up: %v", removeErr)
	}
}
