package repositoryInte

import (
	"LifeNavigator/internal/models"
)

// KanbanRepository 定义看板的数据访问接口。
// 看板用于组织和展示项目（Projects），支持设置默认看板、排序等。
type KanbanRepository interface {
	// Create 创建一个新看板。
	Create(kanban *models.Kanban) error

	// GetByID 根据 ID 查询看板（不预加载关联的项目）。
	GetByID(id uint64) (*models.Kanban, error)

	// GetByIDWithProjects 根据 ID 查询看板，并预加载其包含的所有项目。
	GetByIDWithProjects(id uint64) (*models.Kanban, error)

	// ListByUserID 列出指定用户的所有看板，按 sort_order 升序、创建时间降序排序。
	ListByUserID(userID uint64) ([]models.Kanban, error)

	// Update 更新看板信息。
	Update(kanban *models.Kanban) error

	// Delete 删除看板，同时清理其与项目的关联关系。
	Delete(id uint64) error

	// AddProject 向看板中添加一个项目（建立关联）。
	AddProject(kanbanID, projectID uint64) error

	// RemoveProject 从看板中移除一个项目。
	RemoveProject(kanbanID, projectID uint64) error

	// SetProjects 设置看板包含的项目列表（先删除原有全部关联，再添加新列表）。
	SetProjects(kanbanID uint64, projectIDs []uint64) error

	// GetProjectIDs 获取看板包含的所有项目 ID 列表。
	GetProjectIDs(kanbanID uint64) ([]uint64, error)

	// CheckOwnership 检查用户是否为看板的创建者（拥有者）。
	CheckOwnership(userID, kanbanID uint64) (bool, error)

	// GetDefaultKanban 获取用户的默认看板。
	// 若不存在默认看板，返回 ErrNotFound。
	GetDefaultKanban(userID uint64) (*models.Kanban, error)

	// SetDefaultKanban 将指定看板设为用户的默认看板。
	// 会先将用户的所有看板的 is_default 置为 false，再设置新的默认看板。
	SetDefaultKanban(userID, kanbanID uint64) error
}
