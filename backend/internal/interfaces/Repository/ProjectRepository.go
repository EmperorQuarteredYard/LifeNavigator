package Repository

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/pkg/permission"
	"LifeNavigator/pkg/scheduler"
)

// ProjectRepository 定义了项目相关的数据访问接口。
type ProjectRepository interface {
	// Create 创建新项目并关联指定的用户。
	// 参数 userIDs: 至少一个用户 ID，用于将项目与用户关联（project_users 表）。
	// 错误：如果 userIDs 为空返回 ErrInvalidInput。
	Create(project *models.Project, userIDs []uint64) error

	// GetByID 根据项目 ID 查询项目信息。
	GetByID(id uint64) (*models.Project, error)

	// ListByUserID 分页查询用户有权限访问的项目列表。
	// 返回结果按 id 升序排列。
	ListByUserID(userID uint64, offset, limit int) ([]models.Project, error)

	// Update 更新项目信息（全量更新）。
	// 错误：如果项目不存在返回 ErrNotFound。
	Update(project *models.Project) error

	// Delete 删除项目，同时删除其与用户的关联关系。
	// 错误：如果项目不存在返回 ErrNotFound。
	Delete(id uint64) error

	// CheckOwnership 检查用户是否为项目的所有者（owner 字段匹配）。
	CheckOwnership(userID, projectID uint64) (bool, error)

	// GetRefreshInformation 获取需要刷新时间的项目列表（用于调度器）。
	// 返回按 last_refresh 升序排列的项目 ID 和最后刷新时间。
	// 参数 page, pageSize: 分页参数。
	GetRefreshInformation(page, pageSize int) (list []*scheduler.Schedule, total int64, err error)

	// CheckAccessibility 检查用户对项目是否有指定的操作权限。
	// 参数 operation: 要检查的操作类型（如 permission.Read, permission.Write）。
	// 参数 isUserGuest: 标识用户是否为访客模式（用于权限判断）。
	// 权限判断逻辑：
	//   - 项目所有者：检查 RoleOwner 权限。
	//   - 项目成员（project_users 表中存在）：检查 RoleWorkmate 权限。
	//   - 其他用户：若 isUserGuest 为 true，检查 RoleGuest 权限；否则检查 RoleViewer 权限。
	// 错误：如果项目不存在返回 ErrNotFound。
	CheckAccessibility(userID uint64, projectID uint64, operation permission.Operation, isUserGuest bool) (bool, error)
}
