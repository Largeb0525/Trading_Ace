package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	viper.Set("database.host", "localhost")
	viper.Set("database.port", 5432)
	viper.Set("database.user", "admin")
	viper.Set("database.password", "2222")
	viper.Set("database.dbname", "test_db")

	testDB = InitPostgreSQL()
	db = testDB
	code := m.Run()

	cleanupDatabase()

	db.Close()

	os.Exit(code)
}

func cleanupDatabase() {
	_, err := testDB.Exec("DROP TABLE IF EXISTS user_swaps, user_points_history, user_tasks, tasks, campaigns, users CASCADE")
	if err != nil {
		log.Printf("Failed to clean up database: %v", err)
	}
}
