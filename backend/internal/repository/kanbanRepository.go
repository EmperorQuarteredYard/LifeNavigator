package repository

import (
	"LifeNavigator/internal/interfaces/repositoryInte"
	"LifeNavigator/internal/models"
	"errors"

	"gorm.io/gorm"
)

func NewKanbanRepository(db *gorm.DB) repositoryInte.KanbanRepository {
	return &kanbanRepository{baseRepository: &baseRepository{db: db}}
}

type kanbanRepository struct {
	*baseRepository
}

func (r *kanbanRepository) Create(kanban *models.Kanban) error {
	return r.create(kanban)
}

func (r *kanbanRepository) GetByID(id uint64) (*models.Kanban, error) {
	var kanban models.Kanban
	err := r.db.First(&kanban, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositoryInte.ErrNotFound
		}
		return nil, err
	}
	return &kanban, nil
}

func (r *kanbanRepository) GetByIDWithProjects(id uint64) (*models.Kanban, error) {
	var kanban models.Kanban
	err := r.db.Preload("Projects").First(&kanban, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositoryInte.ErrNotFound
		}
		return nil, err
	}
	return &kanban, nil
}

func (r *kanbanRepository) ListByUserID(userID uint64) ([]models.Kanban, error) {
	var kanbans []models.Kanban
	err := r.db.Where("user_id = ?", userID).Order("sort_order ASC, created_at DESC").Find(&kanbans).Error
	return kanbans, err
}

func (r *kanbanRepository) Update(kanban *models.Kanban) error {
	result := r.db.Save(kanban)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return repositoryInte.ErrNotFound
	}
	return nil
}

func (r *kanbanRepository) Delete(id uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("kanban_id = ?", id).Delete(&models.KanbanProject{}).Error; err != nil {
			return err
		}
		result := tx.Delete(&models.Kanban{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return repositoryInte.ErrNotFound
		}
		return nil
	})
}

func (r *kanbanRepository) AddProject(kanbanID, projectID uint64) error {
	kp := models.KanbanProject{
		KanbanID:  kanbanID,
		ProjectID: projectID,
	}
	return r.db.Create(&kp).Error
}

func (r *kanbanRepository) RemoveProject(kanbanID, projectID uint64) error {
	return r.db.Where("kanban_id = ? AND project_id = ?", kanbanID, projectID).Delete(&models.KanbanProject{}).Error
}

func (r *kanbanRepository) SetProjects(kanbanID uint64, projectIDs []uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("kanban_id = ?", kanbanID).Delete(&models.KanbanProject{}).Error; err != nil {
			return err
		}
		for _, pid := range projectIDs {
			kp := models.KanbanProject{
				KanbanID:  kanbanID,
				ProjectID: pid,
			}
			if err := tx.Create(&kp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *kanbanRepository) GetProjectIDs(kanbanID uint64) ([]uint64, error) {
	var projectIDs []uint64
	err := r.db.Model(&models.KanbanProject{}).Where("kanban_id = ?", kanbanID).Pluck("project_id", &projectIDs).Error
	return projectIDs, err
}

func (r *kanbanRepository) CheckOwnership(userID, kanbanID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Kanban{}).Where("id = ? AND user_id = ?", kanbanID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *kanbanRepository) GetDefaultKanban(userID uint64) (*models.Kanban, error) {
	var kanban models.Kanban
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&kanban).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositoryInte.ErrNotFound
		}
		return nil, err
	}
	return &kanban, nil
}

func (r *kanbanRepository) SetDefaultKanban(userID, kanbanID uint64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Kanban{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
			return err
		}
		result := tx.Model(&models.Kanban{}).Where("id = ? AND user_id = ?", kanbanID, userID).Update("is_default", true)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return repositoryInte.ErrNotFound
		}
		return nil
	})
}
