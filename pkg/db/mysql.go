package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/quietjoy/gocom/pkg/models"
)

func NewMysqlDB(dsn string) (*gorm.DB, error) {
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

func Create(db *gorm.DB, model interface{}) error {
	return db.Create(model).Error
}

func Update(db *gorm.DB, model interface{}) error {
	return db.Save(model).Error
}

func Delete(db *gorm.DB, model interface{}) error {
	return db.Delete(model).Error
}

func FindOne(db *gorm.DB, model interface{}) error {
	return db.First(&model).Error
}

func GetClients(db *gorm.DB) ([]models.Client, error) {
	var clients []models.Client
	err := db.Find(&clients).Error
	return clients, err
}

func GetNextCommand(db *gorm.DB, clientID uint) (models.ControlCommand, error) {
	var command models.ControlCommand
	err := db.Where("client_id = ? AND status = ? ORDER BY created_at desc", clientID, "PENDING").First(&command).Error
	return command, err
}

func GetCommands(db *gorm.DB, clientID uint) ([]models.ControlCommand, error) {
	var commands []models.ControlCommand
	// TODO: ORDER BY created_at desc
	err := db.Where("client_id = ? AND status = ?", clientID, "PENDING").Find(&commands).Error
	return commands, err
}

func GetClientByUUID(db *gorm.DB, uuid string) (models.Client, error) {
	var client models.Client
	err := db.Where("uuid = ?", uuid).First(&client).Error
	return client, err
}

func SaveCommand(db *gorm.DB, command models.ControlCommand) error {
	return db.Save(&command).Error
}

func GetCommandByID(db *gorm.DB, commandID uint) (models.ControlCommand, error) {
	var command models.ControlCommand
	err := db.Where("id = ?", commandID).First(&command).Error
	return command, err
}
