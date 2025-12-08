package impl

import (
	"errors"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/utils"

	"gorm.io/gorm"
)

type UserImpl struct {
	BaseImpl
}

var User = UserImpl{}

// FindByEmail 获取用户信息，包括角色及角色关联的权限和菜单
func (u UserImpl) FindByEmail(email string, user *dto.User) error {
	err := utils.Db.DB.Model(&dto.User{}).
		Where("email = ?", email).
		Preload("Roles").             // 预加载用户的角色
		Preload("Roles.Permissions"). // 预加载角色的权限
		Preload("Roles.Menus").       // 预加载角色的菜单
		First(user).
		Error

	// 明确处理用户不存在的情况
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	}

	return err
}

// CreateUser 创建用户
func (u UserImpl) CreateUser(createUserDto dto.CreateUserDto, isInit bool) (*dto.User, error) {
	// 1. 检查用户是否已存在
	var existingUser dto.User
	err := utils.Db.DB.Where("email = ?", createUserDto.Email).First(&existingUser).Error

	if isInit && err == nil {
		return &existingUser, nil
	}

	if err == nil {
		return nil, errors.New("user already exists")
	}

	// 2. 获取关联角色（这里简化处理，实际应查询角色表）
	// roles := iRoleRepository.findAllById(createUserDto.RoleIds)

	// 3. 创建并保存用户
	salt, _ := utils.GenerateSalt()
	hashedPassword, _ := utils.Encry(createUserDto.Password, salt)

	user := dto.User{
		Email:             createUserDto.Email,
		Password:          hashedPassword,
		Name:              createUserDto.Name,
		Department:        createUserDto.Department,
		EmployeeType:      createUserDto.EmployeeType,
		ProbationStart:    createUserDto.ProbationStart,
		ProbationEnd:      createUserDto.ProbationEnd,
		ProbationDuration: createUserDto.ProbationDuration,
		ProtocolStart:     createUserDto.ProtocolStart,
		ProtocolEnd:       createUserDto.ProtocolEnd,
		Address:           createUserDto.Address,
		Salt:              salt,
		Status:            *createUserDto.Status, // 默认状态
	}

	if createUserDto.Status != nil {
		user.Status = *createUserDto.Status
	}

	result := utils.Db.DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// GetRoleByUserId 根据用户ID获取角色
func (u UserImpl) GetRoleByUserId(userId int64) ([]dto.Permission, error) {
	var user dto.User
	// 这里需要联表查询用户的角色和权限，简化处理
	err := utils.Db.DB.Preload("Roles.Permissions").Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 构建权限列表
	var permissions []dto.Permission
	// 这里需要根据实际的关联关系进行处理

	return permissions, nil
}

// RemoveUserInfo 删除用户信息
func (u UserImpl) RemoveUserInfo(email string) (*dto.User, error) {
	var user dto.User
	err := utils.Db.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	result := utils.Db.DB.Delete(&user)
	if result.Error != nil {
		return nil, errors.New("failed to delete user")
	}

	return &user, nil
}

// UpdateUserInfo 更新用户信息
func (u UserImpl) UpdateUserInfo(updateUserDto dto.UpdateUserDto) (*dto.User, error) {
	var user dto.User
	err := utils.Db.DB.Where("email = ?", updateUserDto.Email).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 更新用户信息
	user.Name = updateUserDto.Name
	user.Department = updateUserDto.Department
	user.EmployeeType = updateUserDto.EmployeeType

	// 处理日期字段
	if updateUserDto.ProbationStart != "" && updateUserDto.ProbationStart != "NaN-NaN-NaN" {
		user.ProbationStart = updateUserDto.ProbationStart
	}

	if updateUserDto.ProbationEnd != "" && updateUserDto.ProbationEnd != "NaN-NaN-NaN" {
		user.ProbationEnd = updateUserDto.ProbationEnd
	}

	user.ProbationDuration = updateUserDto.ProbationDuration

	if updateUserDto.ProtocolStart != "" && updateUserDto.ProtocolStart != "NaN-NaN-NaN" {
		user.ProtocolStart = updateUserDto.ProtocolStart
	}

	if updateUserDto.ProtocolEnd != "" && updateUserDto.ProtocolEnd != "NaN-NaN-NaN" {
		user.ProtocolEnd = updateUserDto.ProtocolEnd
	}

	user.Address = updateUserDto.Address

	if updateUserDto.Status != nil {
		user.Status = *updateUserDto.Status
	}

	result := utils.Db.DB.Save(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// GetAllUser 获取所有用户（分页）
func (u UserImpl) GetAllUser(paginationQuery dto.PaginationQueryDto, name, email string, roles []int) (*dto.PageWrapper[dto.User], error) {
	var users []dto.User
	var total int64

	query := utils.Db.DB.Model(&dto.User{})

	// 添加查询条件
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if email != "" {
		query = query.Where("email LIKE ?", "%"+email+"%")
	}

	// 角色查询条件（简化处理）
	if len(roles) > 0 {
		// 这里需要根据实际的用户角色关联表进行处理
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (paginationQuery.Page - 1) * paginationQuery.Limit
	result := query.Offset(offset).Limit(paginationQuery.Limit).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	// 计算分页信息
	totalPages := 1
	if paginationQuery.Limit > 0 {
		totalPages = int((total + int64(paginationQuery.Limit) - 1) / int64(paginationQuery.Limit))
	}

	pageWrapper := dto.NewPageWrapper[dto.User](
		users,
		total,
		len(users),
		paginationQuery.Limit,
		totalPages,
		paginationQuery.Page,
	)
	return pageWrapper, nil

}

// UpdatePwdAdmin 管理员强制更新密码
func (u UserImpl) UpdatePwdAdmin(updatePwdAdminDto dto.UpdatePwdAdminDto) error {
	var user dto.User
	err := utils.Db.DB.Where("email = ?", updatePwdAdminDto.Email).First(&user).Error
	if err != nil {
		return errors.New("user not found")
	}

	// 生成新密码哈希
	hashedPassword, _ := utils.Encry(updatePwdAdminDto.NewPassword, user.Salt)

	// 更新密码
	result := utils.Db.DB.Model(&dto.User{}).Where("email = ?", updatePwdAdminDto.Email).Update("password", hashedPassword)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdatePwdUser 用户更新自己的密码
func (u UserImpl) UpdatePwdUser(updatePwdUserDto dto.UpdatePwdUserDto) error {
	var user dto.User
	err := utils.Db.DB.Where("email = ?", updatePwdUserDto.Email).First(&user).Error
	if err != nil {
		return errors.New("user not found")
	}

	// 验证旧密码
	isValid, _ := utils.VerifyPassword(updatePwdUserDto.OldPassword, user.Salt, user.Password)
	if !isValid {
		return errors.New("old password is incorrect")
	}

	// 生成新密码哈希
	hashedPassword, _ := utils.Encry(updatePwdUserDto.NewPassword, user.Salt)

	// 更新密码
	result := utils.Db.DB.Model(&dto.User{}).Where("email = ?", updatePwdUserDto.Email).Update("password", hashedPassword)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// BatchDeleteUser 批量删除用户
func (u UserImpl) BatchDeleteUser(emails []string) ([]dto.User, error) {
	var users []dto.User
	var deletedUsers []dto.User

	// 查询要删除的用户
	err := utils.Db.DB.Where("email IN ?", emails).Find(&users).Error
	if err != nil {
		return nil, err
	}

	deletedUsers = append(deletedUsers, users...)

	// 批量删除
	result := utils.Db.DB.Where("email IN ?", emails).Delete(&dto.User{})
	if result.Error != nil {
		return nil, result.Error
	}

	return deletedUsers, nil
}
