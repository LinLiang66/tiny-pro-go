package elastic

import (
	"context"
	"fmt"
	"log"
)

// User 示例用户实体
type User struct {
	BaseModel
	Name   string   `json:"name" es:"name"`
	Age    int      `json:"age" es:"age"`
	Email  string   `json:"email" es:"email"`
	Salary float64  `json:"salary" es:"salary"`
	Active bool     `json:"active" es:"active"`
	Tags   []string `json:"tags" es:"tags"`
}

// Example 使用示例
func Example() {
	ctx := context.Background()

	// 1. 初始化 ES 客户端
	if err := InitClient(); err != nil {
		log.Fatalf("Failed to initialize ES client: %v", err)
	}

	// 2. 创建仓库实例
	userRepo, err := NewBaseRepository[User]()
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	// 3. 插入文档 - 使用字符串ID作为文档ID
	user1 := &User{
		BaseModel: BaseModel{
			ID: "12345678901234567", // 雪花ID转换为字符串
		},
		Name:   "张三",
		Age:    25,
		Email:  "zhangsan@example.com",
		Salary: 5000.50,
		Active: true,
		Tags:   []string{"developer", "go"},
	}

	// 3.1 插入文档 - 使用数字ID（如雪花ID）作为文档ID
	user2 := &User{
		BaseModel: BaseModel{
			ID: "1234567890123456789", // 雪花ID转换为字符串
		},
		Name:   "赵六",
		Age:    40,
		Email:  "zhaoliu@example.com",
		Salary: 8000.00,
		Active: true,
		Tags:   []string{"manager", "java"},
	}

	id, err := userRepo.Insert(ctx, user1)
	if err != nil {
		log.Printf("Failed to insert user1: %v", err)
	} else {
		fmt.Printf("Inserted user1 with ID: %s\n", id)
	}

	// 插入使用雪花ID的用户
	id2, err := userRepo.Insert(ctx, user2)
	if err != nil {
		log.Printf("Failed to insert user2: %v", err)
	} else {
		fmt.Printf("Inserted user2 with ID: %s\n", id2)
	}

	// 4. 批量插入文档
	users := []*User{
		{
			Name:   "李四",
			Age:    30,
			Email:  "lisi@example.com",
			Salary: 6000.00,
			Active: true,
			Tags:   []string{"developer", "java"},
		},
		{
			Name:   "王五",
			Age:    35,
			Email:  "wangwu@example.com",
			Salary: 7000.00,
			Active: false,
			Tags:   []string{"designer"},
		},
	}

	ids, err := userRepo.InsertBatch(ctx, users)
	if err != nil {
		log.Printf("Failed to insert batch users: %v", err)
	} else {
		fmt.Printf("Inserted batch users with IDs: %v\n", ids)
	}

	// 5. 根据 ID 获取文档
	user, err := userRepo.GetById(ctx, id)
	if err != nil {
		log.Printf("Failed to get user by ID: %v", err)
	} else if user != nil {
		fmt.Printf("Get user by ID: %v\n", user.Name)
	}

	// 6. 更新文档
	user.Name = "张三 - 更新"
	user.Salary = 5500.00
	if err := userRepo.Update(ctx, user); err != nil {
		log.Printf("Failed to update user: %v", err)
	} else {
		fmt.Printf("Updated user: %v\n", user.Name)
	}

	// 7. 根据条件查询单个文档
	queryWrapper := NewQueryWrapper[User]().
		Eq("name", "李四").
		Gte("age", 25)

	user, err = userRepo.GetOne(ctx, queryWrapper)
	if err != nil {
		log.Printf("Failed to get one user: %v", err)
	} else if user != nil {
		fmt.Printf("Get one user: %v\n", user.Name)
	}

	// 8. 根据条件查询列表
	queryWrapper = NewQueryWrapper[User]().
		Gte("age", 25).
		Lte("salary", 7000.00).
		Eq("active", true).
		OrderBy("age", true)

	userList, err := userRepo.List(ctx, queryWrapper)
	if err != nil {
		log.Printf("Failed to list users: %v", err)
	} else {
		fmt.Printf("List users count: %d\n", len(userList))
		for _, u := range userList {
			fmt.Printf("  - %v (Age: %d)\n", u.Name, u.Age)
		}
	}

	// 9. 分页查询
	page, size := 1, 2
	pageResult, err := userRepo.Page(ctx, queryWrapper, page, size)
	if err != nil {
		log.Printf("Failed to page users: %v", err)
	} else {
		fmt.Printf("Page result: Total=%d, Pages=%d, Current=%d, Size=%d\n",
			pageResult.Total, pageResult.Pages, pageResult.Current, pageResult.Size)
		for _, u := range pageResult.Records {
			fmt.Printf("  - %v (Age: %d)\n", u.Name, u.Age)
		}
	}

	// 10. 统计数量
	count, err := userRepo.Count(ctx, queryWrapper)
	if err != nil {
		log.Printf("Failed to count users: %v", err)
	} else {
		fmt.Printf("Count users: %d\n", count)
	}

	// 11. 检查文档是否存在
	exists, err := userRepo.Exists(ctx, id)
	if err != nil {
		log.Printf("Failed to check user exists: %v", err)
	} else {
		fmt.Printf("User exists: %v\n", exists)
	}

	// 12. 删除文档
	if err := userRepo.DeleteById(ctx, id); err != nil {
		log.Printf("Failed to delete user: %v", err)
	} else {
		fmt.Printf("Deleted user with ID: %s\n", id)
	}

	// 13. 批量删除文档
	if err := userRepo.DeleteBatch(ctx, ids); err != nil {
		log.Printf("Failed to delete batch users: %v", err)
	} else {
		fmt.Printf("Deleted batch users with IDs: %v\n", ids)
	}

	fmt.Println("Example completed successfully!")
}

