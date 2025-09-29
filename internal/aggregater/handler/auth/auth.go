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
				"message": "bind error",
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
			return
		}
		refreshTokenStr, err := utils.GenerateRefreshToken(resp.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    10001,
				"message": err.Error(),
			})
			return
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
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    10000,
				"message": "Unauthorized",
			})
			return
		}

		// 在此检测refresh_token是否使用过，可以使用redis，保证一个refresh_token只能刷新一次token

		//--------

		refreshTokenClaims, err := utils.ParseRefreshToken(refreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    10001,
				"message": "Unauthorized",
			})
			return
		}

		userId, ok := refreshTokenClaims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    10002,
				"message": "internal server error",
			})
			return
		}

		accessTokenStr, err := utils.GenerateAccessToken(uint64(userId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    10003,
				"message": "internal server error",
			})
			return
		}
		refreshTokenStr, err := utils.GenerateRefreshToken(uint64(userId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    10004,
				"message": "internal server error",
			})
			return
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
