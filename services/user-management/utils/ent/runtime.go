// Code generated by ent, DO NOT EDIT.

package ent

import (
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/ent/role"
	"github.com/Ali-Gorgani/chat-room-project/services/user-management/utils/ent/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	roleFields := schema.Role{}.Fields()
	_ = roleFields
	// roleDescName is the schema descriptor for name field.
	roleDescName := roleFields[0].Descriptor()
	// role.NameValidator is a validator for the "name" field. It is called by the builders before save.
	role.NameValidator = roleDescName.Validators[0].(func(string) error)
	// roleDescPermissions is the schema descriptor for permissions field.
	roleDescPermissions := roleFields[1].Descriptor()
	// role.DefaultPermissions holds the default value on creation for the permissions field.
	role.DefaultPermissions = roleDescPermissions.Default.([]string)
}
