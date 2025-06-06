package domain

import (
	"mms/common/role"
)

type User struct {
	ID             string
	Username       string
	AgeRestriction int
	PinHash        *string
	Role           role.Role // custom enum

}

func (u User) IsAdmin() bool {
	return u.Role == role.RoleAdmin
}
