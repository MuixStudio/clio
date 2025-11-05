package handler

import (
	"context"

	"github.com/muixstudio/clio/services/common/pb/userService"
)

// ProfileImpl implements the last service interface defined in the IDL.
type ProfileImpl struct{}

// FindUserProfile implements the ProfileImpl interface.
func (s *ProfileImpl) FindUserProfile(ctx context.Context, req *userService.FindUserProfileRequest) (resp *userService.FindUserProfileResponse, err error) {
	// TODO: Your code here...
	return
}
