package impl

import (
	"errors"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/utils"
)

type PermissionImpl struct {
	BaseImpl
}

var Permission = PermissionImpl{}

// Create 创建权限
func (p PermissionImpl) Create(createPermissionDto dto.Permission, isInit bool) (*dto.Permission, error) {
	var existingPermission dto.Permission
	err := utils.Db.DB.Where("name = ?", createPermissionDto.Name).First(&existingPermission).Error

	// 情况1：初始化模式且权限已存在
	if isInit && err == nil {
		result := &dto.Permission{
			ID:   existingPermission.ID,
			Name: existingPermission.Name,
			Desc: existingPermission.Desc,
		}
		return result, nil
	}

	// 情况2：非初始化模式且权限已存在
	if err == nil && !isInit {
		return nil, errors.New("permission already exists")
	}

	// 情况3：创建新权限
	newPermission := dto.Permission{
		Name: createPermissionDto.Name,
		Desc: createPermissionDto.Desc,
	}

	result := utils.Db.DB.Create(&newPermission)
	if result.Error != nil {
		return nil, result.Error
	}

	permissionVo := &dto.Permission{
		ID:   newPermission.ID,
		Name: newPermission.Name,
		Desc: newPermission.Desc,
	}

	return permissionVo, nil
}

// UpdatePermission 更新权限
func (p PermissionImpl) UpdatePermission(updatePermissionDto dto.Permission) (*dto.Permission, error) {
	var permission dto.Permission
	err := utils.Db.DB.Where("id = ?", updatePermissionDto.ID).First(&permission).Error
	if err != nil {
		return nil, errors.New("permission not found")
	}

	permission.Name = updatePermissionDto.Name
	permission.Desc = updatePermissionDto.Desc

	result := utils.Db.DB.Save(&permission)
	if result.Error != nil {
		return nil, result.Error
	}

	permissionVo := &dto.Permission{
		ID:   permission.ID,
		Name: permission.Name,
		Desc: permission.Desc,
	}

	return permissionVo, nil
}

// FindPermissions 查询权限列表（分页）
func (p PermissionImpl) FindPermissions(page, limit int, name string) (*dto.PageWrapper[dto.Permission], error) {
	// 处理分页参数
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	// 构建查询
	query := utils.Db.DB.Model(&dto.Permission{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	var permissions []dto.Permission
	var total int64

	// 获取总数
	query.Count(&total)

	// 分页查询
	result := query.Offset(offset).Limit(limit).Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	// 计算分页元数据
	totalPages := int((total + int64(limit) - 1) / int64(limit)) // 向上取整计算总页数

	// 构造分页结果
	pageWrapper := dto.NewPageWrapper[dto.Permission](
		permissions,
		total,
		len(permissions),
		limit,
		totalPages,
		page,
	)

	return pageWrapper, nil
}

// DelPermission 删除权限
func (p PermissionImpl) DelPermission(id int) (*dto.Permission, error) {
	var permission dto.Permission
	err := utils.Db.DB.Where("id = ?", id).First(&permission).Error
	if err != nil {
		return nil, errors.New("permission not found")
	}

	// 删除角色权限关联记录
	utils.Db.DB.Exec("DELETE FROM role_permission WHERE permission_id = ?", permission.ID)

	// 删除权限
	result := utils.Db.DB.Delete(&permission)
	if result.Error != nil {
		return nil, result.Error
	}

	createPermissionDto := &dto.Permission{
		Name: permission.Name,
		Desc: permission.Desc,
	}

	return createPermissionDto, nil
}

// FindAllPermission 查询所有权限
func (p PermissionImpl) FindAllPermission() ([]dto.Permission, error) {
	var permissions []dto.Permission
	result := utils.Db.DB.Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}
