package config

import (
	"github.com/joaop/psiencontra/api/schemas"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&schemas.Session{},
		&schemas.Response{},
		&schemas.Result{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
