package resolvers

import (
	mediaService "mms/modules/media-module/services"
	userService "mms/modules/user-module/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService  *userService.UserService
	MediaService *mediaService.MediaService
}
