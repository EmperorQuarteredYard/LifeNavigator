package repository

import (
	"LifeNavigator/internal/models"
	"errors"

	"gorm.io/gorm"
)

// UserRepository 定义了用户数据访问接口
type UserRepository interface {
	Create(user *models.User) error                      //如果用户名已存在，返回 ErrRecordExist
	GetByID(id uint64) (*models.User, error)             // 根据主键 ID 查询用户,返回用户对象指针，若记录不存在返回 ErrNotFound
	GetByUsername(username string) (*models.User, error) // 返回用户对象指针，若记录不存在返回 ErrNotFound
	GetByEmail(email string) (*models.User, error)       // 返回用户对象指针，若记录不存在返回 ErrNotFound
	Update(user *models.User) error                      // 根据用户对象的 ID 进行更新，若记录不存在返回 ErrNotFound
	UpdateProfile(userID uint64, profile string) error   // profile最大长度为2000字符，超出则裁剪多出部分，若用户不存在返回 ErrNotFound
	Delete(id uint64) error                              // 软删除，若记录不存在返回 ErrNotFound
	List(offset, limit int) ([]models.User, error)       // 分页列出所有用户（未软删除的）
	HardDeleteById(id uint64) error                      // 该操作不可逆，请谨慎使用；若用户不存在返回 ErrNotFound
	SoftDeleteById(id uint64) error                      // 若用户不存在返回 ErrNotFound
}

type userRepository struct {
	*baseRepository
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{baseRepository: &baseRepository{db: db}}
}
func (r *userRepository) HardDeleteById(id uint64) error {
	result := r.db.Unscoped().Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// SoftDeleteById 实现 UserRepository.SoftDeleteById。
// 假设 models.User 包含 gorm.DeletedAt 字段，GORM 会自动执行软删除。
func (r *userRepository) SoftDeleteById(id uint64) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *userRepository) Create(user *models.User) error {
	var count int64
	r.db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return ErrRecordExist
	}

	return r.create(user)
}

func (r *userRepository) GetByID(id uint64) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *userRepository) UpdateProfile(userID uint64, profile string) error {
	const maxProfileLength = 2000
	if len(profile) > maxProfileLength {
		profile = profile[:maxProfileLength]
	}
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Update("profile", profile)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *userRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *userRepository) List(offset, limit int) ([]models.User, error) {
	var users []models.User
	result := r.db.Offset(offset).Limit(limit).Find(&users)
	return users, result.Error
}
