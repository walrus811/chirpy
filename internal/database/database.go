package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DBStructure struct {
	Chirps        map[int]Chirp  `json:"chirps"`
	Users         map[int]User   `json:"users"`
	RefreshTokens map[int]string `json:"refreshTokens"`
}

type DB struct {
	path        string
	mux         *sync.RWMutex
	dbStructure DBStructure
}

func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if os.IsNotExist(err) {
		file, err := os.Create(db.path)
		if err != nil {
			return err
		}

		_, err = file.WriteString(`{"chirps":{}, "users":{}, "refreshTokens":{}}`)
		if err != nil {
			return err
		}

		defer file.Close()
	}

	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	file, err := os.Open(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	defer file.Close()

	structure := DBStructure{}
	err = json.NewDecoder(file).Decode(&structure)
	if err != nil {
		return DBStructure{}, err
	}

	return structure, nil
}

func (db *DB) writeDB(structure DBStructure) error {
	file, err := os.OpenFile(db.path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(structure)
	if err != nil {
		return err
	}

	return nil
}

func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	ensureError := db.ensureDB()
	if ensureError != nil {
		return nil, ensureError
	}

	dbStructure, loadErr := db.loadDB()
	if loadErr != nil {
		return nil, loadErr
	}

	db.dbStructure = dbStructure

	return &db, nil
}
