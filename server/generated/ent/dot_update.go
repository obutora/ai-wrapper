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
	"github.com/obutora/go_graphql_template/generated/ent/dot"
	"github.com/obutora/go_graphql_template/generated/ent/predicate"
)

// DotUpdate is the builder for updating Dot entities.
type DotUpdate struct {
	config
	hooks     []Hook
	mutation  *DotMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the DotUpdate builder.
func (du *DotUpdate) Where(ps ...predicate.Dot) *DotUpdate {
	du.mutation.Where(ps...)
	return du
}

// SetUID sets the "uid" field.
func (du *DotUpdate) SetUID(s string) *DotUpdate {
	du.mutation.SetUID(s)
	return du
}

// SetNillableUID sets the "uid" field if the given value is not nil.
func (du *DotUpdate) SetNillableUID(s *string) *DotUpdate {
	if s != nil {
		du.SetUID(*s)
	}
	return du
}

// SetTitle sets the "title" field.
func (du *DotUpdate) SetTitle(s string) *DotUpdate {
	du.mutation.SetTitle(s)
	return du
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (du *DotUpdate) SetNillableTitle(s *string) *DotUpdate {
	if s != nil {
		du.SetTitle(*s)
	}
	return du
}

// SetCreatedAt sets the "created_at" field.
func (du *DotUpdate) SetCreatedAt(t time.Time) *DotUpdate {
	du.mutation.SetCreatedAt(t)
	return du
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (du *DotUpdate) SetNillableCreatedAt(t *time.Time) *DotUpdate {
	if t != nil {
		du.SetCreatedAt(*t)
	}
	return du
}

// SetUpdatedAt sets the "updated_at" field.
func (du *DotUpdate) SetUpdatedAt(t time.Time) *DotUpdate {
	du.mutation.SetUpdatedAt(t)
	return du
}

// Mutation returns the DotMutation object of the builder.
func (du *DotUpdate) Mutation() *DotMutation {
	return du.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (du *DotUpdate) Save(ctx context.Context) (int, error) {
	du.defaults()
	return withHooks(ctx, du.sqlSave, du.mutation, du.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (du *DotUpdate) SaveX(ctx context.Context) int {
	affected, err := du.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (du *DotUpdate) Exec(ctx context.Context) error {
	_, err := du.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (du *DotUpdate) ExecX(ctx context.Context) {
	if err := du.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (du *DotUpdate) defaults() {
	if _, ok := du.mutation.UpdatedAt(); !ok {
		v := dot.UpdateDefaultUpdatedAt()
		du.mutation.SetUpdatedAt(v)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (du *DotUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *DotUpdate {
	du.modifiers = append(du.modifiers, modifiers...)
	return du
}

func (du *DotUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(dot.Table, dot.Columns, sqlgraph.NewFieldSpec(dot.FieldID, field.TypeInt))
	if ps := du.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := du.mutation.UID(); ok {
		_spec.SetField(dot.FieldUID, field.TypeString, value)
	}
	if value, ok := du.mutation.Title(); ok {
		_spec.SetField(dot.FieldTitle, field.TypeString, value)
	}
	if value, ok := du.mutation.CreatedAt(); ok {
		_spec.SetField(dot.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := du.mutation.UpdatedAt(); ok {
		_spec.SetField(dot.FieldUpdatedAt, field.TypeTime, value)
	}
	_spec.AddModifiers(du.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, du.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{dot.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	du.mutation.done = true
	return n, nil
}

// DotUpdateOne is the builder for updating a single Dot entity.
type DotUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *DotMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetUID sets the "uid" field.
func (duo *DotUpdateOne) SetUID(s string) *DotUpdateOne {
	duo.mutation.SetUID(s)
	return duo
}

// SetNillableUID sets the "uid" field if the given value is not nil.
func (duo *DotUpdateOne) SetNillableUID(s *string) *DotUpdateOne {
	if s != nil {
		duo.SetUID(*s)
	}
	return duo
}

// SetTitle sets the "title" field.
func (duo *DotUpdateOne) SetTitle(s string) *DotUpdateOne {
	duo.mutation.SetTitle(s)
	return duo
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (duo *DotUpdateOne) SetNillableTitle(s *string) *DotUpdateOne {
	if s != nil {
		duo.SetTitle(*s)
	}
	return duo
}

// SetCreatedAt sets the "created_at" field.
func (duo *DotUpdateOne) SetCreatedAt(t time.Time) *DotUpdateOne {
	duo.mutation.SetCreatedAt(t)
	return duo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (duo *DotUpdateOne) SetNillableCreatedAt(t *time.Time) *DotUpdateOne {
	if t != nil {
		duo.SetCreatedAt(*t)
	}
	return duo
}

// SetUpdatedAt sets the "updated_at" field.
func (duo *DotUpdateOne) SetUpdatedAt(t time.Time) *DotUpdateOne {
	duo.mutation.SetUpdatedAt(t)
	return duo
}

// Mutation returns the DotMutation object of the builder.
func (duo *DotUpdateOne) Mutation() *DotMutation {
	return duo.mutation
}

// Where appends a list predicates to the DotUpdate builder.
func (duo *DotUpdateOne) Where(ps ...predicate.Dot) *DotUpdateOne {
	duo.mutation.Where(ps...)
	return duo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (duo *DotUpdateOne) Select(field string, fields ...string) *DotUpdateOne {
	duo.fields = append([]string{field}, fields...)
	return duo
}

// Save executes the query and returns the updated Dot entity.
func (duo *DotUpdateOne) Save(ctx context.Context) (*Dot, error) {
	duo.defaults()
	return withHooks(ctx, duo.sqlSave, duo.mutation, duo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (duo *DotUpdateOne) SaveX(ctx context.Context) *Dot {
	node, err := duo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (duo *DotUpdateOne) Exec(ctx context.Context) error {
	_, err := duo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (duo *DotUpdateOne) ExecX(ctx context.Context) {
	if err := duo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (duo *DotUpdateOne) defaults() {
	if _, ok := duo.mutation.UpdatedAt(); !ok {
		v := dot.UpdateDefaultUpdatedAt()
		duo.mutation.SetUpdatedAt(v)
	}
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (duo *DotUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *DotUpdateOne {
	duo.modifiers = append(duo.modifiers, modifiers...)
	return duo
}

func (duo *DotUpdateOne) sqlSave(ctx context.Context) (_node *Dot, err error) {
	_spec := sqlgraph.NewUpdateSpec(dot.Table, dot.Columns, sqlgraph.NewFieldSpec(dot.FieldID, field.TypeInt))
	id, ok := duo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Dot.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := duo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, dot.FieldID)
		for _, f := range fields {
			if !dot.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != dot.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := duo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := duo.mutation.UID(); ok {
		_spec.SetField(dot.FieldUID, field.TypeString, value)
	}
	if value, ok := duo.mutation.Title(); ok {
		_spec.SetField(dot.FieldTitle, field.TypeString, value)
	}
	if value, ok := duo.mutation.CreatedAt(); ok {
		_spec.SetField(dot.FieldCreatedAt, field.TypeTime, value)
	}
	if value, ok := duo.mutation.UpdatedAt(); ok {
		_spec.SetField(dot.FieldUpdatedAt, field.TypeTime, value)
	}
	_spec.AddModifiers(duo.modifiers...)
	_node = &Dot{config: duo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, duo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{dot.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	duo.mutation.done = true
	return _node, nil
}
