package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muixstudio/clio/internal/aggregater/svc"
	"github.com/muixstudio/clio/internal/aggregater/utils"
	"github.com/muixstudio/clio/internal/user/pb/user"
)

var (
	accessTokenSecret  = []byte("access_token")
	refreshTokenSecret = []byte("refresh_token")
	accessTokenExp     = time.Hour * 24     // 1 day
	refreshTokenExp    = time.Hour * 24 * 7 // 7 days
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
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    10001,
				"message": err.Error(),
			})
			return
		}
		//业务逻辑
		resp, err := ah.loginLogic(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    10001,
				"message": err.Error(),
			})
			return
		}
		accessTokenStr, err := utils.GenerateAccessToken(resp.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    10001,
				"message": err.Error(),
			})
		}
		refreshTokenStr, err := utils.GenerateRefreshToken(resp.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    10001,
				"message": err.Error(),
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "ok",
			"data": gin.H{
				"access_token":  accessTokenStr,
				"refresh_token": refreshTokenStr,
			},
		})
	}
}

func (ah AuthHandler) loginLogic(c context.Context, req LoginReq) (LoginResp, error) {
	resp, err := ah.svcCtx.UserService.VerifyPassword(c, &user.VerifyPasswordRequest{
		UserName: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return LoginResp{}, err
	}
	key := fmt.Sprintf("/auth/%d/sessions", resp.UserID)
	ah.svcCtx.RDB.Set(c, key, resp.UserID, time.Hour)
	return LoginResp{
		UserID: resp.UserID,
	}, nil
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
