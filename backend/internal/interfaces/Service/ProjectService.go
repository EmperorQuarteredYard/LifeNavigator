package Service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/pkg/dto"
)

// ProjectService 定义项目管理接口
type ProjectService interface {
	// Create 创建项目
	// 参数:
	//   userID: 当前用户ID，作为项目所有者
	//   project: 项目模型（应包含名称、描述等）
	// 返回值:
	//   *dto.ProjectResponse: 创建成功的项目详情
	//   error: 可能返回 ErrForbidden、ErrInternal
	Create(userID uint64, project *models.Project) (*dto.ProjectResponse, error)

	// GetByID 获取项目详情
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   id: 项目ID
	// 返回值:
	//   *dto.ProjectResponse: 项目详情（包含预算信息）
	//   error: 可能返回 ErrForbidden、ErrProjectNotFound、ErrInternal
	GetByID(userID, id uint64) (*dto.ProjectResponse, error)

	// ListByUserID 获取用户的项目列表（分页）
	// 参数:
	//   userID: 用户ID
	//   offset: 偏移量
	//   limit: 每页数量
	// 返回值:
	//   *dto.ProjectListResponse: 项目列表
	//   error: 可能返回 ErrInternal
	ListByUserID(userID uint64, offset, limit int) (*dto.ProjectListResponse, error)

	// Update 更新项目信息
	// 参数:
	//   userID: 当前用户ID，用于权限校验（需要项目所有者权限）
	//   project: 更新后的项目模型（必须包含ID）
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrProjectNotFound、ErrInternal
	Update(userID uint64, project *models.Project) error

	// Delete 删除项目
	// 参数:
	//   userID: 当前用户ID，用于权限校验（需要项目所有者权限）
	//   id: 项目ID
	// 返回值:
	//   error: 可能返回 ErrForbidden、ErrProjectNotFound、ErrInternal
	Delete(userID, id uint64) error
}
