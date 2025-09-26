package user

import "github.com/gin-gonic/gin"

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (ah UserHandler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) GetUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) GetUserEmail() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) GetUserPhone() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) GetUserStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) GetUserByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) GetUserProfileByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) GetUserEmailByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) GetUserPhoneByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) GetUserStatusByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) ChangeProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) ChangeUserName() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) ChangeName() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) ChangeEmail() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) ChangeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) ChangeIdentity() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) ChangePhone() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// ChangePassword 将修改自身的密码
func (ah UserHandler) ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// ChangePasswordByUserID 将修改指定用户的密码，该接口只允许管理员调用
func (ah UserHandler) ChangePasswordByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) ChangeProfileByUserID() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func (ah UserHandler) ChangeUserNameByUserID() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func (ah UserHandler) ChangeNameByUserID() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func (ah UserHandler) ChangeEmailByUserID() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func (ah UserHandler) ChangePhoneByUserID() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func (ah UserHandler) ChangeStatusByUserID() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func (ah UserHandler) ChangeIdentityByUserID() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func (ah UserHandler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah UserHandler) DeleteUserByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
