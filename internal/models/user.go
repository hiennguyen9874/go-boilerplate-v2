package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id                 uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name               string     `gorm:"type:varchar(100);not null"`
	Email              string     `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password           string     `gorm:"type:varchar(100);not null"`
	CreatedAt          time.Time  `gorm:"not null;default:now()"`
	UpdatedAt          time.Time  `gorm:"not null;default:now()"`
	IsActive           bool       `gorm:"not null;default:true"`
	IsSuperUser        bool       `gorm:"not null;default:false"`
	Verified           bool       `gorm:"not null;default:false"`
	VerificationCode   *string    `gorm:"type:varchar(32);default:null"`
	PasswordResetToken *string    `gorm:"type:varchar(32);default:null"`
	PasswordResetAt    *time.Time `gorm:"default:null"`
	Items              []Item     `gorm:"foreignKey:OwnerId;references:Id"`
}

func (User) TableName() string {
	return "user"
}
