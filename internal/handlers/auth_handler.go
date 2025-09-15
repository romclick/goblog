package handlers

import (
	"goblog/internal/models"
	"goblog/internal/services"
	"goblog/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSerive *services.AuthService
}

func NewAuthHandler(authSerive *services.AuthService) *AuthHandler {
	return &AuthHandler{authSerive: authSerive}
}

// 新建用户（用户注册）
func (handler *AuthHandler) Register(c *gin.Context) {
	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	var req RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusUnauthorized, 422, "注册参数错误")
		return
	}
	user, err := handler.authSerive.Register(req.Username, req.Password, req.Email)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, 400, err.Error())
		return
	}
	response.Success(c, user)
}

// 用户登录
func (handler *AuthHandler) Login(c *gin.Context) {
	type LoginRequest struct {
		LoginID  string `json:"login_id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var req LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, http.StatusUnauthorized, 400, "登录参数错误")
		return
	}

	token, user, err := handler.authSerive.Login(req.LoginID, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, 400, err.Error())
		return
	}
	type LoginResponse struct {
		Token string       `json:"token"`
		User  *models.User `json:"user"`
	}
	response.Success(c, LoginResponse{
		Token: token,
		User:  user,
	})
}
