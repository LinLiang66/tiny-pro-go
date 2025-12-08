// impl/roleImpl.go
package impl

import (
	"errors"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/utils"
)

type RoleImpl struct {
	BaseImpl
}

var Role = RoleImpl{}

// CreateRole 创建角色
func (r RoleImpl) CreateRole(createRoleDto dto.CreateRoleDto, isInit bool) (*dto.Role, error) {
	// 检查角色是否已存在
	var existingRole dto.Role
	err := utils.Db.DB.Where("name = ?", createRoleDto.Name).First(&existingRole).Error

	if isInit && err == nil {
		return &existingRole, nil
	}

	if err == nil {
		return nil, errors.New("role already exists")
	}

	// 创建新角色
	newRole := dto.Role{
		Name: createRoleDto.Name,
	}

	result := utils.Db.DB.Create(&newRole)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newRole, nil
}

// FindAllDetail 获取角色详细信息
func (r RoleImpl) FindAllDetail(page, limit int, name string) (*dto.RolePMVo, error) {
	// 处理分页参数
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	// 构建查询
	query := utils.Db.DB.Model(&dto.Role{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	var roles []dto.Role
	var total int64

	// 获取总数
	query.Count(&total)

	// 分页查询
	result := query.Offset(offset).Limit(limit).Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	menuList, _ := Menu.FindAllMenu()
	// 计算分页元数据
	totalPages := int((total + int64(limit) - 1) / int64(limit)) // 向上取整计算总页数

	// 构建返回结果
	rolePMVo := &dto.RolePMVo{
		RoleInfo: dto.NewPageWrapper[dto.Role](
			roles,
			total,
			len(roles),
			limit,
			totalPages,
			page,
		),
		MenuTree: menuList,
	}

	return rolePMVo, nil
}

// UpdateRole 更新角色
func (r RoleImpl) UpdateRole(updateRoleDto dto.UpdateRoleDto) (*dto.Role, error) {
	var role dto.Role
	err := utils.Db.DB.Where("id = ?", updateRoleDto.ID).First(&role).Error
	if err != nil {
		return nil, errors.New("role not found")
	}

	if updateRoleDto.Name != "" {
		role.Name = updateRoleDto.Name
	}

	result := utils.Db.DB.Save(&role)
	if result.Error != nil {
		return nil, result.Error
	}

	return &role, nil
}

// RemoveRoleById 根据ID删除角色
func (r RoleImpl) RemoveRoleById(id int) (map[string]string, error) {
	var role dto.Role
	err := utils.Db.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, errors.New("role not found")
	}

	// 检查是否有用户关联该角色
	var userCount int64
	utils.Db.DB.Model(&dto.User{}).Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", id).Count(&userCount)

	if userCount > 0 {
		return nil, errors.New("role is associated with users, cannot delete")
	}

	// 删除角色
	result := utils.Db.DB.Delete(&role)
	if result.Error != nil {
		return nil, result.Error
	}

	resultMap := map[string]string{
		"name": role.Name,
	}

	return resultMap, nil
}

// FindAllRole 获取所有角色
func (r RoleImpl) FindAllRole() ([]dto.Role, error) {
	var roles []dto.Role
	result := utils.Db.DB.Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}

	var roleSimpleVos []dto.Role
	for _, role := range roles {
		roleSimpleVos = append(roleSimpleVos, dto.Role{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	return roleSimpleVos, nil
}

// FindOne 根据ID查找角色，包含关联的权限和菜单信息
func (r RoleImpl) FindOne(id int) (*dto.Role, error) {
	var role dto.Role
	err := utils.Db.DB.Where("id = ?", id).Preload("Permissions").Preload("Menus").First(&role).Error
	if err != nil {
		return nil, errors.New("role not found")
	}

	return &role, nil
}
