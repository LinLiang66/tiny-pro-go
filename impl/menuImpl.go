package impl

import (
	"errors"
	"sort"
	"tiny-admin-api-serve/entity/dto"
	"tiny-admin-api-serve/utils"
)

type MenuImpl struct {
	BaseImpl
}

var Menu = MenuImpl{}

// GetMenubyEmail 根据邮箱获取菜单
func (m MenuImpl) GetMenubyEmail(email string) ([]dto.MenuVo, error) {
	// 1. 通过email获取用户，并预加载角色信息
	var user dto.User
	err := utils.Db.DB.Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. 获取用户的所有角色ID
	var roles = user.Roles
	if len(roles) == 0 {
		return []dto.MenuVo{}, nil
	}
	roleIds := make([]int64, len(roles))
	for i, role := range roles {
		roleIds[i] = role.ID
	}

	// 3. 获取这些角色关联的所有菜单
	menus, err := m.FindMenusByRoleIds(roleIds)
	if err != nil {
		return nil, err
	}

	if len(menus) == 0 {
		return []dto.MenuVo{}, nil
	}

	// 4. 构建菜单树
	menuTree := m.buildMenuTree(menus)
	return menuTree, nil
}

// FindAllMenu 获取所有菜单
func (m MenuImpl) FindAllMenu() ([]dto.MenuVo, error) {
	// 1. 查询所有菜单并按order排序
	var menus []dto.Menu
	result := utils.Db.DB.Order("`order` ASC").Find(&menus)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(menus) == 0 {
		return []dto.MenuVo{}, nil
	}

	// 2. 构建菜单树
	menuTree := m.buildMenuTree(menus)
	return menuTree, nil
}

// CreateMenu 创建菜单
func (m MenuImpl) CreateMenu(createMenuDto dto.Menu, isInit bool) (*dto.Menu, error) {
	// 检查菜单是否已存在 (简化处理)
	var existingMenu dto.Menu
	err := utils.Db.DB.Where("name = ? AND order = ? AND menu_type = ? AND parent_id = ? AND path = ? AND icon = ? AND component = ? AND locale = ?",
		createMenuDto.Name, createMenuDto.Order, createMenuDto.MenuType, createMenuDto.ParentId,
		createMenuDto.Path, createMenuDto.Icon, createMenuDto.Component, createMenuDto.Locale).First(&existingMenu).Error

	if isInit && err == nil {
		return &existingMenu, nil
	}

	if err == nil && !isInit {
		return nil, errors.New("menu already exists")
	}

	// 创建新菜单
	newMenu := dto.Menu{
		Name:      createMenuDto.Name,
		Path:      createMenuDto.Path,
		Component: createMenuDto.Component,
		ParentId:  createMenuDto.ParentId,
		MenuType:  createMenuDto.MenuType,
		Icon:      createMenuDto.Icon,
		Order:     createMenuDto.Order,
		Locale:    createMenuDto.Locale,
	}

	result := utils.Db.DB.Create(&newMenu)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newMenu, nil
}

// UpdateMenu 更新菜单
func (m MenuImpl) UpdateMenu(updateMenuDto dto.Menu) (bool, error) {
	var menu dto.Menu
	err := utils.Db.DB.Where("id = ?", updateMenuDto.ID).First(&menu).Error
	if err != nil {
		return false, errors.New("menu not found")
	}

	menu.Name = updateMenuDto.Name
	menu.Path = updateMenuDto.Path
	menu.Component = updateMenuDto.Component
	menu.ParentId = updateMenuDto.ParentId
	menu.MenuType = updateMenuDto.MenuType
	menu.Icon = updateMenuDto.Icon
	menu.Order = updateMenuDto.Order
	menu.Locale = updateMenuDto.Locale

	result := utils.Db.DB.Save(&menu)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

// DeleteMenu 删除菜单
func (m MenuImpl) DeleteMenu(id int64, parentId int64) (*dto.Menu, error) {
	// 查找要删除的菜单
	var menu dto.Menu
	err := utils.Db.DB.Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, errors.New("menu not found")
	}

	// 查找所有子菜单
	var childMenus []dto.Menu
	utils.Db.DB.Where("parent_id = ?", id).Find(&childMenus)

	// 更新子菜单的parentId
	for _, childMenu := range childMenus {
		if parentId == -1 {
			childMenu.ParentId = nil
		} else {
			childMenu.ParentId = &parentId
		}
		utils.Db.DB.Save(&childMenu)
	}

	// 删除菜单
	result := utils.Db.DB.Delete(&menu)
	if result.Error != nil {
		return nil, result.Error
	}

	return &menu, nil
}

// buildMenuTree 构建菜单树形结构
func (m MenuImpl) buildMenuTree(menus []dto.Menu) []dto.MenuVo {
	// 使用Map存储所有菜单，便于快速查找
	menuMap := make(map[int64]*dto.MenuVo)

	// 先转换所有菜单为VO对象
	for _, menu := range menus {
		menuVo := m.convertToVo(menu)
		menuMap[menu.ID] = &menuVo
	}

	// 构建树形结构
	var rootMenus []dto.MenuVo
	for _, menuVo := range menuMap {
		parentId := menuVo.ParentId
		if parentId == nil {
			rootMenus = append(rootMenus, *menuVo)
		} else {
			if parentMenu, exists := menuMap[*parentId]; exists {
				if parentMenu.Children == nil {
					children := make([]dto.MenuVo, 0)
					parentMenu.Children = &children
				}
				*parentMenu.Children = append(*parentMenu.Children, *menuVo)
			} else {
				rootMenus = append(rootMenus, *menuVo)
			}
		}
	}

	// 对菜单进行排序
	sort.Slice(rootMenus, func(i, j int) bool {
		return rootMenus[i].Order < rootMenus[j].Order
	})

	for _, menu := range rootMenus {
		if menu.Children != nil {
			sort.Slice(*menu.Children, func(i, j int) bool {
				return (*menu.Children)[i].Order < (*menu.Children)[j].Order
			})
		}
	}

	return rootMenus
}

// convertToVo 将Menu实体转换为MenuVo
func (m MenuImpl) convertToVo(menu dto.Menu) dto.MenuVo {
	menuVo := dto.MenuVo{
		ID:         menu.ID,
		Label:      menu.Name,
		ParentId:   menu.ParentId,
		Order:      menu.Order,
		Url:        menu.Path,
		Component:  menu.Component,
		CustomIcon: menu.Icon,
		Locale:     menu.Locale,
		MenuType:   menu.MenuType,
	}

	children := make([]dto.MenuVo, 0)
	menuVo.Children = &children

	return menuVo
}

// FindMenusByRoleIds 根据角色IDs查询关联的菜单
func (m MenuImpl) FindMenusByRoleIds(roleIds []int64) ([]dto.Menu, error) {
	var menus []dto.Menu
	result := utils.Db.DB.
		Distinct().
		Joins("JOIN role_menu ON menu.id = role_menu.menu_id").
		Where("role_menu.role_id IN ?", roleIds).
		Find(&menus)

	if result.Error != nil {
		return nil, result.Error
	}

	return menus, nil
}
