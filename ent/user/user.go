// Code generated by ent, DO NOT EDIT.

package user

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreateTime holds the string denoting the create_time field in the database.
	FieldCreateTime = "create_time"
	// FieldUpdateTime holds the string denoting the update_time field in the database.
	FieldUpdateTime = "update_time"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldPassword holds the string denoting the password field in the database.
	FieldPassword = "password"
	// FieldIsActive holds the string denoting the is_active field in the database.
	FieldIsActive = "is_active"
	// FieldIsSuperUser holds the string denoting the is_super_user field in the database.
	FieldIsSuperUser = "is_super_user"
	// FieldVerified holds the string denoting the verified field in the database.
	FieldVerified = "verified"
	// FieldVerificationCode holds the string denoting the verification_code field in the database.
	FieldVerificationCode = "verification_code"
	// FieldPasswordResetToken holds the string denoting the password_reset_token field in the database.
	FieldPasswordResetToken = "password_reset_token"
	// FieldPasswordResetAt holds the string denoting the password_reset_at field in the database.
	FieldPasswordResetAt = "password_reset_at"
	// EdgeItems holds the string denoting the items edge name in mutations.
	EdgeItems = "items"
	// Table holds the table name of the user in the database.
	Table = "users"
	// ItemsTable is the table that holds the items relation/edge.
	ItemsTable = "items"
	// ItemsInverseTable is the table name for the Item entity.
	// It exists in this package in order to avoid circular dependency with the "item" package.
	ItemsInverseTable = "items"
	// ItemsColumn is the table column denoting the items relation/edge.
	ItemsColumn = "owner_id"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldCreateTime,
	FieldUpdateTime,
	FieldName,
	FieldEmail,
	FieldPassword,
	FieldIsActive,
	FieldIsSuperUser,
	FieldVerified,
	FieldVerificationCode,
	FieldPasswordResetToken,
	FieldPasswordResetAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreateTime holds the default value on creation for the "create_time" field.
	DefaultCreateTime func() time.Time
	// DefaultUpdateTime holds the default value on creation for the "update_time" field.
	DefaultUpdateTime func() time.Time
	// UpdateDefaultUpdateTime holds the default value on update for the "update_time" field.
	UpdateDefaultUpdateTime func() time.Time
	// DefaultIsActive holds the default value on creation for the "is_active" field.
	DefaultIsActive bool
	// DefaultIsSuperUser holds the default value on creation for the "is_super_user" field.
	DefaultIsSuperUser bool
	// DefaultVerified holds the default value on creation for the "verified" field.
	DefaultVerified bool
)

// OrderOption defines the ordering options for the User queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreateTime orders the results by the create_time field.
func ByCreateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreateTime, opts...).ToFunc()
}

// ByUpdateTime orders the results by the update_time field.
func ByUpdateTime(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdateTime, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByEmail orders the results by the email field.
func ByEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmail, opts...).ToFunc()
}

// ByPassword orders the results by the password field.
func ByPassword(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPassword, opts...).ToFunc()
}

// ByIsActive orders the results by the is_active field.
func ByIsActive(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsActive, opts...).ToFunc()
}

// ByIsSuperUser orders the results by the is_super_user field.
func ByIsSuperUser(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsSuperUser, opts...).ToFunc()
}

// ByVerified orders the results by the verified field.
func ByVerified(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVerified, opts...).ToFunc()
}

// ByVerificationCode orders the results by the verification_code field.
func ByVerificationCode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVerificationCode, opts...).ToFunc()
}

// ByPasswordResetToken orders the results by the password_reset_token field.
func ByPasswordResetToken(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPasswordResetToken, opts...).ToFunc()
}

// ByPasswordResetAt orders the results by the password_reset_at field.
func ByPasswordResetAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPasswordResetAt, opts...).ToFunc()
}

// ByItemsCount orders the results by items count.
func ByItemsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newItemsStep(), opts...)
	}
}

// ByItems orders the results by items terms.
func ByItems(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newItemsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newItemsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ItemsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ItemsTable, ItemsColumn),
	)
}
