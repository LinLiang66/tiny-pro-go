package controller

import (
	"net/http"
	"strconv"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/impl"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService impl.RoleImpl
}

func NewRoleController() *RoleController {
	return &RoleController{
		roleService: impl.Role,
	}
}

// Create 创建角色
func (rc *RoleController) Create(c *gin.Context) {
	var createRoleDto dto.CreateRoleDto
	if err := c.ShouldBindJSON(&createRoleDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := rc.roleService.CreateRole(createRoleDto, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// GetAllRole 获取所有角色信息
func (rc *RoleController) GetAllRole(c *gin.Context) {
	roles, err := rc.roleService.FindAllRole()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// GetAllRoleDetail 获取角色详细信息
func (rc *RoleController) GetAllRoleDetail(c *gin.Context) {
	// 获取查询参数
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	name := c.Query("name")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	roleDetails, err := rc.roleService.FindAllDetail(page, limit, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roleDetails)
}

// UpdateRole 更新角色
func (rc *RoleController) UpdateRole(c *gin.Context) {
	var updateRoleDto dto.UpdateRoleDto
	if err := c.ShouldBindJSON(&updateRoleDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := rc.roleService.UpdateRole(updateRoleDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// DeleteRole 删除角色
func (rc *RoleController) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := rc.roleService.RemoveRoleById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetRoleInfo 获取角色详情
func (rc *RoleController) GetRoleInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	role, err := rc.roleService.FindOne(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}
