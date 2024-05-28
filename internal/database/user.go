package database

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (db *DB) toHash(text string) (string, error) {
	hashed, bcryptErr := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if bcryptErr != nil {
		return "", bcryptErr
	}

	return string(hashed), nil
}

func (db *DB) LoginUser(email, password string) (User, error) {
	for _, user := range db.dbStructure.Users {
		if user.Email == email {
			compareErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

			if compareErr != nil {
				return User{}, fmt.Errorf("password is incorrect")
			}

			return user, nil
		}
	}

	return User{}, fmt.Errorf("user not found")
}

func (db *DB) CreateUser(email, password string) (User, error) {

	if len(email) == 0 || len(password) == 0 {
		return User{}, fmt.Errorf("email and password cannot be empty")
	}

	if db.existUser(email) {
		return User{}, fmt.Errorf("user already exists")
	}

	hashed, bcryptErr := db.toHash(password)

	if bcryptErr != nil {
		return User{}, bcryptErr
	}

	newUser := User{
		Id:       len(db.dbStructure.Users) + 1,
		Email:    email,
		Password: string(hashed),
	}

	db.dbStructure.Users[newUser.Id] = newUser

	err := db.writeDB(db.dbStructure)

	if err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (db *DB) UpdateUser(id int, newEmail, newPassword string) (User, error) {

	user, getErr := db.GetUser(id)

	if getErr != nil {
		return User{}, getErr
	}

	if len(newEmail) != 0 {
		user.Email = newEmail
	}

	if len(newPassword) != 0 {
		hashed, bcryptErr := db.toHash(newPassword)
		if bcryptErr != nil {
			return User{}, bcryptErr
		}
		user.Password = string(hashed)
	}

	db.dbStructure.Users[user.Id] = user

	dbErr := db.writeDB(db.dbStructure)
	if dbErr != nil {
		return User{}, dbErr
	}

	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	users := make([]User, 0)

	for _, user := range db.dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

func (db *DB) GetUser(id int) (User, error) {
	user, ok := db.dbStructure.Users[id]

	if !ok {
		return User{}, fmt.Errorf("there's no user of %d", id)
	}

	return user, nil
}

func (db *DB) existUser(email string) bool {
	for _, user := range db.dbStructure.Users {
		if user.Email == email {
			return true
		}
	}

	return false
}
