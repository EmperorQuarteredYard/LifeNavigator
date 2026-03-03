package service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"context"
	"errors"
	"log"
)

type ProjectService interface {
	Create(project *models.Project, currentUserID uint64) error
	GetByID(id uint64, currentUserID uint64) (*models.Project, error)
	ListByUserID(currentUserID uint64, offset, limit int) ([]models.Project, error)
	Update(project *models.Project, currentUserID uint64) error
	Delete(id uint64, currentUserID uint64) error
	AddBudget(projectID uint64, budget *models.ProjectBudget, currentUserID uint64) error
	UpdateBudget(budget *models.ProjectBudget, currentUserID uint64) error
	DeleteBudget(budgetID uint64, currentUserID uint64) error
	GetBudgetSummary(projectID uint64, currentUserID uint64) ([]models.ProjectBudget, float64, float64, error)
	GetTaskBudgetSummary(projectID uint64, currentUserID uint64) ([]models.TaskBudget, error)
}

type projectService struct {
	transactor        repository.Transactor
	projectRepo       repository.ProjectRepository
	projectBudgetRepo repository.ProjectBudgetRepository
	taskBudgetRepo    repository.TaskBudgetRepository
	taskRepo          repository.TaskRepository
}

func NewProjectService(
	transactor repository.Transactor,
	projectRepo repository.ProjectRepository,
	projectBudgetRepo repository.ProjectBudgetRepository,
	taskBudgetRepo repository.TaskBudgetRepository,
	taskRepo repository.TaskRepository,
) ProjectService {
	return &projectService{
		transactor:        transactor,
		projectRepo:       projectRepo,
		projectBudgetRepo: projectBudgetRepo,
		taskBudgetRepo:    taskBudgetRepo,
		taskRepo:          taskRepo,
	}
}

func (s *projectService) checkProjectOwnership(projectID uint64, currentUserID uint64) (*models.Project, error) {
	project, err := s.projectRepo.GetByUserID(currentUserID, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrForbidden
		}
		log.Printf("failed to get project %d: %v", projectID, err)
		return nil, ErrInternal
	}
	return project, nil
}
func (s *projectService) Create(project *models.Project, currentUserID uint64) error {
	project.UserID = currentUserID
	if err := s.projectRepo.Create(project); err != nil {
		log.Printf("failed to create project: %v", err)
		return ErrInternal
	}
	return nil
}

func (s *projectService) GetByID(id uint64, currentUserID uint64) (*models.Project, error) {
	return s.checkProjectOwnership(id, currentUserID)
}

func (s *projectService) ListByUserID(currentUserID uint64, offset, limit int) ([]models.Project, error) {
	projects, err := s.projectRepo.ListByUserID(currentUserID, offset, limit)
	if err != nil {
		log.Printf("failed to list projects for user %d: %v", currentUserID, err)
		return nil, ErrInternal
	}
	return projects, nil
}

func (s *projectService) Update(project *models.Project, currentUserID uint64) error {
	// 检查所有权
	_, err := s.checkProjectOwnership(project.ID, currentUserID)
	if err != nil {
		return err
	}
	if err := s.projectRepo.Update(project); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrProjectNotFound
		}
		log.Printf("failed to update project %d: %v", project.ID, err)
		return ErrInternal
	}
	return nil
}

func (s *projectService) Delete(id uint64, currentUserID uint64) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		// 检查项目所有权
		project, err := txRepo.Project.GetByID(id)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrProjectNotFound
			}
			log.Printf("failed to get project %d in transaction: %v", id, err)
			return ErrInternal
		}
		if project.UserID != currentUserID {
			return ErrForbidden
		}

		// 删除项目预算
		if err := txRepo.ProjectBudget.DeleteByProjectID(id); err != nil {
			log.Printf("failed to delete project budgets: %v", err)
			return ErrInternal
		}

		// 获取项目下所有任务
		tasks, _, err := txRepo.Task.ListByProjectID(id, 0, 0)
		if err != nil {
			log.Printf("failed to get tasks for project %d: %v", id, err)
			return ErrInternal
		}
		for _, task := range tasks {
			if err := txRepo.TaskBudget.DeleteByTaskID(task.ID); err != nil {
				log.Printf("failed to delete task budgets for task %d: %v", task.ID, err)
				return ErrInternal
			}
		}

		// 删除项目
		if err := txRepo.Project.Delete(id); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrProjectNotFound
			}
			log.Printf("failed to delete project %d: %v", id, err)
			return ErrInternal
		}
		return nil
	})
}

