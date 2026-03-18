package repository

import (
	"LifeNavigator/internal/models"
	"errors"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *models.Project, userIDs []uint64) error
	GetByID(id uint64) (*models.Project, error)
	ListByUserID(userID uint64, offset, limit int) ([]models.Project, error)
	Update(project *models.Project) error
	Delete(id uint64) error
	CheckOwnership(userID, projectID uint64) (bool, error)
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

type projectRepository struct {
	db *gorm.DB
}

func (r *projectRepository) CheckOwnership(userID, projectID uint64) (bool, error) {
	var count int64
	err := r.db.Table("project_users").
		Where("user_id = ? AND project_id = ?", userID, projectID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *projectRepository) Create(project *models.Project, userIDs []uint64) error {
	if len(userIDs) == 0 {
		return errors.New("at least one user_id is required")
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(project).Error; err != nil {
			return err
		}

		for _, uid := range userIDs {
			if err := tx.Exec("INSERT INTO project_users (user_id, project_id) VALUES (?, ?)", uid, project.ID).Error; err != nil {
				return err
			}
		}
		return nil
	})
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

func (r *projectRepository) ListByUserID(userID uint64, offset, limit int) ([]models.Project, error) {
	var projects []models.Project
	result := r.db.Joins("JOIN project_users ON project_users.project_id = projects.id").
		Where("project_users.user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Find(&projects)
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
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM project_users WHERE project_id = ?", id).Error; err != nil {
			return err
		}
		result := tx.Delete(&models.Project{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	})
}
