package database

import "fmt"

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	newChirp := Chirp{
		Id:   len(db.dbStructure.Chirps) + 1,
		Body: body,
	}

	db.dbStructure.Chirps[newChirp.Id] = newChirp

	err := db.writeDB(db.dbStructure)

	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := make([]Chirp, 0)

	for _, chirp := range db.dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {

	chirp, ok := db.dbStructure.Chirps[id]

	if !ok {
		return Chirp{}, fmt.Errorf("there's no chirp of %d", id)
	}

	return chirp, nil
}
