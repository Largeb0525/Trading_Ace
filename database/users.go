package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func initUserTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		address VARCHAR(100) UNIQUE NOT NULL,
		created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_user_address ON users(address);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create users table and index: %v", err)
	}
	fmt.Println("Users table and index checked/created.")
}

func CreateUser(address string) (int, error) {
	var userID int
	query := `INSERT INTO users (address, created_at) 
	VALUES ($1, $2) RETURNING user_id`
	err := db.QueryRow(query, address, time.Now().Unix()).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	return userID, nil
}

func GetUserByID(userID int) (*User, error) {
	user := &User{}
	query := `SELECT user_id, address, created_at FROM users WHERE user_id = $1`
	err := db.QueryRow(query, userID).Scan(&user.UserID, &user.Address, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user, nil
}

func GetUserByAddress(address string) (*User, error) {
	user := &User{}
	query := `SELECT user_id, address, created_at FROM users WHERE address = $1`
	err := db.QueryRow(query, address).Scan(&user.UserID, &user.Address, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetOrCreateUserID(address string) (int, error) {
	user, err := GetUserByAddress(address)
	if err == nil {
		return user.UserID, nil
	}
	if err != sql.ErrNoRows {
		return 0, fmt.Errorf("failed to query user by address: %w", err)
	}
	return CreateUser(address)
}
