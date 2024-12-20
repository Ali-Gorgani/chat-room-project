// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/ent/auth"
)

// AuthCreate is the builder for creating a Auth entity.
type AuthCreate struct {
	config
	mutation *AuthMutation
	hooks    []Hook
}

// SetUserID sets the "user_id" field.
func (ac *AuthCreate) SetUserID(u uint) *AuthCreate {
	ac.mutation.SetUserID(u)
	return ac
}

// SetRefreshToken sets the "refresh_token" field.
func (ac *AuthCreate) SetRefreshToken(s string) *AuthCreate {
	ac.mutation.SetRefreshToken(s)
	return ac
}

// SetIsRevoked sets the "is_revoked" field.
func (ac *AuthCreate) SetIsRevoked(b bool) *AuthCreate {
	ac.mutation.SetIsRevoked(b)
	return ac
}

// SetNillableIsRevoked sets the "is_revoked" field if the given value is not nil.
func (ac *AuthCreate) SetNillableIsRevoked(b *bool) *AuthCreate {
	if b != nil {
		ac.SetIsRevoked(*b)
	}
	return ac
}

// SetCreatedAt sets the "created_at" field.
func (ac *AuthCreate) SetCreatedAt(t time.Time) *AuthCreate {
	ac.mutation.SetCreatedAt(t)
	return ac
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (ac *AuthCreate) SetNillableCreatedAt(t *time.Time) *AuthCreate {
	if t != nil {
		ac.SetCreatedAt(*t)
	}
	return ac
}

// SetExpiresAt sets the "expires_at" field.
func (ac *AuthCreate) SetExpiresAt(t time.Time) *AuthCreate {
	ac.mutation.SetExpiresAt(t)
	return ac
}

// SetID sets the "id" field.
func (ac *AuthCreate) SetID(s string) *AuthCreate {
	ac.mutation.SetID(s)
	return ac
}

// Mutation returns the AuthMutation object of the builder.
func (ac *AuthCreate) Mutation() *AuthMutation {
	return ac.mutation
}

// Save creates the Auth in the database.
func (ac *AuthCreate) Save(ctx context.Context) (*Auth, error) {
	ac.defaults()
	return withHooks(ctx, ac.sqlSave, ac.mutation, ac.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ac *AuthCreate) SaveX(ctx context.Context) *Auth {
	v, err := ac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ac *AuthCreate) Exec(ctx context.Context) error {
	_, err := ac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ac *AuthCreate) ExecX(ctx context.Context) {
	if err := ac.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ac *AuthCreate) defaults() {
	if _, ok := ac.mutation.IsRevoked(); !ok {
		v := auth.DefaultIsRevoked
		ac.mutation.SetIsRevoked(v)
	}
	if _, ok := ac.mutation.CreatedAt(); !ok {
		v := auth.DefaultCreatedAt()
		ac.mutation.SetCreatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ac *AuthCreate) check() error {
	if _, ok := ac.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "Auth.user_id"`)}
	}
	if v, ok := ac.mutation.UserID(); ok {
		if err := auth.UserIDValidator(v); err != nil {
			return &ValidationError{Name: "user_id", err: fmt.Errorf(`ent: validator failed for field "Auth.user_id": %w`, err)}
		}
	}
	if _, ok := ac.mutation.RefreshToken(); !ok {
		return &ValidationError{Name: "refresh_token", err: errors.New(`ent: missing required field "Auth.refresh_token"`)}
	}
	if _, ok := ac.mutation.IsRevoked(); !ok {
		return &ValidationError{Name: "is_revoked", err: errors.New(`ent: missing required field "Auth.is_revoked"`)}
	}
	if _, ok := ac.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Auth.created_at"`)}
	}
	if _, ok := ac.mutation.ExpiresAt(); !ok {
		return &ValidationError{Name: "expires_at", err: errors.New(`ent: missing required field "Auth.expires_at"`)}
	}
	if v, ok := ac.mutation.ID(); ok {
		if err := auth.IDValidator(v); err != nil {
			return &ValidationError{Name: "id", err: fmt.Errorf(`ent: validator failed for field "Auth.id": %w`, err)}
		}
	}
	return nil
}

func (ac *AuthCreate) sqlSave(ctx context.Context) (*Auth, error) {
	if err := ac.check(); err != nil {
		return nil, err
	}
	_node, _spec := ac.createSpec()
	if err := sqlgraph.CreateNode(ctx, ac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(string); ok {
			_node.ID = id
		} else {
			return nil, fmt.Errorf("unexpected Auth.ID type: %T", _spec.ID.Value)
		}
	}
	ac.mutation.id = &_node.ID
	ac.mutation.done = true
	return _node, nil
}

func (ac *AuthCreate) createSpec() (*Auth, *sqlgraph.CreateSpec) {
	var (
		_node = &Auth{config: ac.config}
		_spec = sqlgraph.NewCreateSpec(auth.Table, sqlgraph.NewFieldSpec(auth.FieldID, field.TypeString))
	)
	if id, ok := ac.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := ac.mutation.UserID(); ok {
		_spec.SetField(auth.FieldUserID, field.TypeUint, value)
		_node.UserID = value
	}
	if value, ok := ac.mutation.RefreshToken(); ok {
		_spec.SetField(auth.FieldRefreshToken, field.TypeString, value)
		_node.RefreshToken = value
	}
	if value, ok := ac.mutation.IsRevoked(); ok {
		_spec.SetField(auth.FieldIsRevoked, field.TypeBool, value)
		_node.IsRevoked = value
	}
	if value, ok := ac.mutation.CreatedAt(); ok {
		_spec.SetField(auth.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := ac.mutation.ExpiresAt(); ok {
		_spec.SetField(auth.FieldExpiresAt, field.TypeTime, value)
		_node.ExpiresAt = value
	}
	return _node, _spec
}

// AuthCreateBulk is the builder for creating many Auth entities in bulk.
type AuthCreateBulk struct {
	config
	err      error
	builders []*AuthCreate
}

// Save creates the Auth entities in the database.
func (acb *AuthCreateBulk) Save(ctx context.Context) ([]*Auth, error) {
	if acb.err != nil {
		return nil, acb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(acb.builders))
	nodes := make([]*Auth, len(acb.builders))
	mutators := make([]Mutator, len(acb.builders))
	for i := range acb.builders {
		func(i int, root context.Context) {
			builder := acb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*AuthMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, acb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, acb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, acb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (acb *AuthCreateBulk) SaveX(ctx context.Context) []*Auth {
	v, err := acb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (acb *AuthCreateBulk) Exec(ctx context.Context) error {
	_, err := acb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (acb *AuthCreateBulk) ExecX(ctx context.Context) {
	if err := acb.Exec(ctx); err != nil {
		panic(err)
	}
}
