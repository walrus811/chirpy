package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func (db *DB) GetRefreshToken(userId int) (string, error) {
	token, ok := db.dbStructure.RefreshTokens[userId]

	if !ok {
		return "", errors.New("token not found")
	}

	return token, nil
}

func (db *DB) GetUserIdByToken(token string) (int, error) {
	for userId, existingToken := range db.dbStructure.RefreshTokens {
		if existingToken == token {
			return userId, nil
		}
	}

	return 0, errors.New("token not found")
}

func (db *DB) CreateRefreshToken(userId int) (string, error) {
	dummyString := make([]byte, 32)
	_, readErr := rand.Read(dummyString)

	if readErr != nil {
		return "", readErr
	}

	encoded := hex.EncodeToString(dummyString)

	db.dbStructure.RefreshTokens[userId] = encoded

	err := db.writeDB(db.dbStructure)

	if err != nil {
		return "", err
	}

	return encoded, nil
}

func (db *DB) DeleteRefreshToken(token string) error {
	for userId, existingToken := range db.dbStructure.RefreshTokens {
		if existingToken == token {
			delete(db.dbStructure.RefreshTokens, userId)
			err := db.writeDB(db.dbStructure)

			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}
