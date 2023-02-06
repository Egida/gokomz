package db

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/quietjoy/gocom/pkg/models"
)

type GokomzStore interface {
	GetAllClients() ([]models.Client, error)
}

func New() *gorm.DB {
	dsn := GenerateMysqlDSN("mysql", "3306", "root", "password", "gocom")
	mysql, err := openDatabase(dsn)
	if err != nil {
		log.Fatal("Error connecting to DB: ", err)
	}
	log.Println("Successfully connected to DB")

	err = Migrate(mysql, &models.Client{}, &models.ControlCommand{})
	if err != nil {
		log.Fatal("Error migrating models: ", err)
	}
	log.Println("Successfully migrated models")
	return mysql
}

func openDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GenerateMysqlDSN(host, port, user, password, database string) string {
	return user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func Migrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}
