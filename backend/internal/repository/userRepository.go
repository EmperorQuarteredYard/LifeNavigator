package repository

import (
	"LifeNavigator/internal/interfaces/repositoryInte"
	"LifeNavigator/internal/models"
	"errors"

	"gorm.io/gorm"
)

type userRepository struct {
	*baseRepository
}

func NewUserRepository(db *gorm.DB) repositoryInte.UserRepository {
	return &userRepository{baseRepository: &baseRepository{db: db}}
}

func (r *userRepository) HardDeleteById(id uint64) error {
	result := r.db.Unscoped().Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repositoryInte.ErrNotFound
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
		return repositoryInte.ErrNotFound
	}
	return nil
}
func (r *userRepository) UpdateAvatar(userID uint64, avatarURL string) error {
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Update("avatar", avatarURL)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repositoryInte.ErrNotFound
	}
	return nil
}

func (r *userRepository) Create(user *models.User) error {
	var count int64
	r.db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return repositoryInte.ErrRecordExist
	}

	return r.create(user)
}

func (r *userRepository) GetByID(id uint64) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repositoryInte.ErrNotFound
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
			return nil, repositoryInte.ErrNotFound
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
			return nil, repositoryInte.ErrNotFound
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
		return repositoryInte.ErrNotFound
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
		return repositoryInte.ErrNotFound
	}
	return nil
}

func (r *userRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repositoryInte.ErrNotFound
	}
	return nil
}

func (r *userRepository) List(offset, limit int) ([]models.User, error) {
	var users []models.User
	result := r.db.Offset(offset).Limit(limit).Find(&users)
	return users, result.Error
}
