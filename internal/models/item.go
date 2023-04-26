package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Id          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Title       string    `gorm:"type:varchar(100);not null"`
	Description string    `gorm:"type:varchar(200);not null"`
	OwnerId     uuid.UUID
}

func (Item) TableName() string {
	return "item"
}
