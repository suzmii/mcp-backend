package dao

import (
	"errors"
	"fmt"

	mysqldriver "github.com/go-sql-driver/mysql"
	"github.com/lunny/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	DriverName   string
	AutoCreateDB bool
}

func MustNewDB(config DBConfig) *gorm.DB {
	var DB *gorm.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.DatabaseName)
	var err error
	if DB, err = gorm.Open(mysql.Open(dsn)); err != nil {
		if !config.AutoCreateDB {
			log.Fatalf("failed to connect to DB: %v", err)
		}

		var mysqlErr *mysqldriver.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1049 {
			log.Warn(fmt.Sprintf("database %s NOT exist, creating", config.DatabaseName))
			// DataBase [DBName] Not Found

			/*
				1. Connect to server without select DB
				2. Create DB
				3. Use it
			*/

			// Connect to server without DB
			dsnNoDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local",
				config.Username, config.Password, config.Host, config.Port)
			DB, err = gorm.Open(mysql.Open(dsnNoDB), &gorm.Config{})
			if err != nil {
				log.Fatalf("Failed to Open DataBase	 while create DB: %v", err)
			}

			// Create DB
			err = DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", config.DatabaseName)).Error
			if err != nil {
				log.Fatalf("Failed to create DataBase: %v", err)
			}
			log.Infof("Create DB %s Successfully", config.DatabaseName)

			// Use it
			SQLUseDB := fmt.Sprintf(`
					USE %s
				`, config.DatabaseName)
			err = DB.Exec(SQLUseDB).Error
			if err != nil {
				log.Fatalf("Failed to use database %v: %v", config.DatabaseName, err)
			}
		}
	}
	return DB
}
