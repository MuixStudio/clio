package handler

import (
	"context"
	"errors"

	kratosErrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/muixstudio/clio/internal/common/pb/userService"
	"github.com/muixstudio/clio/internal/user/models/dao"
	"github.com/muixstudio/clio/internal/user/svc"
	"gorm.io/gorm"
)

type UserHandler struct {
	userService.UnimplementedUserServer
	svcCtx *svc.ServiceContext
}

func NewUserHandler(serviceContext *svc.ServiceContext) *UserHandler {
	return &UserHandler{
		svcCtx: serviceContext,
	}
}

func (u UserHandler) CreateUser(ctx context.Context, request *userService.CreateUserRequest) (*userService.CreateUserResponse, error) {
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
		return &userService.CreateUserResponse{}, err
	}
	return &userService.CreateUserResponse{
		Id: data.ID,
	}, nil
}

func (u UserHandler) CreateUsers(ctx context.Context, request *userService.CreateUsersRequest) (*userService.CreateUsersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) FindUserByID(ctx context.Context, request *userService.FindUserByIDRequest) (*userService.FindUserByIDResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) FindUsersByIDs(ctx context.Context, request *userService.FindUsersByIDsRequest) (*userService.FindUsersByIDsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) FindUsersByName(ctx context.Context, request *userService.FindUsersByNameRequest) (*userService.FindUsersByNameResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) FindUsers(ctx context.Context, request *userService.FindUsersRequest) (*userService.FindUsersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) VerifyPassword(ctx context.Context, request *userService.VerifyPasswordRequest) (*userService.VerifyPasswordResponse, error) {
	if request.UserName == nil {
		return &userService.VerifyPasswordResponse{}, kratosErrors.BadRequest("BAD_REQUEST", "userName is required")
	}
	user, err := u.svcCtx.UserModel.Take(ctx, &dao.User{
		UserName: request.UserName,
		Password: request.Password,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &userService.VerifyPasswordResponse{}, kratosErrors.New(404, "INVALID_USERNAME_OR_PASSWORD", "invalid username or password").WithCause(err)
		}
		return &userService.VerifyPasswordResponse{}, kratosErrors.InternalServer("INTERNAL_SERVER_ERROR", "internal server error").WithCause(err)
	}
	return &userService.VerifyPasswordResponse{
		UserID: user.ID,
	}, nil
}

func (u UserHandler) DeleteUser(ctx context.Context, request *userService.DeleteUserRequest) (*userService.DeleteUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) DeleteUsers(ctx context.Context, request *userService.DeleteUsersRequest) (*userService.DeleteUsersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdateEmail(ctx context.Context, request *userService.UpdateUserEmailRequest) (*userService.UpdateUserEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdateName(ctx context.Context, request *userService.UpdateNameRequest) (*userService.UpdateNameResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdateUserName(ctx context.Context, request *userService.UpdateUserNameRequest) (*userService.UpdateUserNameResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) UpdatePhone(ctx context.Context, request *userService.UpdateUserPhoneRequest) (*userService.UpdateUserPhoneResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) ChangeAdminStatus(ctx context.Context, request *userService.ChangeAdminStatusRequest) (*userService.ChangeAdminStatusResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserHandler) mustEmbedUnimplementedUserServer() {
	//TODO implement me
	panic("implement me")
}
