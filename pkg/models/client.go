package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPList []string

type Client struct {
	gorm.Model
	UUID        uuid.UUID        `json:"uuid"`
	SourceIP    string           `json:"sourceIP"`
	ForwardedIP IPList           `json:"forwardedIP" gorm:"serializer:json"`
	Commands    []ControlCommand `json:"commands"`
}
