package models

import "gorm.io/gorm"

type ControlCommand struct {
	gorm.Model
	ClientID  uint   `json:"clientId"`
	Command   string `json:"command"`
	Arguments string `json:"arguments" gorm:"serializer:json"`
	Status    string `json:"status" gorm:"default:'PENDING'"`
	Output    string `json:"output"`
}
