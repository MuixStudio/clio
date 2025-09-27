package auth

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/svc"
	"github.com/muixstudio/clio/internal/user/pb/user"
)

type AuthHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAuthHandler(svcCtx *svc.ServiceContext) *AuthHandler {
	return &AuthHandler{
		svcCtx: svcCtx,
	}
}

func (ah AuthHandler) Login() gin.HandlerFunc {

	return func(c *gin.Context) {
		// 解析参数
		var req LoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		//业务逻辑
		resp, err := ah.loginLogic(c, req)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		}
		//返回
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": resp,
		})
	}
}

func (ah AuthHandler) loginLogic(c context.Context, req LoginReq) (LoginResp, error) {

	_, err := ah.svcCtx.UserService.VerifyPassword(c, &user.VerifyPasswordRequest{
		UserName: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return LoginResp{}, err
	}
	return LoginResp{}, nil
}

func (ah AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah AuthHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (ah AuthHandler) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
