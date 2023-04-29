// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/hiennguyen9874/go-boilerplate-v2/ent/item"
	"github.com/hiennguyen9874/go-boilerplate-v2/ent/predicate"
	"github.com/hiennguyen9874/go-boilerplate-v2/ent/user"
)

// UserUpdate is the builder for updating User entities.
type UserUpdate struct {
	config
	hooks    []Hook
	mutation *UserMutation
}

// Where appends a list predicates to the UserUpdate builder.
func (uu *UserUpdate) Where(ps ...predicate.User) *UserUpdate {
	uu.mutation.Where(ps...)
	return uu
}

// SetUpdateTime sets the "update_time" field.
func (uu *UserUpdate) SetUpdateTime(t time.Time) *UserUpdate {
	uu.mutation.SetUpdateTime(t)
	return uu
}

// SetName sets the "name" field.
func (uu *UserUpdate) SetName(s string) *UserUpdate {
	uu.mutation.SetName(s)
	return uu
}

// SetEmail sets the "email" field.
func (uu *UserUpdate) SetEmail(s string) *UserUpdate {
	uu.mutation.SetEmail(s)
	return uu
}

// SetPassword sets the "password" field.
func (uu *UserUpdate) SetPassword(s string) *UserUpdate {
	uu.mutation.SetPassword(s)
	return uu
}

// SetIsActive sets the "is_active" field.
func (uu *UserUpdate) SetIsActive(b bool) *UserUpdate {
	uu.mutation.SetIsActive(b)
	return uu
}

// SetNillableIsActive sets the "is_active" field if the given value is not nil.
func (uu *UserUpdate) SetNillableIsActive(b *bool) *UserUpdate {
	if b != nil {
		uu.SetIsActive(*b)
	}
	return uu
}

// SetIsSuperUser sets the "is_super_user" field.
func (uu *UserUpdate) SetIsSuperUser(b bool) *UserUpdate {
	uu.mutation.SetIsSuperUser(b)
	return uu
}

// SetNillableIsSuperUser sets the "is_super_user" field if the given value is not nil.
func (uu *UserUpdate) SetNillableIsSuperUser(b *bool) *UserUpdate {
	if b != nil {
		uu.SetIsSuperUser(*b)
	}
	return uu
}

// SetVerified sets the "verified" field.
func (uu *UserUpdate) SetVerified(b bool) *UserUpdate {
	uu.mutation.SetVerified(b)
	return uu
}

// SetNillableVerified sets the "verified" field if the given value is not nil.
func (uu *UserUpdate) SetNillableVerified(b *bool) *UserUpdate {
	if b != nil {
		uu.SetVerified(*b)
	}
	return uu
}

// SetVerificationCode sets the "verification_code" field.
func (uu *UserUpdate) SetVerificationCode(s string) *UserUpdate {
	uu.mutation.SetVerificationCode(s)
	return uu
}

// SetNillableVerificationCode sets the "verification_code" field if the given value is not nil.
func (uu *UserUpdate) SetNillableVerificationCode(s *string) *UserUpdate {
	if s != nil {
		uu.SetVerificationCode(*s)
	}
	return uu
}

// ClearVerificationCode clears the value of the "verification_code" field.
func (uu *UserUpdate) ClearVerificationCode() *UserUpdate {
	uu.mutation.ClearVerificationCode()
	return uu
}

// SetPasswordResetToken sets the "password_reset_token" field.
func (uu *UserUpdate) SetPasswordResetToken(s string) *UserUpdate {
	uu.mutation.SetPasswordResetToken(s)
	return uu
}

// SetNillablePasswordResetToken sets the "password_reset_token" field if the given value is not nil.
func (uu *UserUpdate) SetNillablePasswordResetToken(s *string) *UserUpdate {
	if s != nil {
		uu.SetPasswordResetToken(*s)
	}
	return uu
}

// ClearPasswordResetToken clears the value of the "password_reset_token" field.
func (uu *UserUpdate) ClearPasswordResetToken() *UserUpdate {
	uu.mutation.ClearPasswordResetToken()
	return uu
}

// SetPasswordResetAt sets the "password_reset_at" field.
func (uu *UserUpdate) SetPasswordResetAt(t time.Time) *UserUpdate {
	uu.mutation.SetPasswordResetAt(t)
	return uu
}

// SetNillablePasswordResetAt sets the "password_reset_at" field if the given value is not nil.
func (uu *UserUpdate) SetNillablePasswordResetAt(t *time.Time) *UserUpdate {
	if t != nil {
		uu.SetPasswordResetAt(*t)
	}
	return uu
}

