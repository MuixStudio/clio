package handler

import (
	"context"
	"errors"

	"github.com/muixstudio/clio/internal/common/pb/userService"
	"github.com/muixstudio/clio/internal/user/models/dao"
	"github.com/muixstudio/clio/internal/user/svc"
)

// UserImpl implements the last service interface defined in the IDL.
type UserImpl struct {
	svcCtx *svc.ServiceContext
}

func NewUserImpl(svcCtx *svc.ServiceContext) *UserImpl {
	return &UserImpl{
		svcCtx: svcCtx,
	}
}

// CreateUser implements the UserImpl interface.
func (s *UserImpl) CreateUser(ctx context.Context, req *userService.CreateUserRequest) (resp *userService.CreateUserResponse, err error) {
	// TODO: Your code here...
	return
}

// CreateUsers implements the UserImpl interface.
func (s *UserImpl) CreateUsers(ctx context.Context, req *userService.CreateUsersRequest) (resp *userService.CreateUsersResponse, err error) {
	// TODO: Your code here...
	return
}

// FindUserByID implements the UserImpl interface.
func (s *UserImpl) FindUserByID(ctx context.Context, req *userService.FindUserByIDRequest) (resp *userService.FindUserByIDResponse, err error) {
	// TODO: Your code here...
	return
}

// FindUsersByIDs implements the UserImpl interface.
func (s *UserImpl) FindUsersByIDs(ctx context.Context, req *userService.FindUsersByIDsRequest) (resp *userService.FindUsersByIDsResponse, err error) {
	// TODO: Your code here...
	return
}

// FindUsersByName implements the UserImpl interface.
func (s *UserImpl) FindUsersByName(ctx context.Context, req *userService.FindUsersByNameRequest) (resp *userService.FindUsersByNameResponse, err error) {
	// TODO: Your code here...
	return
}

// FindUsers implements the UserImpl interface.
func (s *UserImpl) FindUsers(ctx context.Context, req *userService.FindUsersRequest) (resp *userService.FindUsersResponse, err error) {
	// TODO: Your code here...
	return
}

// VerifyPassword implements the UserImpl interface.
func (s *UserImpl) VerifyPassword(ctx context.Context, req *userService.VerifyPasswordRequest) (resp *userService.VerifyPasswordResponse, err error) {
	us, err := s.svcCtx.UserModel.Find(ctx, &dao.User{
		UserName: &req.UserName,
		Password: &req.Password,
	}, -1, 1)
	if err != nil {
		return &userService.VerifyPasswordResponse{}, err
	}
	if len(us) == 0 {
		return &userService.VerifyPasswordResponse{}, errors.New("user not found")
	}

	return &userService.VerifyPasswordResponse{
		UserID: us[0].ID,
	}, nil
}

// DeleteUser implements the UserImpl interface.
func (s *UserImpl) DeleteUser(ctx context.Context, req *userService.DeleteUserRequest) (resp *userService.DeleteUserResponse, err error) {
	// TODO: Your code here...
	return
}

// DeleteUsers implements the UserImpl interface.
func (s *UserImpl) DeleteUsers(ctx context.Context, req *userService.DeleteUsersRequest) (resp *userService.DeleteUsersResponse, err error) {
	// TODO: Your code here...
	return
}

// UpdateEmail implements the UserImpl interface.
func (s *UserImpl) UpdateEmail(ctx context.Context, req *userService.UpdateUserEmailRequest) (resp *userService.UpdateUserEmailResponse, err error) {
	// TODO: Your code here...
	return
}

// UpdateName implements the UserImpl interface.
func (s *UserImpl) UpdateName(ctx context.Context, req *userService.UpdateNameRequest) (resp *userService.UpdateNameResponse, err error) {
	// TODO: Your code here...
	return
}

// UpdateUserName implements the UserImpl interface.
func (s *UserImpl) UpdateUserName(ctx context.Context, req *userService.UpdateUserNameRequest) (resp *userService.UpdateUserNameResponse, err error) {
	// TODO: Your code here...
	return
}

// UpdatePhone implements the UserImpl interface.
func (s *UserImpl) UpdatePhone(ctx context.Context, req *userService.UpdateUserPhoneRequest) (resp *userService.UpdateUserPhoneResponse, err error) {
	// TODO: Your code here...
	return
}

// ChangeAdminStatus implements the UserImpl interface.
func (s *UserImpl) ChangeAdminStatus(ctx context.Context, req *userService.ChangeAdminStatusRequest) (resp *userService.ChangeAdminStatusResponse, err error) {
	// TODO: Your code here...
	return
}
