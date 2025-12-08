package dto

type MenuVo struct {
	ID         int64     `json:"id"`
	Label      string    `json:"label"`
	Order      int       `json:"order"`
	ParentId   *int64    `json:"parentId"`
	MenuType   string    `json:"menuType"`
	CustomIcon string    `json:"customIcon"`
	Component  string    `json:"component"`
	Url        string    `json:"url"`
	Locale     string    `json:"locale"`
	Children   *[]MenuVo `json:"children"`
}

type Menu struct {
	ID        int64  `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Name      string `json:"name" gorm:"column:name"`
	Order     int    `json:"order" gorm:"column:order"`
	ParentId  *int64 `json:"parentId" gorm:"column:parentId"`
	MenuType  string `json:"menuType" gorm:"column:menuType"`
	Icon      string `json:"icon" gorm:"column:icon"`
	Component string `json:"component" gorm:"column:component"`
	Path      string `json:"path" gorm:"column:path"`
	Locale    string `json:"locale" gorm:"column:locale"`
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menu"
}

type MenuTreeVo struct {
	ID         int64        `json:"id"`
	Label      string       `json:"label"`
	Children   []MenuTreeVo `json:"children"`
	Url        string       `json:"url"`
	Component  string       `json:"component"`
	CustomIcon string       `json:"customIcon"`
	MenuType   string       `json:"menuType"`
	ParentId   *int64       `json:"parentId"`
	Order      int          `json:"order"`
	Locale     string       `json:"locale"`
}

// ToDoubleArrayFormat 返回二维数组格式
func (m *MenuTreeVo) ToDoubleArrayFormat() [][]MenuTreeVo {
	mainList := m.ToSingleTree()
	emptyList := []MenuTreeVo{}

	return [][]MenuTreeVo{mainList, emptyList}
}

// ToSingleTree 将当前对象转为单节点树结构
func (m *MenuTreeVo) ToSingleTree() []MenuTreeVo {
	// 如果当前是根节点，直接返回包含自身的列表
	if m.ParentId == nil {
		return []MenuTreeVo{*m}
	}

	// 如果不是根节点，构建一个虚拟根节点
	virtualRoot := &MenuTreeVo{
		Children: []MenuTreeVo{*m},
	}

	return []MenuTreeVo{*virtualRoot}
}