// ClearPasswordResetAt clears the value of the "password_reset_at" field.
func (uu *UserUpdate) ClearPasswordResetAt() *UserUpdate {
	uu.mutation.ClearPasswordResetAt()
	return uu
}

// AddItemIDs adds the "items" edge to the Item entity by IDs.
func (uu *UserUpdate) AddItemIDs(ids ...uint) *UserUpdate {
	uu.mutation.AddItemIDs(ids...)
	return uu
}

// AddItems adds the "items" edges to the Item entity.
func (uu *UserUpdate) AddItems(i ...*Item) *UserUpdate {
	ids := make([]uint, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return uu.AddItemIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uu *UserUpdate) Mutation() *UserMutation {
	return uu.mutation
}

// ClearItems clears all "items" edges to the Item entity.
func (uu *UserUpdate) ClearItems() *UserUpdate {
	uu.mutation.ClearItems()
	return uu
}

// RemoveItemIDs removes the "items" edge to Item entities by IDs.
func (uu *UserUpdate) RemoveItemIDs(ids ...uint) *UserUpdate {
	uu.mutation.RemoveItemIDs(ids...)
	return uu
}

// RemoveItems removes "items" edges to Item entities.
func (uu *UserUpdate) RemoveItems(i ...*Item) *UserUpdate {
	ids := make([]uint, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return uu.RemoveItemIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (uu *UserUpdate) Save(ctx context.Context) (int, error) {
	uu.defaults()
	return withHooks[int, UserMutation](ctx, uu.sqlSave, uu.mutation, uu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uu *UserUpdate) SaveX(ctx context.Context) int {
	affected, err := uu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uu *UserUpdate) Exec(ctx context.Context) error {
	_, err := uu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uu *UserUpdate) ExecX(ctx context.Context) {
	if err := uu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uu *UserUpdate) defaults() {
	if _, ok := uu.mutation.UpdateTime(); !ok {
		v := user.UpdateDefaultUpdateTime()
		uu.mutation.SetUpdateTime(v)
	}
}

func (uu *UserUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(user.Table, user.Columns, sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint))
	if ps := uu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uu.mutation.UpdateTime(); ok {
		_spec.SetField(user.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := uu.mutation.Name(); ok {
		_spec.SetField(user.FieldName, field.TypeString, value)
	}
	if value, ok := uu.mutation.Email(); ok {
		_spec.SetField(user.FieldEmail, field.TypeString, value)
	}
	if value, ok := uu.mutation.Password(); ok {
		_spec.SetField(user.FieldPassword, field.TypeString, value)
	}
	if value, ok := uu.mutation.IsActive(); ok {
		_spec.SetField(user.FieldIsActive, field.TypeBool, value)
	}
	if value, ok := uu.mutation.IsSuperUser(); ok {
		_spec.SetField(user.FieldIsSuperUser, field.TypeBool, value)
	}
	if value, ok := uu.mutation.Verified(); ok {
		_spec.SetField(user.FieldVerified, field.TypeBool, value)
	}
	if value, ok := uu.mutation.VerificationCode(); ok {
		_spec.SetField(user.FieldVerificationCode, field.TypeString, value)
	}
	if uu.mutation.VerificationCodeCleared() {
		_spec.ClearField(user.FieldVerificationCode, field.TypeString)
	}
	if value, ok := uu.mutation.PasswordResetToken(); ok {
		_spec.SetField(user.FieldPasswordResetToken, field.TypeString, value)
	}
	if uu.mutation.PasswordResetTokenCleared() {
		_spec.ClearField(user.FieldPasswordResetToken, field.TypeString)
	}
	if value, ok := uu.mutation.PasswordResetAt(); ok {
		_spec.SetField(user.FieldPasswordResetAt, field.TypeTime, value)
	}
	if uu.mutation.PasswordResetAtCleared() {
		_spec.ClearField(user.FieldPasswordResetAt, field.TypeTime)
	}
	if uu.mutation.ItemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.ItemsTable,
			Columns: []string{user.ItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeUint),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.RemovedItemsIDs(); len(nodes) > 0 && !uu.mutation.ItemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.ItemsTable,
			Columns: []string{user.ItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeUint),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.ItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.ItemsTable,
			Columns: []string{user.ItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeUint),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, uu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	uu.mutation.done = true
	return n, nil
}

// UserUpdateOne is the builder for updating a single User entity.
type UserUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *UserMutation
}

// SetUpdateTime sets the "update_time" field.
func (uuo *UserUpdateOne) SetUpdateTime(t time.Time) *UserUpdateOne {
	uuo.mutation.SetUpdateTime(t)
	return uuo
}

// SetName sets the "name" field.
func (uuo *UserUpdateOne) SetName(s string) *UserUpdateOne {
	uuo.mutation.SetName(s)
	return uuo
}

// SetEmail sets the "email" field.
func (uuo *UserUpdateOne) SetEmail(s string) *UserUpdateOne {
	uuo.mutation.SetEmail(s)
	return uuo
}

// SetPassword sets the "password" field.
func (uuo *UserUpdateOne) SetPassword(s string) *UserUpdateOne {
	uuo.mutation.SetPassword(s)
	return uuo
}

// SetIsActive sets the "is_active" field.
func (uuo *UserUpdateOne) SetIsActive(b bool) *UserUpdateOne {
	uuo.mutation.SetIsActive(b)
	return uuo
}

// SetNillableIsActive sets the "is_active" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableIsActive(b *bool) *UserUpdateOne {
	if b != nil {
		uuo.SetIsActive(*b)
	}
	return uuo
}

// SetIsSuperUser sets the "is_super_user" field.
func (uuo *UserUpdateOne) SetIsSuperUser(b bool) *UserUpdateOne {
	uuo.mutation.SetIsSuperUser(b)
	return uuo
}

// SetNillableIsSuperUser sets the "is_super_user" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableIsSuperUser(b *bool) *UserUpdateOne {
	if b != nil {
		uuo.SetIsSuperUser(*b)
	}
	return uuo
}

// SetVerified sets the "verified" field.
func (uuo *UserUpdateOne) SetVerified(b bool) *UserUpdateOne {
	uuo.mutation.SetVerified(b)
	return uuo
}

// SetNillableVerified sets the "verified" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableVerified(b *bool) *UserUpdateOne {
	if b != nil {
		uuo.SetVerified(*b)
	}
	return uuo
}

// SetVerificationCode sets the "verification_code" field.
func (uuo *UserUpdateOne) SetVerificationCode(s string) *UserUpdateOne {
	uuo.mutation.SetVerificationCode(s)
	return uuo
}

// SetNillableVerificationCode sets the "verification_code" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableVerificationCode(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetVerificationCode(*s)
	}
	return uuo
}

