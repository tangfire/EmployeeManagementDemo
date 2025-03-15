package services

import (
	"EmployeeManagementDemo/config"
	"EmployeeManagementDemo/models"
)

func GetDepartmentAvgSalaries() ([]models.DepartmentAvgSalaryDTO, error) {
	var results []models.DepartmentAvgSalaryDTO

	err := config.DB.Model(&models.Employee{}).
		Select(`
			departments.dep_id,
			departments.depart as depart,
			AVG(employees.salary) as avg_salary
		`).
		Joins("LEFT JOIN departments ON employees.dep_id = departments.dep_id").
		Group("departments.dep_id, departments.depart").
		Order("avg_salary DESC").
		Scan(&results).Error

	return results, err
}

func GetDepartmentHeadcounts() ([]models.DepartmentHeadcountDTO, error) {
	var results []models.DepartmentHeadcountDTO

	// Step 1: 获取各部门人数
	err := config.DB.Model(&models.Employee{}).
		Select(`
            departments.dep_id,
            departments.depart as depart,
            COUNT(employees.emp_id) as headcount
        `).
		Joins("LEFT JOIN departments ON employees.dep_id = departments.dep_id").
		Group("departments.dep_id, departments.depart").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	// Step 2: 计算总人数
	total := 0
	for _, dept := range results {
		total += dept.Headcount
	}

	// Step 3: 计算百分比
	for i := range results {
		if total > 0 {
			results[i].Percentage = float64(results[i].Headcount) / float64(total) * 100
		} else {
			results[i].Percentage = 0.0
		}
	}

	return results, nil
}
