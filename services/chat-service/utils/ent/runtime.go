// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ent/message"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ent/room"
	"github.com/Ali-Gorgani/chat-room-project/services/chat-service/utils/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	messageFields := schema.Message{}.Fields()
	_ = messageFields
	// messageDescContent is the schema descriptor for content field.
	messageDescContent := messageFields[0].Descriptor()
	// message.ContentValidator is a validator for the "content" field. It is called by the builders before save.
	message.ContentValidator = messageDescContent.Validators[0].(func(string) error)
	// messageDescRoomID is the schema descriptor for room_id field.
	messageDescRoomID := messageFields[1].Descriptor()
	// message.RoomIDValidator is a validator for the "room_id" field. It is called by the builders before save.
	message.RoomIDValidator = messageDescRoomID.Validators[0].(func(string) error)
	// messageDescUsername is the schema descriptor for username field.
	messageDescUsername := messageFields[2].Descriptor()
	// message.UsernameValidator is a validator for the "username" field. It is called by the builders before save.
	message.UsernameValidator = messageDescUsername.Validators[0].(func(string) error)
	roomFields := schema.Room{}.Fields()
	_ = roomFields
	// roomDescName is the schema descriptor for name field.
	roomDescName := roomFields[0].Descriptor()
	// room.NameValidator is a validator for the "name" field. It is called by the builders before save.
	room.NameValidator = roomDescName.Validators[0].(func(string) error)
}
