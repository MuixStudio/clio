package handler

import (
	"context"

	"github.com/muixstudio/clio/internal/user/pb/user"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u UserHandler) CreateUser(ctx context.Context, request *user.CreateUserRequest, response *user.CreateUserResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) FindSingleUserByUserID(ctx context.Context, request *user.FindSingleUserByUserIDRequest, response *user.FindSingleUserByUserIDResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) FindSingleUserByUsernameAndPassword(ctx context.Context, request *user.FindSingleUserByUsernameAndPasswordRequest, response *user.FindSingleUserByUsernameAndPasswordResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) FindAllUser(ctx context.Context, request *user.FindAllUserRequest, response *user.FindAllUserResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) DeleteUser(ctx context.Context, request *user.DeleteUserRequest, response *user.DeleteUserResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdateUserInfo(ctx context.Context, request *user.UpdateUserInfoRequest, response *user.UpdateUserInfoResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdateUserEmail(ctx context.Context, request *user.UpdateUserEmailRequest, response *user.UpdateUserEmailResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdateUserName(ctx context.Context, request *user.UpdateUserNameRequest, response *user.UpdateUserNameResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdateUserPhone(ctx context.Context, request *user.UpdateUserPhoneRequest, response *user.UpdateUserPhoneResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdateUserEmployeeNumber(ctx context.Context, request *user.UpdateUserEmployeeNumberRequest, response *user.UpdateUserEmployeeNumberResponse) error {
	//TODO implement me
	panic("implement me")
}
