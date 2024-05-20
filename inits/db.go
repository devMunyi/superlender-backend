package inits

import (
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var CurrentDB *gorm.DB
var ArchiveDB *gorm.DB

func DBInit() {

	// Current database connection
	currentDSN := os.Getenv("CURRENT_DB_URL")

	currentMaxIdleConns, err := strconv.Atoi(os.Getenv("CURRENT_DB_MAX_IDLE_CONNS"))
	if err != nil {
		currentMaxIdleConns = 5 // Set a default value
	}
	currentMaxOpenConns, err := strconv.Atoi(os.Getenv("CURRENT_DB_MAX_OPEN_CONNS"))
	if err != nil {
		currentMaxOpenConns = 50 // Set a default value
	}
	CurrentDB = initDB(currentDSN, currentMaxIdleConns, currentMaxOpenConns)

	// Archive database connection
	archiveDSN := os.Getenv("ARCHIVE_DB_URL")
	archiveMaxIdleConns, err := strconv.Atoi(os.Getenv("ARCHIVE_DB_MAX_IDLE_CONNS"))
	if err != nil {
		archiveMaxIdleConns = 3 // Set a default value
	}
	archiveMaxOpenConns, err := strconv.Atoi(os.Getenv("ARCHIVE_DB_MAX_OPEN_CONNS"))
	if err != nil {
		archiveMaxOpenConns = 50 // Set a default value
	}
	ArchiveDB = initDB(archiveDSN, archiveMaxIdleConns, archiveMaxOpenConns)
}

func initDB(dsn string, maxIdleConns, maxOpenConns int) *gorm.DB {

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,  // DSN data source name
		DefaultStringSize:         256,  // Optional, set the default string size for MySQL columns
		SkipInitializeWithVersion: true, // Optional, disable auto-migration with database version check
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database handle")
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Enable debug mode
	db = db.Debug()

	return db
}
