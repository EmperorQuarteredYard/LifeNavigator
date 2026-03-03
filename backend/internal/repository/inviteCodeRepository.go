package repository

import (
	"LifeNavigator/internal/models"
	"errors"

	"gorm.io/gorm"
)

type InviteCodeRepository interface {
	//Create 创建一个InviteCode对象，其中ID强制置零后添加
	Create(code *models.InviteCode) error

	// UseCodeByToken 使用给定的令牌来消耗一个邀请码。
	//
	// 返回的错误包括：
	//   - ErrNotFound：指定的 token 不存在。
	//   - ErrInviteCodeUsed：邀请码已被用完（Count >= Budget）。
	//   - 其他数据库错误。
	UseCodeByToken(token string) error

	// Delete 删除指定的邀请码记录。
	// 该方法执行硬删除（直接从数据库中移除记录）。
	// 如果记录不存在，返回 ErrNotFound；其他数据库错误原样返回。
	Delete(code *models.InviteCode) error

	// FindByToken 根据令牌查询邀请码。
	// 返回匹配的邀请码指针；如果记录不存在，返回 ErrNotFound。
	FindByToken(token string) (*models.InviteCode, error)

	// FindByID 根据ID查询邀请码。
	// 返回匹配的邀请码指针；如果记录不存在，返回 ErrNotFound。
	FindByID(id int64) (*models.InviteCode, error)
	GetByUserID(uid uint64, offset, limit int) ([]models.InviteCode, error)
}

type inviteCodeRepository struct {
	db *gorm.DB
}

func NewInviteCodeRepository(db *gorm.DB) InviteCodeRepository {
	return &inviteCodeRepository{db: db}
}

func (r *inviteCodeRepository) GetByUserID(userID uint64, offset, limit int) ([]models.InviteCode, error) {
	var codes []models.InviteCode
	err := r.db.Where("invited_by = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&codes).Error
	if err != nil {
		return nil, err
	}
	return codes, nil
}

func (r *inviteCodeRepository) Create(code *models.InviteCode) error {
	err := r.db.Create(code).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *inviteCodeRepository) UseCodeByToken(token string) error {
	result := r.db.Model(&models.InviteCode{}).
		Where("token = ? AND count < amount", token).
		Update("count", gorm.Expr("count + ?", 1))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		var code models.InviteCode
		err := r.db.Where("token = ?", token).First(&code).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if code.Count >= code.Amount {
			return ErrInviteCodeUsed
		}
		// 理论上不会走到这里，但为安全返回未知错误
		return ErrUnexpected
	}

	return nil
}
func (r *inviteCodeRepository) Delete(code *models.InviteCode) error {
	err := r.db.Delete(code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (r *inviteCodeRepository) FindByToken(token string) (*models.InviteCode, error) {
	var code models.InviteCode
	err := r.db.Where("token = ?", token).First(&code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &code, nil
}

func (r *inviteCodeRepository) FindByID(id int64) (*models.InviteCode, error) {
	var code models.InviteCode
	err := r.db.First(&code, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &code, nil
}
