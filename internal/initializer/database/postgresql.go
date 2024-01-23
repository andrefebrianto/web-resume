package database

import (
	"fmt"
	"log"
	"net/url"

	"devoratio.dev/web-resume/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgreSQL(dbConfig config.PostgreSQL) *gorm.DB {
	username := dbConfig.Credential.Username
	password := dbConfig.Credential.Password
	instance := dbConfig.Primary

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		username,
		url.QueryEscape(password),
		instance.Host,
		instance.Port,
		instance.DBName,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to postgresql instances: %s", err)
	}

	return db
}
