package database

import (
	"fmt"
	"github.com/esirangelomub/go-chat-application/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitializeDB(config *configs.Conf) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	dsn := GetDSN(config)

	switch config.DBDriver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite3":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		log.Fatalf("Unsupported DB driver: %s", config.DBDriver)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetDSN(config *configs.Conf) string {
	switch config.DBDriver {
	case "postgres":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			config.DBHost,
			config.DBUser,
			config.DBPassword,
			config.DBName,
			config.DBPort,
		)
	case "sqlite3":
		return fmt.Sprintf(
			"%s",
			config.DBName,
		)
	default:
		log.Fatalf("Unsupported DB driver: %s", config.DBDriver)
		return ""
	}
}
