// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/ent/auth"
	"github.com/Ali-Gorgani/chat-room-project/services/auth-service/utils/ent/predicate"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeAuth = "Auth"
)

// AuthMutation represents an operation that mutates the Auth nodes in the graph.
type AuthMutation struct {
	config
	op            Op
	typ           string
	id            *string
	user_id       *uint
	adduser_id    *int
	refresh_token *string
	is_revoked    *bool
	created_at    *time.Time
	expires_at    *time.Time
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Auth, error)
	predicates    []predicate.Auth
}

var _ ent.Mutation = (*AuthMutation)(nil)

// authOption allows management of the mutation configuration using functional options.
type authOption func(*AuthMutation)

// newAuthMutation creates new mutation for the Auth entity.
func newAuthMutation(c config, op Op, opts ...authOption) *AuthMutation {
	m := &AuthMutation{
		config:        c,
		op:            op,
		typ:           TypeAuth,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withAuthID sets the ID field of the mutation.
func withAuthID(id string) authOption {
	return func(m *AuthMutation) {
		var (
			err   error
			once  sync.Once
			value *Auth
		)
		m.oldValue = func(ctx context.Context) (*Auth, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Auth.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withAuth sets the old Auth of the mutation.
func withAuth(node *Auth) authOption {
	return func(m *AuthMutation) {
		m.oldValue = func(context.Context) (*Auth, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m AuthMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m AuthMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Auth entities.
func (m *AuthMutation) SetID(id string) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *AuthMutation) ID() (id string, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *AuthMutation) IDs(ctx context.Context) ([]string, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []string{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Auth.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetUserID sets the "user_id" field.
func (m *AuthMutation) SetUserID(u uint) {
	m.user_id = &u
	m.adduser_id = nil
}

// UserID returns the value of the "user_id" field in the mutation.
func (m *AuthMutation) UserID() (r uint, exists bool) {
	v := m.user_id
	if v == nil {
		return
	}
	return *v, true
}

// OldUserID returns the old "user_id" field's value of the Auth entity.
// If the Auth object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AuthMutation) OldUserID(ctx context.Context) (v uint, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUserID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUserID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUserID: %w", err)
	}
	return oldValue.UserID, nil
}

// AddUserID adds u to the "user_id" field.
func (m *AuthMutation) AddUserID(u int) {
	if m.adduser_id != nil {
		*m.adduser_id += u
	} else {
		m.adduser_id = &u
	}
}

// AddedUserID returns the value that was added to the "user_id" field in this mutation.
func (m *AuthMutation) AddedUserID() (r int, exists bool) {
	v := m.adduser_id
	if v == nil {
		return
	}
	return *v, true
}

// ResetUserID resets all changes to the "user_id" field.
func (m *AuthMutation) ResetUserID() {
	m.user_id = nil
	m.adduser_id = nil
}

// SetRefreshToken sets the "refresh_token" field.
func (m *AuthMutation) SetRefreshToken(s string) {
	m.refresh_token = &s
}

// RefreshToken returns the value of the "refresh_token" field in the mutation.
func (m *AuthMutation) RefreshToken() (r string, exists bool) {
	v := m.refresh_token
	if v == nil {
		return
	}
	return *v, true
}

// OldRefreshToken returns the old "refresh_token" field's value of the Auth entity.
// If the Auth object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AuthMutation) OldRefreshToken(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRefreshToken is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRefreshToken requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRefreshToken: %w", err)
	}
	return oldValue.RefreshToken, nil
}

// ResetRefreshToken resets all changes to the "refresh_token" field.
func (m *AuthMutation) ResetRefreshToken() {
	m.refresh_token = nil
}

// SetIsRevoked sets the "is_revoked" field.
func (m *AuthMutation) SetIsRevoked(b bool) {
	m.is_revoked = &b
}

// IsRevoked returns the value of the "is_revoked" field in the mutation.
func (m *AuthMutation) IsRevoked() (r bool, exists bool) {
	v := m.is_revoked
	if v == nil {
		return
	}
	return *v, true
}

// OldIsRevoked returns the old "is_revoked" field's value of the Auth entity.
// If the Auth object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AuthMutation) OldIsRevoked(ctx context.Context) (v bool, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldIsRevoked is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldIsRevoked requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldIsRevoked: %w", err)
	}
	return oldValue.IsRevoked, nil
}

// ResetIsRevoked resets all changes to the "is_revoked" field.
func (m *AuthMutation) ResetIsRevoked() {
	m.is_revoked = nil
}

// SetCreatedAt sets the "created_at" field.
func (m *AuthMutation) SetCreatedAt(t time.Time) {
	m.created_at = &t
}

// CreatedAt returns the value of the "created_at" field in the mutation.
func (m *AuthMutation) CreatedAt() (r time.Time, exists bool) {
	v := m.created_at
	if v == nil {
		return
	}
	return *v, true
}

// OldCreatedAt returns the old "created_at" field's value of the Auth entity.
// If the Auth object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AuthMutation) OldCreatedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldCreatedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldCreatedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldCreatedAt: %w", err)
	}
	return oldValue.CreatedAt, nil
}

// ResetCreatedAt resets all changes to the "created_at" field.
func (m *AuthMutation) ResetCreatedAt() {
	m.created_at = nil
}

// SetExpiresAt sets the "expires_at" field.
func (m *AuthMutation) SetExpiresAt(t time.Time) {
	m.expires_at = &t
}

// ExpiresAt returns the value of the "expires_at" field in the mutation.
func (m *AuthMutation) ExpiresAt() (r time.Time, exists bool) {
	v := m.expires_at
	if v == nil {
		return
	}
	return *v, true
}

// OldExpiresAt returns the old "expires_at" field's value of the Auth entity.
// If the Auth object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *AuthMutation) OldExpiresAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldExpiresAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldExpiresAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldExpiresAt: %w", err)
	}
	return oldValue.ExpiresAt, nil
}