// ComplexQueryExample 复杂查询示例
func ComplexQueryExample() {
	ctx := context.Background()

	// 创建仓库实例
	userRepo, err := NewBaseRepository[User]()
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	// 复杂查询示例：
	// 1. 年龄大于等于 25 且小于等于 35
	// 2. 薪资大于 5000 或状态为活跃
	// 3. 名字包含 "张" 或邮箱包含 "example.com"
	// 4. 按年龄降序排序
	queryWrapper := NewQueryWrapper[User]().
		Between("age", 25, 35).
		Or(
			NewQueryWrapper[User]().Gt("salary", 5000.00),
			NewQueryWrapper[User]().Eq("active", true),
		).
		Or(
			NewQueryWrapper[User]().Like("name", "张"),
			NewQueryWrapper[User]().Like("email", "example.com"),
		).
		OrderBy("age", false)

	userList, err := userRepo.List(ctx, queryWrapper)
	if err != nil {
		log.Printf("Failed to list users: %v", err)
	} else {
		fmt.Printf("Complex query result count: %d\n", len(userList))
		for _, u := range userList {
			fmt.Printf("  - %v (Age: %d, Salary: %.2f, Active: %v)\n", u.Name, u.Age, u.Salary, u.Active)
		}
	}
}

// AnnotationQueryExample 注解查询示例（类似Java的注解查询）
func AnnotationQueryExample() {
	ctx := context.Background()

	// 创建仓库实例
	userRepo, err := NewBaseRepository[User]()
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	// 定义查询结构体，类似Java的注解查询
	// 注意：这里使用小写字母开头的字段名，因为Go的反射无法访问未导出的字段
	type UserQuery struct {
		// 机器人ID - 等于查询
		ChatId string `query:"type=EQ,field=chatId"`
		// 机器人类型 - 等于查询
		ChatType string `query:"type=EQ,field=chatType"`
		// 登录IP - 等于查询
		LoginIp string `query:"type=EQ,field=loginIp"`
		// 登录端口 - 等于查询
		LoginPort int `query:"type=EQ,field=loginPort"`
		// 在线状态 - 等于查询
		LoginStatus int `query:"type=EQ,field=loginStatus"`
		// 年龄 - 大于等于查询
		AgeGte int `query:"type=GTE,field=age"`
		// 薪资 - 小于等于查询
		SalaryLte float64 `query:"type=LTE,field=salary"`
		// 名字 - 模糊查询
		NameLike string `query:"type=LIKE,field=name"`
	}

	// 创建查询实例并设置条件
	query := UserQuery{
		AgeGte:    25,      // 年龄大于等于25
		SalaryLte: 7000.00, // 薪资小于等于7000
		NameLike:  "张",     // 名字包含"张"
	}

	// 使用注解查询获取列表
	fmt.Println("\n=== 使用注解查询获取用户列表 ===")
	userList, err := userRepo.ListByQueryStruct(ctx, query)
	if err != nil {
		log.Printf("Failed to list users by query struct: %v", err)
	} else {
		fmt.Printf("Annotation query result count: %d\n", len(userList))
		for _, u := range userList {
			fmt.Printf("  - %v (Age: %d, Salary: %.2f)\n", u.Name, u.Age, u.Salary)
		}
	}

	// 使用注解查询进行分页
	fmt.Println("\n=== 使用注解查询进行分页 ===")
	pageResult, err := userRepo.PageByQueryStruct(ctx, query, 1, 2)
	if err != nil {
		log.Printf("Failed to page users by query struct: %v", err)
	} else {
		fmt.Printf("Page result: Total=%d, Pages=%d, Current=%d, Size=%d\n",
			pageResult.Total, pageResult.Pages, pageResult.Current, pageResult.Size)
		for _, u := range pageResult.Records {
			fmt.Printf("  - %v (Age: %d, Salary: %.2f)\n", u.Name, u.Age, u.Salary)
		}
	}

	// 使用注解查询统计数量
	fmt.Println("\n=== 使用注解查询统计数量 ===")
	count, err := userRepo.CountByQueryStruct(ctx, query)
	if err != nil {
		log.Printf("Failed to count users by query struct: %v", err)
	} else {
		fmt.Printf("Annotation query count: %d\n", count)
	}

	fmt.Println("\nAnnotation query example completed successfully!")
}
