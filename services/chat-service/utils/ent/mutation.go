// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ent/message"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ent/predicate"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ent/room"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeMessage = "Message"
	TypeRoom    = "Room"
)

// MessageMutation represents an operation that mutates the Message nodes in the graph.
type MessageMutation struct {
	config
	op            Op
	typ           string
	id            *int
	content       *string
	room_id       *string
	username      *string
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Message, error)
	predicates    []predicate.Message
}

var _ ent.Mutation = (*MessageMutation)(nil)

// messageOption allows management of the mutation configuration using functional options.
type messageOption func(*MessageMutation)

// newMessageMutation creates new mutation for the Message entity.
func newMessageMutation(c config, op Op, opts ...messageOption) *MessageMutation {
	m := &MessageMutation{
		config:        c,
		op:            op,
		typ:           TypeMessage,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withMessageID sets the ID field of the mutation.
func withMessageID(id int) messageOption {
	return func(m *MessageMutation) {
		var (
			err   error
			once  sync.Once
			value *Message
		)
		m.oldValue = func(ctx context.Context) (*Message, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Message.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withMessage sets the old Message of the mutation.
func withMessage(node *Message) messageOption {
	return func(m *MessageMutation) {
		m.oldValue = func(context.Context) (*Message, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m MessageMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m MessageMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *MessageMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *MessageMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Message.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetContent sets the "content" field.
func (m *MessageMutation) SetContent(s string) {
	m.content = &s
}

// Content returns the value of the "content" field in the mutation.
func (m *MessageMutation) Content() (r string, exists bool) {
	v := m.content
	if v == nil {
		return
	}
	return *v, true
}

// OldContent returns the old "content" field's value of the Message entity.
// If the Message object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MessageMutation) OldContent(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldContent is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldContent requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldContent: %w", err)
	}
	return oldValue.Content, nil
}

// ResetContent resets all changes to the "content" field.
func (m *MessageMutation) ResetContent() {
	m.content = nil
}

// SetRoomID sets the "room_id" field.
func (m *MessageMutation) SetRoomID(s string) {
	m.room_id = &s
}

// RoomID returns the value of the "room_id" field in the mutation.
func (m *MessageMutation) RoomID() (r string, exists bool) {
	v := m.room_id
	if v == nil {
		return
	}
	return *v, true
}

// OldRoomID returns the old "room_id" field's value of the Message entity.
// If the Message object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MessageMutation) OldRoomID(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRoomID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRoomID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRoomID: %w", err)
	}
	return oldValue.RoomID, nil
}

// ResetRoomID resets all changes to the "room_id" field.
func (m *MessageMutation) ResetRoomID() {
	m.room_id = nil
}

// SetUsername sets the "username" field.
func (m *MessageMutation) SetUsername(s string) {
	m.username = &s
}

// Username returns the value of the "username" field in the mutation.
func (m *MessageMutation) Username() (r string, exists bool) {
	v := m.username
	if v == nil {
		return
	}
	return *v, true
}

// OldUsername returns the old "username" field's value of the Message entity.
// If the Message object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *MessageMutation) OldUsername(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldUsername is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldUsername requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldUsername: %w", err)
	}
	return oldValue.Username, nil
}

// ResetUsername resets all changes to the "username" field.
func (m *MessageMutation) ResetUsername() {
	m.username = nil
}

// Where appends a list predicates to the MessageMutation builder.
func (m *MessageMutation) Where(ps ...predicate.Message) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the MessageMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *MessageMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Message, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *MessageMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *MessageMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Message).
func (m *MessageMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *MessageMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.content != nil {
		fields = append(fields, message.FieldContent)
	}
	if m.room_id != nil {
		fields = append(fields, message.FieldRoomID)
	}
	if m.username != nil {
		fields = append(fields, message.FieldUsername)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *MessageMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case message.FieldContent:
		return m.Content()
	case message.FieldRoomID:
		return m.RoomID()
	case message.FieldUsername:
		return m.Username()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *MessageMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case message.FieldContent:
		return m.OldContent(ctx)
	case message.FieldRoomID:
		return m.OldRoomID(ctx)
	case message.FieldUsername:
		return m.OldUsername(ctx)
	}
	return nil, fmt.Errorf("unknown Message field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *MessageMutation) SetField(name string, value ent.Value) error {
	switch name {
	case message.FieldContent:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetContent(v)
		return nil
	case message.FieldRoomID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRoomID(v)
		return nil
	case message.FieldUsername:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUsername(v)
		return nil
	}
	return fmt.Errorf("unknown Message field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *MessageMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *MessageMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *MessageMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Message numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *MessageMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *MessageMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *MessageMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Message nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *MessageMutation) ResetField(name string) error {
	switch name {
	case message.FieldContent:
		m.ResetContent()
		return nil
	case message.FieldRoomID:
		m.ResetRoomID()
		return nil
	case message.FieldUsername:
		m.ResetUsername()
		return nil
	}
	return fmt.Errorf("unknown Message field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *MessageMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *MessageMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *MessageMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *MessageMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *MessageMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *MessageMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *MessageMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Message unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *MessageMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Message edge %s", name)
}

// RoomMutation represents an operation that mutates the Room nodes in the graph.
type RoomMutation struct {
	config
	op            Op
	typ           string
	id            *int
	name          *string
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Room, error)
	predicates    []predicate.Room
}

var _ ent.Mutation = (*RoomMutation)(nil)

// roomOption allows management of the mutation configuration using functional options.
type roomOption func(*RoomMutation)

// newRoomMutation creates new mutation for the Room entity.
func newRoomMutation(c config, op Op, opts ...roomOption) *RoomMutation {
	m := &RoomMutation{
		config:        c,
		op:            op,
		typ:           TypeRoom,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withRoomID sets the ID field of the mutation.
func withRoomID(id int) roomOption {
	return func(m *RoomMutation) {
		var (
			err   error
			once  sync.Once
			value *Room
		)
		m.oldValue = func(ctx context.Context) (*Room, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Room.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withRoom sets the old Room of the mutation.
func withRoom(node *Room) roomOption {
	return func(m *RoomMutation) {
		m.oldValue = func(context.Context) (*Room, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m RoomMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m RoomMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *RoomMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *RoomMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Room.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetName sets the "name" field.
func (m *RoomMutation) SetName(s string) {
	m.name = &s
}

// Name returns the value of the "name" field in the mutation.
func (m *RoomMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// OldName returns the old "name" field's value of the Room entity.
// If the Room object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RoomMutation) OldName(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldName is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldName requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldName: %w", err)
	}
	return oldValue.Name, nil
}

// ResetName resets all changes to the "name" field.
func (m *RoomMutation) ResetName() {
	m.name = nil
}

// Where appends a list predicates to the RoomMutation builder.
func (m *RoomMutation) Where(ps ...predicate.Room) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the RoomMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *RoomMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Room, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *RoomMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *RoomMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Room).
func (m *RoomMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *RoomMutation) Fields() []string {
	fields := make([]string, 0, 1)
	if m.name != nil {
		fields = append(fields, room.FieldName)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *RoomMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case room.FieldName:
		return m.Name()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *RoomMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case room.FieldName:
		return m.OldName(ctx)
	}
	return nil, fmt.Errorf("unknown Room field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *RoomMutation) SetField(name string, value ent.Value) error {
	switch name {
	case room.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	}
	return fmt.Errorf("unknown Room field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *RoomMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *RoomMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *RoomMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Room numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *RoomMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *RoomMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *RoomMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Room nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *RoomMutation) ResetField(name string) error {
	switch name {
	case room.FieldName:
		m.ResetName()
		return nil
	}
	return fmt.Errorf("unknown Room field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *RoomMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *RoomMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *RoomMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *RoomMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *RoomMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *RoomMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *RoomMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Room unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *RoomMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Room edge %s", name)
}
