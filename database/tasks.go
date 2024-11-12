package database

import (
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

func initTaskTable() {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		task_id SERIAL PRIMARY KEY,
		campaign_id INT REFERENCES campaigns(campaign_id) ON DELETE CASCADE,
		type VARCHAR(50) NOT NULL CHECK (type <> ''),
		description TEXT,
		onboarding_reward FLOAT DEFAULT 0 CHECK (onboarding_reward >= 0),
    	onboarding_threshold FLOAT DEFAULT 0 CHECK (onboarding_threshold >= 0),
		points_pool FLOAT DEFAULT 0 CHECK (points_pool >= 0),
		start_time BIGINT NOT NULL CHECK (start_time >= 0),
		end_time BIGINT NOT NULL CHECK (end_time > start_time),
		created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
		updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
		);
	CREATE INDEX IF NOT EXISTS idx_campaign_id ON tasks(campaign_id);
	CREATE INDEX IF NOT EXISTS idx_task_type ON tasks(type);
	CREATE INDEX IF NOT EXISTS idx_task_time ON tasks(start_time, end_time);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create tasks table and indexes: %v", err)
	}
	fmt.Println("Tasks table and indexes checked/created.")
}

func CreateTask(campaignID int, taskType, description string, onboardingReward float64, onboardingThreshold float64, pointsPool float64, startTime, endTime int64) (int, error) {
	var taskID int
	query := `INSERT INTO tasks (campaign_id, type, description, onboarding_reward, onboarding_threshold, points_pool, start_time, end_time, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING task_id`
	err := db.QueryRow(query, campaignID, taskType, description, onboardingReward, onboardingThreshold, pointsPool, startTime, endTime, time.Now().Unix()).Scan(&taskID)
	return taskID, err
}

func CreateOnboardingTask(campaignID int, description string, onboardingReward float64, onboardingThreshold float64, startTime, endTime int64) (int, error) {
	return CreateTask(campaignID, "onboarding", description, onboardingReward, onboardingThreshold, 0, startTime, endTime)
}

func CreateSharePoolTask(campaignID int, description string, pointsPool float64, startTime, endTime int64) (int, error) {
	return CreateTask(campaignID, "share_pool", description, 0, 0, pointsPool, startTime, endTime)
}

func GetTasksByTaskIDs(taskIDs []int) ([]Task, error) {
	query := `SELECT task_id, campaign_id, type, description, onboarding_reward, onboarding_threshold, points_pool, start_time, end_time FROM tasks WHERE task_id = ANY($1)`
	return queryTasks(query, pq.Array(taskIDs))
}

func GetTasksByCampaignID(campaignID int) ([]Task, error) {
	query := `SELECT task_id, campaign_id, type, description, onboarding_reward, onboarding_threshold, points_pool, start_time, end_time FROM tasks WHERE campaign_id = $1`
	return queryTasks(query, campaignID)
}

func GetActiveTasksByCampaignID(campaignID int, timestamp int64) ([]Task, error) {
	query := `SELECT task_id, campaign_id, type, description, onboarding_reward, onboarding_threshold, points_pool, start_time, end_time 
	FROM tasks WHERE campaign_id = $1 AND end_time > $2 AND start_time < $2`
	return queryTasks(query, campaignID, timestamp)
}

func GetOnboardingTaskByCampaignID(campaignID int) (*Task, error) {
	var task Task
	query := `SELECT task_id, campaign_id, type, description, onboarding_reward, onboarding_threshold, points_pool, start_time, end_time 
	FROM tasks WHERE campaign_id = $1 AND type = 'onboarding'`
	err := db.QueryRow(query, campaignID).Scan(&task.TaskID, &task.CampaignID, &task.Type, &task.Description, &task.OnboardingReward, &task.OnboardingThreshold, &task.PointsPool, &task.StartTime, &task.EndTime)
	return &task, err

}

func GetExpiredSharePoolTasks(now int64, lastCheckTime int64) ([]Task, error) {
	query := `SELECT task_id, campaign_id, type, description, onboarding_reward, onboarding_threshold, points_pool, start_time, end_time 
	FROM tasks WHERE end_time < $1 AND end_time > $2 AND type = 'share_pool'`
	return queryTasks(query, now, lastCheckTime)
}

func queryTasks(query string, args ...interface{}) ([]Task, error) {
	var tasks []Task
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err = rows.Scan(&task.TaskID, &task.CampaignID, &task.Type, &task.Description, &task.OnboardingReward, &task.OnboardingThreshold, &task.PointsPool, &task.StartTime, &task.EndTime)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
