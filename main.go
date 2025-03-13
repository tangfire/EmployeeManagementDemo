package main

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
	"log"
	"time"
)

func main() {
	// 初始化 MySQL
	config.InitMySQL()
	defer config.CloseMySQL()

	// 测试 MySQL 连接
	//if err := config.PingMySQL(); err != nil {
	//	log.Fatalf("MySQL 连接失败: %v", err)
	//} else {
	//	log.Println("MySQL 连接成功")
	//}

	// 自动迁移表结构（根据模型创建或更新表）
	if err := config.DB.AutoMigrate(&models.Employee{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	} else {
		log.Println("员工表已创建/更新")
	}

	// 自动迁移表结构（创建或更新 admin 表）
	if err := config.DB.AutoMigrate(&models.Admin{}); err != nil {
		log.Fatalf("管理员表迁移失败: %v", err)
	} else {
		log.Println("管理员表已创建/更新")
	}

	// 自动迁移表结构（创建或更新 departments 表）
	if err := config.DB.AutoMigrate(&models.Department{}); err != nil {
		log.Fatalf("部门表迁移失败: %v", err)
	} else {
		log.Println("部门表已创建/更新")
	}

	// 自动迁移表结构（创建或更新表）
	if err := config.DB.AutoMigrate(
		&models.LeaveRequest{},
		// 确保已迁移依赖的表（如 Admin 和 Employee）
		&models.Admin{},
		&models.Employee{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	} else {
		log.Println("请假申请表已创建/更新")
	}

	// 自动迁移表结构（创建或更新表）
	if err := config.DB.AutoMigrate(&models.SignRecord{}); err != nil {
		log.Fatalf("签到表迁移失败: %v", err)
	} else {
		log.Println("签到表已创建/更新")
	}

	//createEmployee()                  // 创建
	//getEmployeeByUsername("john_doe") // 查询
	//updateSalary(1, 18000)            // 更新
	//deleteEmployee(1)                 // 删除

	//createAdmin()                          // 创建
	//getAdminByName("admin01")              // 查询
	//updateAdminEmail(1, "new@example.com") // 更新
	//deleteAdmin(1)                         // 删除

	//createDepartment()              // 创建
	//getDepartmentByName("技术部")      // 查询
	//updateDepartmentName(1, "研发中心") // 更新
	//deleteDepartment(1)             // 删除

	//createLeaveRequest()             // 创建申请
	//getLeaveRequestByID(1)           // 查询
	//updateLeaveStatus(1, "approved") // 审批
	//deleteLeaveRequest(1)            // 删除

	//createSignRecord()                                // 创建记录
	//getTodaySignRecords(1001)                         // 查询今日记录
	//updateSignOutTime(1, time.Now().Add(9*time.Hour)) // 更新
	//deleteSignRecord(1)                               // 删除

	// 初始化 Redis
	//config.InitRedis()
	//defer config.Rdb.Close()
	//
	//// 测试 Redis 连接
	//if err := config.PingRedis(); err != nil {
	//	log.Fatalf("Redis 连接失败: %v", err)
	//} else {
	//	log.Println("Redis 连接成功")
	//}

}

// 在 main.go 中添加
func createEmployee() {
	emp := models.Employee{
		Username: "john_doe",
		Password: "secure123",
		Gender:   "男",
		Email:    "john@example.com",
		Salary:   15000,
	}

	result := config.DB.Create(&emp)
	if result.Error != nil {
		log.Fatalf("创建员工失败: %v", result.Error)
	}
	log.Printf("员工创建成功，ID: %d", emp.EmpID)
}

// 根据用户名查询
func getEmployeeByUsername(username string) {
	var emp models.Employee
	result := config.DB.Where("username = ?", username).First(&emp)
	if result.Error != nil {
		log.Printf("查询员工失败: %v", result.Error)
		return
	}
	log.Printf("查询结果: %+v", emp)
}

func updateSalary(EmpID uint, newSalary int) {
	result := config.DB.Model(&models.Employee{}).
		Where("emp_id = ?", EmpID).
		Update("salary", newSalary)
	if result.Error != nil {
		log.Printf("更新薪资失败: %v", result.Error)
	} else {
		log.Printf("薪资更新成功，影响行数: %d", result.RowsAffected)
	}
}

func deleteEmployee(EmpID uint) {
	result := config.DB.Delete(&models.Employee{}, EmpID)
	if result.Error != nil {
		log.Printf("删除员工失败: %v", result.Error)
	} else {
		log.Printf("员工删除成功，影响行数: %d", result.RowsAffected)
	}
}

// 在 main.go 中添加
func createAdmin() {
	admin := models.Admin{
		AdminName:     "admin01",
		AdminPassword: "password123",
		AdminEmail:    "admin01@example.com",
		AdminPhone:    "13800138000",
		Avatar:        "/avatars/admin01.png",
	}

	result := config.DB.Create(&admin)
	if result.Error != nil {
		log.Fatalf("创建管理员失败: %v", result.Error)
	}
	log.Printf("管理员创建成功，ID: %d", admin.AdminID)
}

func getAdminByName(name string) {
	var admin models.Admin
	result := config.DB.Where("admin_name = ?", name).First(&admin)
	if result.Error != nil {
		log.Printf("查询管理员失败: %v", result.Error)
		return
	}
	log.Printf("查询结果: %+v", admin)
}

func updateAdminEmail(adminID uint, newEmail string) {
	result := config.DB.Model(&models.Admin{}).
		Where("admin_id = ?", adminID).
		Update("admin_email", newEmail)
	if result.Error != nil {
		log.Printf("更新邮箱失败: %v", result.Error)
	} else {
		log.Printf("邮箱更新成功，影响行数: %d", result.RowsAffected)
	}
}

func deleteAdmin(adminID uint) {
	result := config.DB.Delete(&models.Admin{}, adminID)
	if result.Error != nil {
		log.Printf("删除管理员失败: %v", result.Error)
	} else {
		log.Printf("管理员删除成功，影响行数: %d", result.RowsAffected)
	}
}

// 在 main.go 中添加
func createDepartment() {
	dept := models.Department{
		Depart: "技术部",
	}

	result := config.DB.Create(&dept)
	if result.Error != nil {
		log.Fatalf("创建部门失败: %v", result.Error)
	}
	log.Printf("部门创建成功，ID: %d", dept.DepID)
}

func getDepartmentByName(name string) {
	var dept models.Department
	result := config.DB.Where("depart = ?", name).First(&dept)
	if result.Error != nil {
		log.Printf("查询部门失败: %v", result.Error)
		return
	}
	log.Printf("查询结果: %+v", dept)
}

func updateDepartmentName(depID int, newName string) {
	result := config.DB.Model(&models.Department{}).
		Where("dep_id = ?", depID).
		Update("depart", newName)
	if result.Error != nil {
		log.Printf("更新部门名称失败: %v", result.Error)
	} else {
		log.Printf("部门名称更新成功，影响行数: %d", result.RowsAffected)
	}
}

func deleteDepartment(depID int) {
	result := config.DB.Delete(&models.Department{}, depID)
	if result.Error != nil {
		log.Printf("删除部门失败: %v", result.Error)
	} else {
		log.Printf("部门删除成功，影响行数: %d", result.RowsAffected)
	}
}

func createLeaveRequest() {
	req := models.LeaveRequest{
		AdminID:   1,    // 管理员ID（需存在）
		EmpID:     1001, // 员工ID（需存在）
		Reason:    "家庭事务",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
		Status:    "pending",
	}

	result := config.DB.Create(&req)
	if result.Error != nil {
		log.Fatalf("创建请假申请失败: %v", result.Error)
	}
	log.Printf("申请已提交，ID: %d", req.ID)
}

func getLeaveRequestByID(id uint) {
	var req models.LeaveRequest
	result := config.DB.Preload("Admin").Preload("User").First(&req, id)
	if result.Error != nil {
		log.Printf("查询失败: %v", result.Error)
		return
	}
	log.Printf("申请详情: %+v", req)
}

func updateLeaveStatus(id uint, newStatus string) {
	// 检查状态合法性
	validStatus := map[string]bool{"pending": true, "approved": true, "rejected": true}
	if !validStatus[newStatus] {
		log.Fatalf("非法状态值: %s", newStatus)
	}

	result := config.DB.Model(&models.LeaveRequest{}).
		Where("id = ?", id).
		Update("status", newStatus)
	if result.Error != nil {
		log.Printf("状态更新失败: %v", result.Error)
	} else {
		log.Printf("状态已更新，影响行数: %d", result.RowsAffected)
	}
}

func deleteLeaveRequest(id uint) {
	result := config.DB.Delete(&models.LeaveRequest{}, id)
	if result.Error != nil {
		log.Printf("删除失败: %v", result.Error)
	} else {
		log.Printf("删除成功，影响行数: %d", result.RowsAffected)
	}
}

// 在 main.go 中添加
func createSignRecord() {
	record := models.SignRecord{
		EmpID:       1001,                          // 员工ID（需存在）
		SignInTime:  time.Now(),                    // 当前时间签到
		SignOutTime: time.Now().Add(8 * time.Hour), // 8小时后签退
	}

	result := config.DB.Create(&record)
	if result.Error != nil {
		log.Fatalf("创建签到记录失败: %v", result.Error)
	}
	log.Printf("签到记录已创建，ID: %d", record.ID)
}

func getTodaySignRecords(empID uint) {
	var records []models.SignRecord
	today := time.Now().Format("2006-01-02")

	result := config.DB.Where("emp_id = ? AND DATE(sign_in_time) = ?", empID, today).
		Preload("Employee"). // 预加载关联的员工信息
		Find(&records)

	if result.Error != nil {
		log.Printf("查询失败: %v", result.Error)
		return
	}
	log.Printf("查询结果: %+v", records)
}

func updateSignOutTime(recordID uint, newSignOutTime time.Time) {
	result := config.DB.Model(&models.SignRecord{}).
		Where("id = ?", recordID).
		Update("sign_out_time", newSignOutTime)
	if result.Error != nil {
		log.Printf("更新签退时间失败: %v", result.Error)
	} else {
		log.Printf("签退时间已更新，影响行数: %d", result.RowsAffected)
	}
}

func deleteSignRecord(recordID uint) {
	result := config.DB.Delete(&models.SignRecord{}, recordID)
	if result.Error != nil {
		log.Printf("删除失败: %v", result.Error)
	} else {
		log.Printf("删除成功，影响行数: %d", result.RowsAffected)
	}
}
