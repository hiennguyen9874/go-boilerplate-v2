package models

import (
	"time"
)

type User struct {
	Id                 uint
	CreateTime         time.Time
	UpdateTime         time.Time
	Name               string
	Email              string
	Password           string
	IsActive           bool
	IsSuperUser        bool
	Verified           bool
	VerificationCode   *string
	PasswordResetToken *string
	PasswordResetAt    *time.Time
}

type UserCreate struct {
	Name        string
	Email       string
	Password    string
	IsActive    *bool
	IsSuperUser *bool
}

type UserUpdate struct {
	Name               *string
	Email              *string
	Password           *string
	IsActive           *bool
	IsSuperUser        *bool
	Verified           *bool
	VerificationCode   *string
	PasswordResetToken *string
	PasswordResetAt    *time.Time
}
