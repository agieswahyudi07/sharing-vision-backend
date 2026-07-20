package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "article"
	}

	// 1. Connect to MySQL server (without DB name) to create database if not exists
	tlsParamRaw := ""
	if dbHost != "127.0.0.1" && dbHost != "localhost" {
		tlsParamRaw = "&tls=skip-verify"
	}

	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local%s",
		dbUser, dbPassword, dbHost, dbPort, tlsParamRaw)

	rawDB, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		log.Fatalf("Failed to open connection to MySQL server: %v", err)
	}
	defer rawDB.Close()

	// Try pinging the server
	if err := rawDB.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL server: %v", err)
	}

	_, err = rawDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName))
	if err != nil {
		log.Fatalf("Failed to create database %s: %v", dbName, err)
	}

	// 2. Connect to the specific database using GORM
	tlsParam := ""
	if dbHost != "127.0.0.1" && dbHost != "localhost" {
		tlsParam = "&tls=skip-verify"
	}

	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, tlsParam)

	// Configure logging
	newLogger := glogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		glogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  glogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	DB, err = gorm.Open(gmysql.Open(dsnWithDB), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database using GORM: %v", err)
	}

	log.Printf("Successfully connected to database: %s", dbName)
}
