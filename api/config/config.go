package config

import (
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	Log *Logger
)

func Init() error {
	// Load .env from project root (../. relative to api/)
	_ = godotenv.Load("../.env")

	Log = NewLogger()

	dsn := GetEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/psiencontra?sslmode=disable")
	db, err := NewDB(dsn)
	if err != nil {
		return err
	}
	DB = db

	Log.Info.Println("Database connected and migrated")
	return nil
}
