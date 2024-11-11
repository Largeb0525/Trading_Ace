package database

import (
	"fmt"
	"log"
	"time"
)

func initUserSwapTable() {
	query := `
	CREATE TABLE IF NOT EXISTS user_swaps (
		swap_id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
		transaction_hash VARCHAR(100) NOT NULL,
		pool_address VARCHAR(100) NOT NULL,
		amount_usdc FLOAT DEFAULT 0,
		amount_weth FLOAT DEFAULT 0,
		swap_time BIGINT NOT NULL,
		created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
	);
	CREATE INDEX IF NOT EXISTS idx_user_id ON user_swaps(user_id);
	CREATE INDEX IF NOT EXISTS idx_pool_address ON user_swaps(pool_address);
	CREATE INDEX IF NOT EXISTS idx_transaction_hash ON user_swaps(transaction_hash);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create user_swaps table and indexes: %v", err)
	}
	fmt.Println("UserSwaps table and indexes checked/created.")
}

func InsertSwapEvent(userID int, poolAddress string, usdc float64, swapTime int64, txHash string) error {
	query := `
	INSERT INTO user_swaps (user_id, pool_address, amount_usdc, swap_time, transaction_hash, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, userID, poolAddress, usdc, swapTime, txHash, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to insert swap event: %w", err)
	}
	return nil
}