// ClearVerificationCode clears the value of the "verification_code" field.
func (uuo *UserUpdateOne) ClearVerificationCode() *UserUpdateOne {
	uuo.mutation.ClearVerificationCode()
	return uuo
}

// SetPasswordResetToken sets the "password_reset_token" field.
func (uuo *UserUpdateOne) SetPasswordResetToken(s string) *UserUpdateOne {
	uuo.mutation.SetPasswordResetToken(s)
	return uuo
}

// SetNillablePasswordResetToken sets the "password_reset_token" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillablePasswordResetToken(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetPasswordResetToken(*s)
	}
	return uuo
}

// ClearPasswordResetToken clears the value of the "password_reset_token" field.
func (uuo *UserUpdateOne) ClearPasswordResetToken() *UserUpdateOne {
	uuo.mutation.ClearPasswordResetToken()
	return uuo
}

// SetPasswordResetAt sets the "password_reset_at" field.
func (uuo *UserUpdateOne) SetPasswordResetAt(t time.Time) *UserUpdateOne {
	uuo.mutation.SetPasswordResetAt(t)
	return uuo
}

// SetNillablePasswordResetAt sets the "password_reset_at" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillablePasswordResetAt(t *time.Time) *UserUpdateOne {
	if t != nil {
		uuo.SetPasswordResetAt(*t)
	}
	return uuo
}

// ClearPasswordResetAt clears the value of the "password_reset_at" field.
func (uuo *UserUpdateOne) ClearPasswordResetAt() *UserUpdateOne {
	uuo.mutation.ClearPasswordResetAt()
	return uuo
}

// AddItemIDs adds the "items" edge to the Item entity by IDs.
func (uuo *UserUpdateOne) AddItemIDs(ids ...uint) *UserUpdateOne {
	uuo.mutation.AddItemIDs(ids...)
	return uuo
}

