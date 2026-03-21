package repository

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/pkg/permission"
	"LifeNavigator/pkg/scheduler"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(project *models.Project, userIDs []uint64) error
	GetByID(id uint64) (*models.Project, error)
	ListByUserID(userID uint64, offset, limit int) ([]models.Project, error)
	Update(project *models.Project) error
	Delete(id uint64) error
	CheckOwnership(userID, projectID uint64) (bool, error)
	GetRefreshInformation(page, pageSize int) (list []*scheduler.Schedule, total int64, err error)
	CheckAccessibility(userID uint64, projectID uint64, operation permission.Operation, isUserGuest bool) (bool, error)
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{baseRepository: &baseRepository{db: db}}
}

type projectRepository struct {
	*baseRepository
}

func (r *projectRepository) CheckAccessibility(userID uint64, projectID uint64, operation permission.Operation, isUserGuest bool) (bool, error) {
	var project models.Project

	err := r.db.Where("id = ?", projectID).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, ErrNotFound
		}
		return false, err
	}
	switch {
	case project.Owner == userID:
		if project.Permission.Has(permission.RoleOwner, operation) {
			return true, nil
		}
		return false, nil
	default:
		var count int64
		err = r.db.Table("project_users").
			Where("project_id = ? AND user_id = ?", projectID, userID).
			Count(&count).Error
		if err != nil {
			return false, err
		}
		if count > 0 {
			if project.Permission.Has(permission.RoleWorkmate, operation) {
				return true, nil
			}
		}

		if isUserGuest {
			if project.Permission.Has(permission.RoleGuest, operation) {
				return true, nil
			}
			return false, nil
		}
		if project.Permission.Has(permission.RoleViewer, operation) {
			return true, nil
		}
		return false, nil
	}
}

func (r *projectRepository) GetRefreshInformation(page, pageSize int) (list []*scheduler.Schedule, total int64, err error) {
	var limit, offset int
	if err = r.db.Model(&models.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if int64(page*pageSize) > total {
		limit = int(total % int64(pageSize))
		offset = int(total - int64(limit))
	} else {
		limit = pageSize
		offset = page * pageSize
	}

	rows, err := r.db.
		Model(&models.Project{}).
		Select("id", "last_refresh").
		Offset(offset).
		Limit(limit).
		Order("last_refresh asc").
		Rows()
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	var id uint64
	var lastRefresh time.Time
	list = make([]*scheduler.Schedule, limit)
	index := 0
	for rows.Next() {
		if err = rows.Scan(&id, &lastRefresh); err != nil {
			return nil, 0, err
		}
		list[index] = &scheduler.Schedule{}
		index++
	}
	return list, total, nil
}

func (r *projectRepository) CheckOwnership(userID, projectID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Project{}).
		Where("owner = ? AND id = ?", userID, projectID).
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
		if err := r.createWithTX(tx, project); err != nil {
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
		Order("id asc").
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
