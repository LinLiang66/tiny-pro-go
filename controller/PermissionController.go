package controller

import (
	"net/http"
	"strconv"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/impl"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	permissionService impl.PermissionImpl
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		permissionService: impl.Permission,
	}
}

// Create 创建权限
func (pc *PermissionController) Create(c *gin.Context) {
	var createPermissionDto dto.Permission
	if err := c.ShouldBindJSON(&createPermissionDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissionVo, err := pc.permissionService.Create(createPermissionDto, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissionVo)
}

// Update 更新权限
func (pc *PermissionController) Update(c *gin.Context) {
	var updatePermissionDto dto.Permission
	if err := c.ShouldBindJSON(&updatePermissionDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permissionVo, err := pc.permissionService.UpdatePermission(updatePermissionDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissionVo)
}

// GetAll 查询权限列表
func (pc *PermissionController) GetAll(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	name := c.Query("name")

	// 如果没有分页参数，返回所有权限
	if pageStr == "" && limitStr == "" {
		permissions, err := pc.permissionService.FindAllPermission()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, permissions)
		return
	}

	// 处理分页参数
	page := 1
	limit := 10
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	result, err := pc.permissionService.FindPermissions(page, limit, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Delete 删除权限
func (pc *PermissionController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	createPermissionDto, err := pc.permissionService.DelPermission(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createPermissionDto)
}
