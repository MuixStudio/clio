package user

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/middleware"
)

func Register(r *gin.RouterGroup) {

	ar := r.Group("")
	ar.Use(middleware.TestMid())

	r.GET("/users", NewUserHandler().GetUsers())

	r.GET("/user/user_info", NewUserHandler().GetUserInfo())
	r.GET("/user/profile", NewUserHandler().GetUserProfile())
	r.GET("/user/email", NewUserHandler().GetUserEmail())
	r.GET("/user/phone", NewUserHandler().GetUserPhone())
	r.GET("/user/status", NewUserHandler().GetUserStatus())

	r.GET("/user/:user_id", NewUserHandler().GetUserByUserID())
	r.GET("/user/:user_id/profile", NewUserHandler().GetUserProfileByUserID())
	r.GET("/user/:user_id/email", NewUserHandler().GetUserEmailByUserID())
	r.GET("/user/:user_id/phone", NewUserHandler().GetUserPhoneByUserID())
	r.GET("/user/:user_id/status", NewUserHandler().GetUserStatusByUserID())

	r.POST("/user", NewUserHandler().CreateUser())

	r.PUT("/user/profile", NewUserHandler().ChangeProfile())
	r.PUT("/user/password", NewUserHandler().ChangePassword())
	r.PUT("/user/username", NewUserHandler().ChangeUserName())
	r.PUT("/user/name", NewUserHandler().ChangeName())
	r.PUT("/user/email", NewUserHandler().ChangeEmail())
	r.PUT("/user/phone", NewUserHandler().ChangePhone())
	r.PUT("/user/status", NewUserHandler().ChangeStatus())
	r.PUT("/user/identity", NewUserHandler().ChangeIdentity())

	ar.PUT("/user/:user_id/profile", NewUserHandler().ChangeProfileByUserID())
	ar.PUT("/user/:user_id/password", NewUserHandler().ChangePasswordByUserID())
	ar.PUT("/user/:user_id/username", NewUserHandler().ChangeUserNameByUserID())
	ar.PUT("/user/:user_id/name", NewUserHandler().ChangeNameByUserID())
	ar.PUT("/user/:user_id/email", NewUserHandler().ChangeEmailByUserID())
	ar.PUT("/user/:user_id/phone", NewUserHandler().ChangePhoneByUserID())
	ar.PUT("/user/:user_id/status", NewUserHandler().ChangeStatusByUserID())
	ar.PUT("/user/:user_id/identity", NewUserHandler().ChangeIdentityByUserID())

	r.DELETE("/user", NewUserHandler().DeleteUser())
	r.DELETE("/user/:user_id", NewUserHandler().DeleteUserByUserID())
}
