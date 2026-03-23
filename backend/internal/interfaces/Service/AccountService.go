package Service

import (
	"time"

	"LifeNavigator/internal/models"
	"LifeNavigator/pkg/dto"
)

// AccountService 定义账户相关的业务接口
// 负责账户的创建、删除、余额调整、查询及关联任务列表等功能
type AccountService interface {
	// CreateAccount 创建新账户
	// 参数:
	//   userID: 用户ID，用于关联账户所有者
	//   account: 待创建的账户模型，应包含账户名称、类型、单位等必要信息
	// 返回值:
	//   *dto.Account: 创建成功后的账户信息（包含自动生成的ID、净余额等）
	//   error: 错误类型可能为 ErrInternal（内部错误）、ErrForbidden（权限不足）等
	CreateAccount(userID uint64, account *models.Account) (*dto.Account, error)

	// DeleteAccount 删除账户
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   accountID: 要删除的账户ID
	// 返回值:
	//   error: 可能返回 ErrForbidden（非所有者）、ErrAccountNotFound（账户不存在）、ErrInternal（内部错误）
	DeleteAccount(userID, accountID uint64) error

	// AdjustBalance 调整账户余额（增加或减少）
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   accountID: 要调整的账户ID
	//   amount: 调整金额（正数为增加，负数为减少）
	// 返回值:
	//   float64: 调整后的最新余额
	//   error: 可能返回 ErrForbidden、ErrAccountNotFound、ErrConcurrentUpdate（并发冲突）、ErrInternal
	AdjustBalance(userID, accountID uint64, amount float64) (float64, error)

	// GetByAccountID 根据账户ID获取账户详情
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   accountID: 账户ID
	// 返回值:
	//   *dto.Account: 账户详情
	//   error: 可能返回 ErrForbidden、ErrAccountNotFound、ErrInternal
	GetByAccountID(userID, accountID uint64) (*dto.Account, error)

	// ListByUserID 获取用户的所有账户列表
	// 参数:
	//   userID: 用户ID
	// 返回值:
	//   *dto.AccountList: 账户列表
	//   error: 可能返回 ErrInternal
	ListByUserID(userID uint64) (*dto.AccountList, error)

	// ListLinkedTask 获取与指定账户关联的任务列表（通过任务预算关联的账户）
	// 参数:
	//   userID: 当前用户ID，用于权限校验
	//   accountID: 账户ID
	//   startTime: 任务创建时间范围开始
	//   endTime: 任务创建时间范围结束
	// 返回值:
	//   *dto.TaskList: 任务列表，包含任务的详细信息及关联的支付项
	//   error: 可能返回 ErrForbidden、ErrInternal
	ListLinkedTask(userID, accountID uint64, startTime, endTime time.Time) (*dto.TaskList, error)

	// GetUserName 获取用户名
	// 参数:
	//   userID: 用户ID
	// 返回值:
	//   string: 用户名
	//   error: 可能返回 ErrUserNotFound、ErrInternal
	GetUserName(userID uint64) (string, error)
}
