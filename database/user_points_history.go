package database

import (
	"fmt"
	"log"
	"time"
)

func initUserPointsHistoryTable() {
	query := `
	CREATE TABLE IF NOT EXISTS user_points_history (
		history_id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
		task_id INT REFERENCES tasks(task_id) ON DELETE CASCADE,
		campaign_id INT REFERENCES campaigns(campaign_id) ON DELETE CASCADE,
		points FLOAT NOT NULL,
		created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
		);
	CREATE INDEX IF NOT EXISTS idx_user_id ON user_points_history(user_id);
	CREATE INDEX IF NOT EXISTS idx_task_id ON user_points_history(task_id);
	CREATE INDEX IF NOT EXISTS idx_campaign_id ON user_points_history(campaign_id);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create user_points_history table and indexes: %v", err)
	}
	fmt.Println("UserPointsHistory table and indexes checked/created.")
}

func CreateUserPointsHistory(userID, taskID, campaignID int, points float64) error {
	query := `INSERT INTO user_points_history (user_id, task_id, campaign_id, points, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, userID, taskID, campaignID, points, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to create user_points_history: %w", err)
	}
	return nil
}

func GetUserPointsHistoryByUserID(userID int) ([]UserPointsHistory, error) {
	query := `SELECT history_id, user_id, task_id, campaign_id, points, created_at FROM user_points_history WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user points history for user_id %d: %w", userID, err)
	}
	defer rows.Close()

	var histories []UserPointsHistory
	for rows.Next() {
		var history UserPointsHistory
		if err := rows.Scan(&history.HistoryID, &history.UserID, &history.TaskID, &history.CampaignID, &history.Points, &history.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user points history: %w", err)
		}
		histories = append(histories, history)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return histories, nil
}
