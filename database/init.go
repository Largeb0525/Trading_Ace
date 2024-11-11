package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db *sql.DB

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func loadConfig() DBConfig {
	return DBConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
	}
}

func connectToTargetDB(config DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	return sql.Open("postgres", connStr)
}

func connectToDefaultDB(config DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		config.Host, config.Port, config.User, config.Password)
	return sql.Open("postgres", connStr)
}

func createDatabase(db *sql.DB, dbName string) error {
	query := fmt.Sprintf("CREATE DATABASE %s", dbName)
	_, err := db.Exec(query)
	return err
}

func InitPostgreSQL() *sql.DB {
	var err error
	config := loadConfig()
	db, err = connectToTargetDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database %s: %v", config.DBName, err)
	}

	err = db.Ping()
	if err != nil && !strings.Contains(err.Error(), "does not exist") {
		log.Fatalf("Failed to ping database %s: %v", config.DBName, err)
	} else if err != nil && strings.Contains(err.Error(), "does not exist") {
		defaultDB, err := connectToDefaultDB(config)
		if err != nil {
			log.Fatalf("Failed to connect to default database: %v", err)
		}

		err = createDatabase(defaultDB, config.DBName)
		if err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		fmt.Printf("Database %s created successfully.\n", config.DBName)
		defaultDB.Close()

		db, err = connectToTargetDB(config)
		if err != nil {
			log.Fatalf("Failed to connect to database %s: %v", config.DBName, err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatalf("Failed to ping database %s: %v", config.DBName, err)
		}
	}

	fmt.Printf("Connected to database %s successfully.\n", config.DBName)

	initTable()
	return db
}

func initTable() {
	initUserTable()
	initCampaignTable()
	initTaskTable()
	initUserTaskTable()
	initUserPointsHistoryTable()
	initUserSwapTable()
}
