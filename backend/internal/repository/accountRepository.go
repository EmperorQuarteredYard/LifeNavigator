package repository

import (
	"LifeNavigator/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(*models.Account) (*models.Account, error)
	Delete(*models.Account) error
	AdjustBalance(userID, accountID uint64, amount float64) (float64, error)
	GetByID(userID, accountID uint64) (*models.Account, error)
	ListByUserID(userID uint64) ([]models.Account, error)
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return accountRepository{db: db}
}

type accountRepository struct {
	db *gorm.DB
}

func (a accountRepository) AdjustBalance(userID, accountID uint64, amount float64) (float64, error) {
	var maxRetries = 3
	if amount == 0 {
		account, err := a.GetByID(userID, accountID)
		if err != nil {
			return 0, err
		}
		return account.Balance, nil
	}

	for i := 0; i < maxRetries; i++ {
		account := &models.Account{}
		err := a.db.Where("user_id = ? AND id = ?", userID, accountID).First(account).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, ErrNotFound
			}
			return 0, err
		}

		newBalance := account.Balance + amount

		result := a.db.Model(&models.Account{}).
			Where("user_id = ? AND id = ? AND version = ?", userID, accountID, account.Version).
			Updates(map[string]interface{}{
				"balance": newBalance,
				"version": account.Version + 1,
			})

		if result.Error != nil {
			return 0, result.Error
		}

		if result.RowsAffected == 0 {
			if i == maxRetries-1 {
				return 0, ErrConcurrentUpdate
			}
			time.Sleep(10 * time.Millisecond) // 短暂等待后重试
			continue
		}

		return newBalance, nil
	}

	return 0, ErrConcurrentUpdate
}

// GetByID 根据用户ID和账户ID查询单个账户
func (a accountRepository) GetByID(userID, accountID uint64) (*models.Account, error) {
	account := &models.Account{}
	err := a.db.Where("user_id = ? AND id = ?", userID, accountID).First(account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return account, nil
}

// ListByUserID 查询指定用户的所有账户
func (a accountRepository) ListByUserID(userID uint64) ([]models.Account, error) {
	var accounts []models.Account
	err := a.db.Where("user_id = ?", userID).Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

// Create 创建新账户
// 注意：传入的 account 对象中应包含必要字段（如 UserID, Type, Balance 等）
// 创建后，account 的 ID 和 Version 会被自动填充（ID自增，Version默认0）
func (a accountRepository) Create(account *models.Account) (*models.Account, error) {
	// 可在此添加业务校验（如检查必填字段）
	if account.UserID == 0 {
		return nil, errors.New("user_id is required")
	}
	if account.Type == "" {
		return nil, errors.New("type is required")
	}

	err := a.db.Create(account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

// Delete 删除指定账户
func (a accountRepository) Delete(account *models.Account) error {
	// 确保 account 对象包含有效的主键（ID）
	if account.ID == 0 {
		return errors.New("account ID is required for delete")
	}
	err := a.db.Delete(account).Error
	if err != nil {
		return err
	}
	return nil
}
