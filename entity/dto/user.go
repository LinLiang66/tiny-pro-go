package dto

type LoginBody struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	ID                int64  `json:"id" gorm:"primaryKey;autoIncrement" form:"id"`
	Address           string `json:"address" form:"address"`
	CreateTime        string `json:"createTime" form:"create_time"`
	Department        string `json:"department" form:"department"`
	Email             string `json:"email" form:"email"`
	EmployeeType      string `json:"employeeType" form:"employee_type"`
	Name              string `json:"name" form:"name"`
	Password          string `json:"-" form:"password"`
	ProbationDuration string `json:"probationDuration" form:"probation_duration"`
	ProbationEnd      string `json:"probationEnd" form:"probation_end"`
	ProbationStart    string `json:"probationStart" form:"probation_start"`
	ProtocolEnd       string `json:"protocolEnd" form:"protocol_end"`
	ProtocolStart     string `json:"protocolStart" form:"protocol_start"`
	Salt              string `json:"-" form:"salt"`
	Status            int    `json:"status" form:"status"`
	UpdateTime        string `json:"updateTime" form:"update_time"`
	Roles             []Role `json:"role" gorm:"many2many:user_role;foreignKey:id;joinForeignKey:user_id;References:id;joinReferences:role_id"`
}

func (User) TableName() string {
	return "user"

}

type GetInfo struct {
	User        *User    `json:"user"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
type CreateUserDto struct {
	Name              string  `json:"name" binding:"required"`
	Email             string  `json:"email" binding:"required"`
	Password          string  `json:"password" binding:"required"`
	RoleIds           []int64 `json:"roleIds"`
	Department        string  `json:"department"`
	EmployeeType      string  `json:"employeeType"`
	ProbationStart    string  `json:"probationStart" binding:"required"`
	ProbationEnd      string  `json:"probationEnd" binding:"required"`
	ProbationDuration string  `json:"probationDuration"`
	ProtocolStart     string  `json:"protocolStart" binding:"required"`
	ProtocolEnd       string  `json:"protocolEnd" binding:"required"`
	Address           string  `json:"address"`
	Status            *int    `json:"status"`
}

type UpdateUserDto struct {
	OldPassword       string  `json:"oldPassword" binding:"required"`
	NewPassword       string  `json:"newPassword" binding:"required"`
	Email             string  `json:"email" binding:"required"`
	RoleIds           []int64 `json:"roleIds" binding:"required"`
	Department        string  `json:"department" binding:"required"`
	EmployeeType      string  `json:"employeeType" binding:"required"`
	ProbationStart    string  `json:"probationStart" binding:"required"`
	ProbationEnd      string  `json:"probationEnd" binding:"required"`
	ProbationDuration string  `json:"probationDuration" binding:"required"`
	ProtocolStart     string  `json:"protocolStart" binding:"required"`
	ProtocolEnd       string  `json:"protocolEnd" binding:"required"`
	Address           string  `json:"address" binding:"required"`
	Status            *int    `json:"status" binding:"required"`
	Name              string  `json:"name" binding:"required"`
}
type UpdatePwdAdminDto struct {
	Email              string `json:"email" binding:"required"`
	NewPassword        string `json:"newPassword" binding:"required"`
	ConfirmNewPassword string `json:"confirmNewPassword"`
}

type UpdatePwdUserDto struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	NewPassword string `json:"newPassword" binding:"required"`
	OldPassword string `json:"oldPassword" binding:"required"`
}
