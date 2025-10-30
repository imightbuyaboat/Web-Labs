package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"restapi/user"
)

func (ps *PostgresStore) InsertUser(data *user.UserData) (int, error) {
	var userID int
	query := "insert into users (login, hash) values ($1, $2) returning id"

	err := ps.db.QueryRow(query, data.Login, createHash(data.Password)).Scan(&userID)
	if err != nil {
		return -1, fmt.Errorf("failed to insert user %s to DB: %v", data.Login, err)
	}

	return userID, nil
}

func (ps *PostgresStore) CheckUser(data *user.UserData) (int, error) {
	var userID int
	var hashFromDb string
	query := "select id, hash from users where login = $1"

	err := ps.db.QueryRow(query, data.Login).Scan(&userID, &hashFromDb)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, ErrUserNotFound
		}
		return -1, fmt.Errorf("failed to select user %s from DB: %v", data.Login, err)
	}

	if hashFromDb != createHash(data.Password) {
		return -1, ErrIncorrectPassword
	}

	return userID, nil
}

func createHash(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}