func (s *projectService) AddBudget(projectID uint64, budget *models.ProjectBudget, currentUserID uint64) error {
	// 检查项目所有权
	_, err := s.checkProjectOwnership(projectID, currentUserID)
	if err != nil {
		return err
	}
	budget.ProjectID = projectID
	if err := s.projectBudgetRepo.Create(budget); err != nil {
		log.Printf("failed to create project budget: %v", err)
		return ErrInternal
	}
	return nil
}

func (s *projectService) UpdateBudget(budget *models.ProjectBudget, currentUserID uint64) error {
	// 获取预算对应的项目ID
	existing, err := s.projectBudgetRepo.GetByID(budget.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to get budget %d: %v", budget.ID, err)
		return ErrInternal
	}
	// 检查项目所有权
	_, err = s.checkProjectOwnership(existing.ProjectID, currentUserID)
	if err != nil {
		return err
	}
	if err := s.projectBudgetRepo.Update(budget); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to update project budget %d: %v", budget.ID, err)
		return ErrInternal
	}
	return nil
}

func (s *projectService) DeleteBudget(budgetID uint64, currentUserID uint64) error {
	// 获取预算对应的项目ID
	existing, err := s.projectBudgetRepo.GetByID(budgetID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to get budget %d: %v", budgetID, err)
		return ErrInternal
	}
	// 检查项目所有权
	_, err = s.checkProjectOwnership(existing.ProjectID, currentUserID)
	if err != nil {
		return err
	}
	if err := s.projectBudgetRepo.Delete(budgetID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to delete project budget %d: %v", budgetID, err)
		return ErrInternal
	}
	return nil
}

func (s *projectService) GetBudgetSummary(projectID uint64, currentUserID uint64) ([]models.ProjectBudget, float64, float64, error) {
	// 检查项目所有权
	_, err := s.checkProjectOwnership(projectID, currentUserID)
	if err != nil {
		return nil, 0, 0, err
	}
	budgets, err := s.projectBudgetRepo.GetByProjectID(projectID)
	if err != nil {
		log.Printf("failed to get project budgets: %v", err)
		return nil, 0, 0, ErrInternal
	}
	var totalBudget, totalUsed float64
	for _, b := range budgets {
		totalBudget += b.Budget
		totalUsed += b.Used
	}
	return budgets, totalBudget, totalUsed, nil
}

func (s *projectService) GetTaskBudgetSummary(projectID uint64, currentUserID uint64) ([]models.TaskBudget, error) {
	// 检查项目所有权
	_, err := s.checkProjectOwnership(projectID, currentUserID)
	if err != nil {
		return nil, err
	}
	tasks, _, err := s.taskRepo.ListByProjectID(projectID, 0, 0)
	if err != nil {
		log.Printf("failed to get tasks for project %d: %v", projectID, err)
		return nil, ErrInternal
	}
	var allTaskBudgets []models.TaskBudget
	for _, task := range tasks {
		budgets, err := s.taskBudgetRepo.GetByTaskID(task.ID)
		if err != nil {
			log.Printf("failed to get task budgets for task %d: %v", task.ID, err)
			return nil, ErrInternal
		}
		allTaskBudgets = append(allTaskBudgets, budgets...)
	}
	merged := models.MergeBudgetItems(allTaskBudgets)
	return merged, nil
}
