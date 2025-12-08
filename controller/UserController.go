package controller

import (
	"net/http"
	"strconv"
	"strings"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/impl"
	"tiny-admin-api-serve/middleware"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService impl.UserImpl
}

func NewUserController() *UserController {
	return &UserController{
		userService: impl.User,
	}
}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		// 从JWT获取当前用户邮箱
		claims, _ := c.Get("claims")
		userClaims := claims.(*middleware.UserClaims)
		email = userClaims.Email
	}
	// 查询用户信息
	var user dto.User
	if err := uc.userService.FindByEmail(email, &user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, user)
}

// Register 用户注册
func (uc *UserController) Register(c *gin.Context) {
	var createUserDto dto.CreateUserDto
	if err := c.ShouldBindJSON(&createUserDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userVo, err := uc.userService.CreateUser(createUserDto, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userVo)
}

// DelUser 删除用户
func (uc *UserController) DelUser(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	userVo, err := uc.userService.RemoveUserInfo(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userVo)
}

// UpdateUser 更新用户信息
func (uc *UserController) UpdateUser(c *gin.Context) {
	var updateUserDto dto.UpdateUserDto
	if err := c.ShouldBindJSON(&updateUserDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userVo, err := uc.userService.UpdateUserInfo(updateUserDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userVo)
}

// GetAllUser 获取所有用户（分页）
func (uc *UserController) GetAllUser(c *gin.Context) {
	// 获取查询参数
	var paginationQuery = dto.NewPaginationQueryDto()
	paginationQuery.Page = 1
	paginationQuery.Limit = 10

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			paginationQuery.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			paginationQuery.Limit = limit
		}
	}

	name := c.Query("name")
	email := c.Query("email")

	// 处理role数组参数
	var roles []int
	roleParams := c.QueryArray("role")
	for _, roleStr := range roleParams {
		if role, err := strconv.Atoi(roleStr); err == nil {
			roles = append(roles, role)
		}
	}

	users, err := uc.userService.GetAllUser(paginationQuery, name, email, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdatePwdAdmin 管理员强制更新密码
func (uc *UserController) UpdatePwdAdmin(c *gin.Context) {
	var updatePwdAdminDto dto.UpdatePwdAdminDto
	if err := c.ShouldBindJSON(&updatePwdAdminDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.userService.UpdatePwdAdmin(updatePwdAdminDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "ok")
}

// UpdatePwdUser 用户更新自己的密码
func (uc *UserController) UpdatePwdUser(c *gin.Context) {
	var updatePwdUserDto dto.UpdatePwdUserDto
	if err := c.ShouldBindJSON(&updatePwdUserDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.userService.UpdatePwdUser(updatePwdUserDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// 调用Logout方法将token加入黑名单
	middleware.Auth.Logout(tokenString)
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

// BatchRemoveUser 批量删除用户
func (uc *UserController) BatchRemoveUser(c *gin.Context) {
	var emails []string
	if err := c.ShouldBindJSON(&emails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userVos, err := uc.userService.BatchDeleteUser(emails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userVos)
}
