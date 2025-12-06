package controller

import (
	"net/http"
	"strings"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/impl"
	"tiny-admin-api-serve/middleware"
	"tiny-admin-api-serve/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService impl.UserImpl
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: impl.User,
	}
}

func (a *AuthController) Login(c *gin.Context) {
	// 验证用户名密码后生成token
	var loginBody dto.LoginBody
	err := c.ShouldBindJSON(&loginBody)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
		return
	}
	var user dto.User
	impl.User.FindByEmail(loginBody.Email, &user)
	if &user == nil {
		c.JSON(http.StatusOK, gin.H{"msg": "用户不存在"})
		return
	}
	isValid, err := utils.VerifyPassword(loginBody.Password, user.Salt, user.Password)
	if err != nil || !isValid {
		c.JSON(http.StatusOK, gin.H{"msg": "用户名或密码错误"})
		return
	}
	token, err := middleware.Auth.GenerateToken(user.ID, user.Email, "user")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"token":   token,
		"message": "Login successful",
	})
}

func (a *AuthController) Profile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	username := c.MustGet("username").(string)
	role := c.MustGet("role").(string)

	c.JSON(200, gin.H{
		"user_id":  userID,
		"username": username,
		"role":     role,
	})
}

func (a *AuthController) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// 调用Logout方法将token加入黑名单
	middleware.Auth.Logout(tokenString)
	c.JSON(200, gin.H{"message": "Logged out successfully"})
}

func (a *AuthController) Refresh(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	oldToken := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := middleware.Auth.RefreshToken(oldToken)
	if err != nil {
		return
	}
	c.JSON(200, gin.H{"message": "Token refreshed", "token": token})
}

func (a *AuthController) Register(c *gin.Context) {
	// 注册逻辑...
	c.JSON(200, gin.H{"message": "User registered"})
}