// AddItems adds the "items" edges to the Item entity.
func (uuo *UserUpdateOne) AddItems(i ...*Item) *UserUpdateOne {
	ids := make([]uint, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return uuo.AddItemIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uuo *UserUpdateOne) Mutation() *UserMutation {
	return uuo.mutation
}

// ClearItems clears all "items" edges to the Item entity.
func (uuo *UserUpdateOne) ClearItems() *UserUpdateOne {
	uuo.mutation.ClearItems()
	return uuo
}

// RemoveItemIDs removes the "items" edge to Item entities by IDs.
func (uuo *UserUpdateOne) RemoveItemIDs(ids ...uint) *UserUpdateOne {
	uuo.mutation.RemoveItemIDs(ids...)
	return uuo
}

// RemoveItems removes "items" edges to Item entities.
func (uuo *UserUpdateOne) RemoveItems(i ...*Item) *UserUpdateOne {
	ids := make([]uint, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return uuo.RemoveItemIDs(ids...)
}

// Where appends a list predicates to the UserUpdate builder.
func (uuo *UserUpdateOne) Where(ps ...predicate.User) *UserUpdateOne {
	uuo.mutation.Where(ps...)
	return uuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (uuo *UserUpdateOne) Select(field string, fields ...string) *UserUpdateOne {
	uuo.fields = append([]string{field}, fields...)
	return uuo
}

// Save executes the query and returns the updated User entity.
func (uuo *UserUpdateOne) Save(ctx context.Context) (*User, error) {
	uuo.defaults()
	return withHooks[*User, UserMutation](ctx, uuo.sqlSave, uuo.mutation, uuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uuo *UserUpdateOne) SaveX(ctx context.Context) *User {
	node, err := uuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (uuo *UserUpdateOne) Exec(ctx context.Context) error {
	_, err := uuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uuo *UserUpdateOne) ExecX(ctx context.Context) {
	if err := uuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uuo *UserUpdateOne) defaults() {
	if _, ok := uuo.mutation.UpdateTime(); !ok {
		v := user.UpdateDefaultUpdateTime()
		uuo.mutation.SetUpdateTime(v)
	}
}

func (uuo *UserUpdateOne) sqlSave(ctx context.Context) (_node *User, err error) {
	_spec := sqlgraph.NewUpdateSpec(user.Table, user.Columns, sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint))
	id, ok := uuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "User.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := uuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, user.FieldID)
		for _, f := range fields {
			if !user.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != user.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := uuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uuo.mutation.UpdateTime(); ok {
		_spec.SetField(user.FieldUpdateTime, field.TypeTime, value)
	}
	if value, ok := uuo.mutation.Name(); ok {
		_spec.SetField(user.FieldName, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Email(); ok {
		_spec.SetField(user.FieldEmail, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Password(); ok {
		_spec.SetField(user.FieldPassword, field.TypeString, value)
	}
	if value, ok := uuo.mutation.IsActive(); ok {
		_spec.SetField(user.FieldIsActive, field.TypeBool, value)
	}
	if value, ok := uuo.mutation.IsSuperUser(); ok {
		_spec.SetField(user.FieldIsSuperUser, field.TypeBool, value)
	}
	if value, ok := uuo.mutation.Verified(); ok {
		_spec.SetField(user.FieldVerified, field.TypeBool, value)
	}
	if value, ok := uuo.mutation.VerificationCode(); ok {
		_spec.SetField(user.FieldVerificationCode, field.TypeString, value)
	}
	if uuo.mutation.VerificationCodeCleared() {
		_spec.ClearField(user.FieldVerificationCode, field.TypeString)
	}
	if value, ok := uuo.mutation.PasswordResetToken(); ok {
		_spec.SetField(user.FieldPasswordResetToken, field.TypeString, value)
	}
	if uuo.mutation.PasswordResetTokenCleared() {
		_spec.ClearField(user.FieldPasswordResetToken, field.TypeString)
	}
	if value, ok := uuo.mutation.PasswordResetAt(); ok {
		_spec.SetField(user.FieldPasswordResetAt, field.TypeTime, value)
	}
	if uuo.mutation.PasswordResetAtCleared() {
		_spec.ClearField(user.FieldPasswordResetAt, field.TypeTime)
	}
	if uuo.mutation.ItemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.ItemsTable,
			Columns: []string{user.ItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeUint),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.RemovedItemsIDs(); len(nodes) > 0 && !uuo.mutation.ItemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.ItemsTable,
			Columns: []string{user.ItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeUint),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.ItemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.ItemsTable,
			Columns: []string{user.ItemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(item.FieldID, field.TypeUint),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &User{config: uuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, uuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	uuo.mutation.done = true
	return _node, nil
}