// ResetExpiresAt resets all changes to the "expires_at" field.
func (m *AuthMutation) ResetExpiresAt() {
	m.expires_at = nil
}

// Where appends a list predicates to the AuthMutation builder.
func (m *AuthMutation) Where(ps ...predicate.Auth) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the AuthMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *AuthMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Auth, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *AuthMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *AuthMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Auth).
func (m *AuthMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *AuthMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m.user_id != nil {
		fields = append(fields, auth.FieldUserID)
	}
	if m.refresh_token != nil {
		fields = append(fields, auth.FieldRefreshToken)
	}
	if m.is_revoked != nil {
		fields = append(fields, auth.FieldIsRevoked)
	}
	if m.created_at != nil {
		fields = append(fields, auth.FieldCreatedAt)
	}
	if m.expires_at != nil {
		fields = append(fields, auth.FieldExpiresAt)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *AuthMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case auth.FieldUserID:
		return m.UserID()
	case auth.FieldRefreshToken:
		return m.RefreshToken()
	case auth.FieldIsRevoked:
		return m.IsRevoked()
	case auth.FieldCreatedAt:
		return m.CreatedAt()
	case auth.FieldExpiresAt:
		return m.ExpiresAt()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *AuthMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case auth.FieldUserID:
		return m.OldUserID(ctx)
	case auth.FieldRefreshToken:
		return m.OldRefreshToken(ctx)
	case auth.FieldIsRevoked:
		return m.OldIsRevoked(ctx)
	case auth.FieldCreatedAt:
		return m.OldCreatedAt(ctx)
	case auth.FieldExpiresAt:
		return m.OldExpiresAt(ctx)
	}
	return nil, fmt.Errorf("unknown Auth field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *AuthMutation) SetField(name string, value ent.Value) error {
	switch name {
	case auth.FieldUserID:
		v, ok := value.(uint)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUserID(v)
		return nil
	case auth.FieldRefreshToken:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRefreshToken(v)
		return nil
	case auth.FieldIsRevoked:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIsRevoked(v)
		return nil
	case auth.FieldCreatedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreatedAt(v)
		return nil
	case auth.FieldExpiresAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetExpiresAt(v)
		return nil
	}
	return fmt.Errorf("unknown Auth field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *AuthMutation) AddedFields() []string {
	var fields []string
	if m.adduser_id != nil {
		fields = append(fields, auth.FieldUserID)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *AuthMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case auth.FieldUserID:
		return m.AddedUserID()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *AuthMutation) AddField(name string, value ent.Value) error {
	switch name {
	case auth.FieldUserID:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddUserID(v)
		return nil
	}
	return fmt.Errorf("unknown Auth numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *AuthMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *AuthMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *AuthMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Auth nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *AuthMutation) ResetField(name string) error {
	switch name {
	case auth.FieldUserID:
		m.ResetUserID()
		return nil
	case auth.FieldRefreshToken:
		m.ResetRefreshToken()
		return nil
	case auth.FieldIsRevoked:
		m.ResetIsRevoked()
		return nil
	case auth.FieldCreatedAt:
		m.ResetCreatedAt()
		return nil
	case auth.FieldExpiresAt:
		m.ResetExpiresAt()
		return nil
	}
	return fmt.Errorf("unknown Auth field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *AuthMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *AuthMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *AuthMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *AuthMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *AuthMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *AuthMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *AuthMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Auth unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *AuthMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Auth edge %s", name)
}
