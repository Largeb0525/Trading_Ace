package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func initCampaignTable() {
	query := `
	CREATE TABLE IF NOT EXISTS campaigns (
		campaign_id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		pool_address VARCHAR(100) NOT NULL CHECK (pool_address <> ''),
		start_time BIGINT NOT NULL,
		end_time BIGINT NOT NULL,
		created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW()),
		updated_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
		);
	CREATE INDEX IF NOT EXISTS idx_pool_address ON campaigns(pool_address);
	CREATE INDEX IF NOT EXISTS idx_campaign_time ON campaigns(start_time, end_time);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create campaigns table and indexes: %v", err)
	}
	fmt.Println("Campaigns table and indexes checked/created.")
}

func GetCampaignByID(id int) (*Campaign, error) {
	var campaign Campaign
	query := `SELECT campaign_id, name, pool_address, start_time, end_time FROM campaigns WHERE campaign_id = $1`
	err := db.QueryRow(query, id).Scan(&campaign.CampaignID, &campaign.Name, &campaign.PoolAddress, &campaign.StartTime, &campaign.EndTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("campaign with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to query campaign by ID: %w", err)
	}
	return &campaign, nil
}

func GetCampaignsByAddress(address string) ([]Campaign, error) {
	var campaigns []Campaign
	query := `SELECT campaign_id, name, pool_address, start_time, end_time FROM campaigns WHERE pool_address = $1`
	rows, err := db.Query(query, address)
	if err != nil {
		return nil, fmt.Errorf("failed to query campaigns by address: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var campaign Campaign
		if err := rows.Scan(&campaign.CampaignID, &campaign.Name, &campaign.PoolAddress, &campaign.StartTime, &campaign.EndTime); err != nil {
			return nil, fmt.Errorf("failed to scan campaign: %w", err)
		}
		campaigns = append(campaigns, campaign)
	}
	return campaigns, nil
}

func CreateCampaign(name, poolAddress string, startTime, endTime int64) (int, error) {
	var campaignID int
	query := `INSERT INTO campaigns (name, pool_address, start_time, end_time, created_at) 
	VALUES ($1, $2, $3, $4, $5) RETURNING campaign_id`
	err := db.QueryRow(query, name, poolAddress, startTime, endTime, time.Now().Unix()).Scan(&campaignID)
	if err != nil {
		return 0, fmt.Errorf("failed to create campaign: %w", err)
	}
	return campaignID, nil
}

func GetActiveCampaignAddresses() ([]string, error) {
	now := time.Now().Unix()
	query := `SELECT DISTINCT pool_address FROM campaigns WHERE start_time <= $1 AND end_time >= $1`
	rows, err := db.Query(query, now)
	if err != nil {
		return nil, fmt.Errorf("failed to query active campaigns: %w", err)
	}
	defer rows.Close()

	var addresses []string
	for rows.Next() {
		var address string
		if err := rows.Scan(&address); err != nil {
			log.Printf("Failed to scan address: %v", err)
			continue
		}
		addresses = append(addresses, address)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return addresses, nil
}
