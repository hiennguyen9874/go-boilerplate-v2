package models

type Item struct {
	Id          uint
	Title       string
	Description string
	OwnerId     uint
}

type ItemCreate struct {
	Title       string
	Description string
}

type ItemUpdate struct {
	Title       *string
	Description *string
}
