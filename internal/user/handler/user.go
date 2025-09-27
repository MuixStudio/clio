package handler

import (
	"context"
	"errors"

	"github.com/muixstudio/clio/internal/user/config"
	"github.com/muixstudio/clio/internal/user/models/dao"
	"github.com/muixstudio/clio/internal/user/pb/user"
	"github.com/muixstudio/clio/internal/user/svc"
	"go-micro.dev/v5/logger"
)

type UserHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUserHandler(c config.Config) *UserHandler {
	return &UserHandler{
		svcCtx: svc.NewServiceContext(c),
	}
}

func (u UserHandler) CreateUser(ctx context.Context, request *user.CreateUserRequest, response *user.CreateUserResponse) error {
	l := dao.LOCAL
	data := dao.User{
		Name:         request.Name,
		UserName:     request.UserName,
		Password:     request.Password,
		CountryCode:  request.CountryCode,
		Phone:        request.Phone,
		Email:        request.Email,
		AuthProvider: &l,
	}
	err := u.svcCtx.UserModel.Create(ctx, &data)
	if err != nil {
		logger.Error(ctx, err)
		return err
	}
	response.Id = data.ID
	return nil
}

func (u UserHandler) CreateUsers(ctx context.Context, request *user.CreateUsersRequest, response *user.CreateUsersResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) FindUserByID(ctx context.Context, request *user.FindUserByIDRequest, response *user.FindUserByIDResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) FindUsersByIDs(ctx context.Context, request *user.FindUsersByIDsRequest, response *user.FindUsersByIDsResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) FindUsersByName(ctx context.Context, request *user.FindUsersByNameRequest, response *user.FindUsersByNameResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) FindUsers(ctx context.Context, request *user.FindUsersRequest, response *user.FindUsersResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) VerifyPassword(ctx context.Context, request *user.VerifyPasswordRequest, response *user.VerifyPasswordResponse) error {
	username := request.UserName
	password := request.Password
	us, err := u.svcCtx.UserModel.Find(ctx, &dao.User{
		UserName: &username,
		Password: &password,
	}, -1, 1)
	if err != nil {
		return err
	}
	if len(us) == 0 {
		return errors.New("user not found")
	}
	response.UserID = us[0].ID
	return nil
}

func (u UserHandler) DeleteUser(ctx context.Context, request *user.DeleteUserRequest, response *user.DeleteUserResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) DeleteUsers(ctx context.Context, request *user.DeleteUsersRequest, response *user.DeleteUsersResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdateEmail(ctx context.Context, request *user.UpdateUserEmailRequest, response *user.UpdateUserEmailResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdateName(ctx context.Context, request *user.UpdateNameRequest, response *user.UpdateNameResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdateUserName(ctx context.Context, request *user.UpdateUserNameRequest, response *user.UpdateUserNameResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdatePhone(ctx context.Context, request *user.UpdateUserPhoneRequest, response *user.UpdateUserPhoneResponse) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) ChangeAdminStatus(ctx context.Context, request *user.ChangeAdminStatusRequest, response *user.ChangeAdminStatusResponse) error {
	//TODO implement me
	panic("implement me")
}
