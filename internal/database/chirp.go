package database

import "fmt"

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	_, getUserErr := db.GetUser(authorId)

	if getUserErr != nil {
		return Chirp{}, fmt.Errorf("user not found")
	}

	newChirp := Chirp{
		Id:       len(db.dbStructure.Chirps) + 1,
		Body:     body,
		AuthorId: authorId,
	}

	db.dbStructure.Chirps[newChirp.Id] = newChirp

	err := db.writeDB(db.dbStructure)

	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	_, ok := db.dbStructure.Chirps[id]

	if !ok {
		return fmt.Errorf("chirp not found")
	}

	delete(db.dbStructure.Chirps, id)

	err := db.writeDB(db.dbStructure)

	if err != nil {
		return err
	}

	return nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := make([]Chirp, 0)

	for _, chirp := range db.dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpsByAuthorId(authorId int) ([]Chirp, error) {
	chirps := make([]Chirp, 0)

	for _, chirp := range db.dbStructure.Chirps {
		if chirp.AuthorId == authorId {
			chirps = append(chirps, chirp)
		}
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
