package models

type Tabler interface {
	// TableName overrides the table name used by User to `profiles`
	TableName() string
}
