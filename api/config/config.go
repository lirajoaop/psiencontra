package config

import "gorm.io/gorm"

var (
	DB  *gorm.DB
	Log *Logger
)

func Init() error {
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
