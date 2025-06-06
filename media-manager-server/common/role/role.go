package role

import (
	"fmt"
	"strings"
)

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleViewer Role = "viewer"
	RoleEditor Role = "editor"
)

var AllRoles = []Role{
	RoleAdmin,
	RoleViewer,
	RoleEditor,
}

func IsValid(role string) bool {
	switch Role(role) {
	case RoleAdmin, RoleViewer, RoleEditor:
		return true
	default:
		return false
	}
}

func AllowedRolesString() string {
	var roles []string
	for _, r := range AllRoles {
		roles = append(roles, fmt.Sprintf("'%s'", r))
	}
	return strings.Join(roles, ", ")
}
