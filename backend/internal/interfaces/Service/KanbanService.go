package Service

import (
	"LifeNavigator/pkg/dto"
)

// KanbanService 定义看板管理接口
type KanbanService interface {
	// Create 创建看板
	// 参数:
	//   userID: 当前用户ID
	//   req: 创建看板的请求参数（包含名称、描述、项目ID列表）
	// 返回值:
	//   *dto.KanbanResponse: 创建成功的看板详情
	//   error: 可能返回 ErrForbidden、ErrInternal
	Create(userID uint64, req *dto.CreateKanbanRequest) (*dto.KanbanResponse, error)

	// GetByID 获取看板详情
	// 参数:
	//   userID: 当前用户ID
	//   id: 看板ID
	// 返回值:
	//   *dto.KanbanResponse: 看板详情（包含关联的项目信息）
	//   error: 可能返回 ErrForbidden、ErrKanbanNotFound、ErrInternal
	GetByID(userID, id uint64) (*dto.KanbanResponse, error)

	// ListByUserID 获取用户的所有看板列表
	// 参数:
	//   userID: 用户ID
	// 返回值:
	//   *dto.KanbanListResponse: 看板列表
	//   error: 可能返回 ErrInternal
	ListByUserID(userID uint64) (*dto.KanbanListResponse, error)

	// Update 更新看板信息
	// 参数:
	//   userID: 当前用户ID
	//   id: 看板ID
	//   req: 更新请求（可修改名称、描述、排序、关联项目等）
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrKanbanNotFound、ErrInternal
	Update(userID uint64, id uint64, req *dto.UpdateKanbanRequest) error

	// Delete 删除看板
	// 参数:
	//   userID: 当前用户ID
	//   id: 看板ID
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrKanbanNotFound、ErrInternal
	Delete(userID, id uint64) error

	// GetKanbanTasks 获取看板下的任务列表
	// 参数:
	//   userID: 当前用户ID
	//   kanbanID: 看板ID
	//   status: 可选的任务状态过滤（nil表示不过滤）
	//   page: 分页页码
	//   pageSize: 每页数量
	// 返回值:
	//   *dto.KanbanTaskListResponse: 任务列表（包含任务详细信息、前置任务依赖、支付记录等）
	//   error: 可能返回 ErrForbidden、ErrKanbanNotFound、ErrInternal
	GetKanbanTasks(userID, kanbanID uint64, status *uint8, page, pageSize int) (*dto.KanbanTaskListResponse, error)

	// GetDefaultKanbanTasks 获取用户默认看板的任务列表
	// 参数:
	//   userID: 当前用户ID
	//   status: 可选的任务状态过滤
	//   page: 分页页码
	//   pageSize: 每页数量
	// 返回值:
	//   *dto.KanbanTaskListResponse: 任务列表
	//   error: 可能返回 ErrInternal
	GetDefaultKanbanTasks(userID uint64, status *uint8, page, pageSize int) (*dto.KanbanTaskListResponse, error)

	// SetDefaultKanban 设置用户的默认看板
	// 参数:
	//   userID: 当前用户ID
	//   kanbanID: 看板ID
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrKanbanNotFound、ErrInternal
	SetDefaultKanban(userID, kanbanID uint64) error

	// AddProjectInCheck 将项目添加到看板的待检清单（具体业务含义待实现）
	// 参数:
	//   userID: 当前用户ID
	//   projectID: 项目ID
	//   kanbanID: 看板ID
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrProjectNotFound、ErrKanbanNotFound、ErrInternal
	AddProjectInCheck(userID, projectID, kanbanID uint64) error
}
