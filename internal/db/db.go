package db

import (
	"fmt"
	"github.com/k4zb3k/pethub/pkg/logging"

	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var logger = logging.GetLogger()

func GetDbConnection() (*gorm.DB, error) {
	host := "localhost"
	port := "5432"
	user := "humo"
	password := "pass"
	dbname := "humo_db"

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Dushanbe",
		host, port, user, password, dbname)

	conn, err := gorm.Open(postgresDriver.Open(connString))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info("Successful connection to DB", host)

	return conn, nil
}
