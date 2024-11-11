package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func initUserTaskTable() {
	query := `
	CREATE TABLE IF NOT EXISTS user_tasks (
		user_task_id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
		task_id INT REFERENCES tasks(task_id) ON DELETE CASCADE,
		completed BOOLEAN DEFAULT FALSE,
		amount FLOAT DEFAULT 0,
		points FLOAT DEFAULT 0,
		created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
		updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
	);
	CREATE INDEX IF NOT EXISTS idx_user_id_task_id ON user_tasks(user_id, task_id);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create user_tasks table and indexes: %v", err)
	}
	fmt.Println("UserTasks table and indexes checked/created.")
}

func CreateUserTask(userID, taskID int, completed bool, amount, points float64) (int, error) {
	var userTaskID int
	query := `INSERT INTO user_tasks (user_id, task_id, completed, amount, points, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_task_id`
	err := db.QueryRow(query, userID, taskID, completed, amount, points, time.Now().Unix()).Scan(&userTaskID)
	return userTaskID, err
}

func GetUserTasksByUserID(userID int) ([]UserTask, error) {
	query := `SELECT user_task_id, user_id, task_id, completed, amount, points, created_at, updated_at FROM user_tasks WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTasks []UserTask
	for rows.Next() {
		var userTask UserTask
		if err := rows.Scan(&userTask.UserTaskID, &userTask.UserID, &userTask.TaskID, &userTask.Completed, &userTask.Amount, &userTask.Points, &userTask.CreatedAt, &userTask.UpdatedAt); err != nil {
			return nil, err
		}
		userTasks = append(userTasks, userTask)
	}
	return userTasks, nil
}

func GetUserTasksByUserIDTaskIDs(userID int, taskIDs []int) ([]UserTask, error) {
	query := `SELECT user_task_id, user_id, task_id, completed, amount, points, created_at, updated_at FROM user_tasks WHERE user_id = $1 AND task_id = ANY($2)`
	rows, err := db.Query(query, userID, taskIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTasks []UserTask
	for rows.Next() {
		var userTask UserTask
		if err := rows.Scan(&userTask.UserTaskID, &userTask.UserID, &userTask.TaskID, &userTask.Completed, &userTask.Amount, &userTask.Points, &userTask.CreatedAt, &userTask.UpdatedAt); err != nil {
			return nil, err
		}
		userTasks = append(userTasks, userTask)
	}
	return userTasks, nil
}

func GetOrCreateOnboardingUserTask(userID, taskID int) (UserTask, error) {
	userTask, err := GetUserTaskByUserIDTaskID(userID, taskID)
	if err == sql.ErrNoRows {
		userTaskID, err := CreateUserTask(userID, taskID, false, 0, 0)
		if err != nil {
			return UserTask{}, fmt.Errorf("failed to create onboarding user task: %w", err)
		}
		return UserTask{UserTaskID: userTaskID, UserID: userID, TaskID: taskID, Completed: false, Amount: 0, Points: 0}, nil
	} else if err != nil {
		return UserTask{}, fmt.Errorf("failed to get user task: %w", err)
	}
	return userTask, nil
}

func GetUserTaskByUserIDTaskID(userID, taskID int) (UserTask, error) {
	query := `SELECT user_task_id, user_id, task_id, completed, amount, points, created_at, updated_at FROM user_tasks WHERE user_id = $1 AND task_id = $2`
	var userTask UserTask
	err := db.QueryRow(query, userID, taskID).Scan(&userTask.UserTaskID, &userTask.UserID, &userTask.TaskID, &userTask.Completed, &userTask.Amount, &userTask.Points, &userTask.CreatedAt, &userTask.UpdatedAt)
	if err != nil {
		return UserTask{}, err
	}
	return userTask, nil
}

func UpdateUserTask(userTaskID int, completed bool, amount, points float64) error {
	query := `UPDATE user_tasks SET completed = $2, amount = $3, points = $4, updated_at = $5 WHERE user_task_id = $1`
	_, err := db.Exec(query, userTaskID, completed, amount, points, time.Now().Unix())
	return err
}

func IncreaseUserTaskAmount(taskID int, userID int, amount float64) error {
	query := `UPDATE user_tasks SET amount = amount + $3, updated_at = $4 WHERE task_id = $1 AND user_id = $2`
	_, err := db.Exec(query, taskID, userID, amount, time.Now().Unix())
	return err
}

func UpdateUserTaskByUserIDTaskID(userID, taskID int, completed bool, amount, points float64) error {
	query := `UPDATE user_tasks SET completed = $3, amount = $4, points = $5, updated_at = $6 WHERE user_id = $1 AND task_id = $2`
	_, err := db.Exec(query, userID, taskID, completed, amount, points, time.Now().Unix())
	return err
}
