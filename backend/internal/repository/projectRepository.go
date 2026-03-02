package repository

import (
	"LifeNavigator/backend/internal/models"
	"errors"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *models.Project) error
	GetByID(id uint64) (*models.Project, error)
	GetByUserID(userID uint64, offset, limit int) ([]models.Project, error)
	Update(project *models.Project) error
	Delete(id uint64) error
	// 可根据需要增加其他查询，如根据刷新间隔等
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *models.Project) error {
	result := r.db.Create(project)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *projectRepository) GetByID(id uint64) (*models.Project, error) {
	var project models.Project
	result := r.db.First(&project, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, result.Error
	}
	return &project, nil
}

func (r *projectRepository) GetByUserID(userID uint64, offset, limit int) ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&projects)
	return projects, result.Error
}

func (r *projectRepository) Update(project *models.Project) error {
	result := r.db.Save(project)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *projectRepository) Delete(id uint64) error {
	result := r.db.Delete(&models.Project{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
