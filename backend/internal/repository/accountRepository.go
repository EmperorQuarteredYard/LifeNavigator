package repository

import (
	"LifeNavigator/internal/interfaces/repositoryInte"
	"LifeNavigator/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

func NewAccountRepository(db *gorm.DB) repositoryInte.AccountRepository {
	return &accountRepository{baseRepository: &baseRepository{db: db}, updateMaxTry: 3}
}

type accountRepository struct {
	*baseRepository
	updateMaxTry int
}

func (a *accountRepository) SetUpdateMaxTry(maxTry int) {
	a.updateMaxTry = maxTry
}

func (a *accountRepository) CheckOwnership(userID, accountID uint64) (bool, error) {
	var count int64
	err := a.db.Table("account_users").
		Where("user_id = ? AND account_id = ?", userID, accountID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (a *accountRepository) AdjustNetBalance(accountID uint64, amount float64) (float64, error) {
	if amount == 0 {
		account, err := a.GetByID(accountID)
		if err != nil {
			return 0, err
		}
		return account.NetBalance, nil
	}

	for i := 0; i < a.updateMaxTry; i++ {
		account := &models.Account{}
		err := a.db.First(account, accountID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, repositoryInte.ErrNotFound
			}
			return 0, err
		}

		newNetBalance := account.NetBalance + amount

		result := a.db.Model(&models.Account{}).
			Where("id = ? AND version = ?", accountID, account.Version).
			Updates(map[string]interface{}{
				"net_balance": newNetBalance,
				"version":     account.Version + 1,
			})

		if result.Error != nil {
			return 0, result.Error
		}

		if result.RowsAffected == 0 {
			if i == a.updateMaxTry-1 {
				return 0, repositoryInte.ErrConcurrentUpdate
			}
			time.Sleep(10 * time.Millisecond)
			continue
		}

		return newNetBalance, nil
	}

	return 0, repositoryInte.ErrConcurrentUpdate
}

func (a *accountRepository) AdjustBalance(accountID uint64, amount float64) (float64, error) {
	if amount == 0 {
		account, err := a.GetByID(accountID)
		if err != nil {
			return 0, err
		}
		return account.Balance, nil
	}

	for i := 0; i < a.updateMaxTry; i++ {
		account := &models.Account{}
		err := a.db.First(account, accountID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, repositoryInte.ErrNotFound
			}
			return 0, err
		}

		newBalance := account.Balance + amount

		result := a.db.Model(&models.Account{}).
			Where("id = ? AND version = ?", accountID, account.Version).
			Updates(map[string]interface{}{
				"balance": newBalance,
				"version": account.Version + 1,
			})

		if result.Error != nil {
			return 0, result.Error
		}

		if result.RowsAffected == 0 {
			if i == a.updateMaxTry-1 {
				return 0, repositoryInte.ErrConcurrentUpdate
			}
			time.Sleep(10 * time.Millisecond)
			continue
		}

		return newBalance, nil
	}

	return 0, repositoryInte.ErrConcurrentUpdate
}

func (a *accountRepository) GetByID(accountID uint64) (*models.Account, error) {
	account := &models.Account{}
	err := a.db.First(account, accountID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositoryInte.ErrNotFound
		}
		return nil, err
	}
	return account, nil
}

func (a *accountRepository) ListByUserID(userID uint64) ([]models.Account, error) {
	var accounts []models.Account
	err := a.db.Joins("JOIN account_users ON account_users.account_id = accounts.id").
		Where("account_users.user_id = ?", userID).
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (a *accountRepository) Create(account *models.Account, userIDs []uint64) (*models.Account, error) {
	if len(userIDs) == 0 {
		return nil, errors.New("at least one user_id is required")
	}
	if account.Type == "" {
		return nil, errors.New("type is required")
	}

	err := a.db.Transaction(func(tx *gorm.DB) error {
		if err := a.createWithTX(tx, account); err != nil {
			return err
		}

		for _, uid := range userIDs {
			if err := tx.Exec("INSERT INTO account_users (user_id, account_id) VALUES (?, ?)", uid, account.ID).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a *accountRepository) Delete(account *models.Account) error {
	if account.ID == 0 {
		return errors.New("account ID is required for delete")
	}

	err := a.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM account_users WHERE account_id = ?", account.ID).Error; err != nil {
			return err
		}
		if err := tx.Delete(account).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}
