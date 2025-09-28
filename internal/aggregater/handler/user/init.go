package user

import (
	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/middleware"
	"github.com/muixstudio/clio/internal/aggregater/svc"
)

func Register(r *gin.RouterGroup, svcCtx *svc.ServiceContext) {

	userHandler := NewUserHandler(svcCtx)

	ar := r.Group("")
	ar.Use(middleware.TestMid())

	r.GET("/users", userHandler.GetUsers())

	r.GET("/user/userinfo", userHandler.GetUserInfo())
	r.GET("/user/profile", userHandler.GetUserProfile())
	r.GET("/user/email", userHandler.GetUserEmail())
	r.GET("/user/phone", userHandler.GetUserPhone())
	r.GET("/user/status", userHandler.GetUserStatus())

	r.GET("/user/:user_id", userHandler.GetUserByUserID())
	r.GET("/user/:user_id/profile", userHandler.GetUserProfileByUserID())
	r.GET("/user/:user_id/email", userHandler.GetUserEmailByUserID())
	r.GET("/user/:user_id/phone", userHandler.GetUserPhoneByUserID())
	r.GET("/user/:user_id/status", userHandler.GetUserStatusByUserID())

	r.POST("/user", userHandler.CreateUser())

	r.PUT("/user/profile", userHandler.ChangeProfile())
	r.PUT("/user/password", userHandler.ChangePassword())
	r.PUT("/user/username", userHandler.ChangeUserName())
	r.PUT("/user/name", userHandler.ChangeName())
	r.PUT("/user/email", userHandler.ChangeEmail())
	r.PUT("/user/phone", userHandler.ChangePhone())
	r.PUT("/user/status", userHandler.ChangeStatus())
	r.PUT("/user/identity", userHandler.ChangeIdentity())

	ar.PUT("/user/:user_id/profile", userHandler.ChangeProfileByUserID())
	ar.PUT("/user/:user_id/password", userHandler.ChangePasswordByUserID())
	ar.PUT("/user/:user_id/username", userHandler.ChangeUserNameByUserID())
	ar.PUT("/user/:user_id/name", userHandler.ChangeNameByUserID())
	ar.PUT("/user/:user_id/email", userHandler.ChangeEmailByUserID())
	ar.PUT("/user/:user_id/phone", userHandler.ChangePhoneByUserID())
	ar.PUT("/user/:user_id/status", userHandler.ChangeStatusByUserID())
	ar.PUT("/user/:user_id/identity", userHandler.ChangeIdentityByUserID())

	r.DELETE("/user", userHandler.DeleteUser())
	r.DELETE("/user/:user_id", userHandler.DeleteUserByUserID())
}
