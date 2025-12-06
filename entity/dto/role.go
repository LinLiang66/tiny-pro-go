package dto

type CreateRoleDto struct {
	Name          string  `json:"name" binding:"required"`
	PermissionIds []int64 `json:"permissionIds" binding:"required"`
	MenuIds       []int64 `json:"menuIds" binding:"required"`
}

type Role struct {
	ID   int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:name"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "role"
}

type UpdateRoleDto struct {
	ID            int     `json:"id" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	PermissionIds []int64 `json:"permissionIds" binding:"required"`
	MenuIds       []int64 `json:"menuIds" binding:"required"`
}
type RolePMVo struct {
	RoleInfo *PageWrapper[Role] `json:"roleInfo"`
	MenuTree []MenuVo           `json:"menuTree"`
}
