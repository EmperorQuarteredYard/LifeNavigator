package service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"LifeNavigator/pkg/refresh"
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

// TODO 这里刷新的逻辑根本不对，应当当设置一个定时器来查询是否需要更新
func (s *projectService) checkIfNeedRefreshBudget(tx repository.TxRepositories, projectID, currentUserID uint64) (bool, error) {
	project, err := s.projectRepo.GetByUserID(currentUserID, projectID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, ErrForbidden
		}
		log.Printf("failed to get project %d: %v", projectID, err)
		return false, ErrInternal
	}
	budgets, err := tx.ProjectBudget.GetByProjectID(projectID)
	if err != nil {
		log.Println(err)
		return false, ErrInternal
	}
	if !refresh.ShouldRefresh(project.RefreshInterval, project.LastRefresh) {
		return false, nil
	}
	for _, budget := range budgets {
		_, err = tx.Account.AdjustNetBalance(currentUserID, budget.AccountID, -1*budget.Used)
		if err != nil {
			log.Print("fail to auto refresh budget of project %d ,account %d :%v\n", projectID, budget.AccountID, err)
			return false, ErrInternal
		}
		err = tx.ProjectBudget.SetUsedZero(budget.ID, budget.ProjectID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return false, ErrBudgetNotFound
			}
			log.Printf("fail to auto refresh budget of project %d ,account %d :%v\n", projectID, budget.AccountID, err)
			return false, ErrInternal
		}
	}
	return true, nil
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

		// 删除项目预算并返还未使用的
		var budgets []models.ProjectBudget
		if budgets, err = txRepo.ProjectBudget.GetByProjectID(id); err != nil {
			return ErrProjectBudgetNotFound
		}
		for _, budget := range budgets {
			err = txRepo.ProjectBudget.Delete(budget.ID)
			if err != nil {
				log.Printf("failed to delete project %d: %v", budget.ID, err)
				return ErrInternal
			}
			_, err = txRepo.Account.AdjustNetBalance(currentUserID, budget.AccountID, budget.Budget)
			if err != nil {
				log.Printf("failed to update project %d: %v", budget.ID, err)
				return ErrInternal
			}
		}

		// 获取项目下所有任务
		tasks, _, err := txRepo.Task.ListByProjectID(id, 0, 0)
		if err != nil {
			log.Printf("failed to get tasks for project %d: %v", id, err)
			return ErrInternal
		}
		for _, task := range tasks {
			if err := txRepo.TaskPayment.DeleteByTaskID(task.ID); err != nil {
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
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		_, err := s.checkProjectOwnership(projectID, currentUserID)
		if err != nil {
			return err
		}
		budget.ProjectID = projectID
		if err := txRepo.ProjectBudget.Create(budget); err != nil {
			log.Printf("failed to create project budget: %v", err)
			return ErrInternal
		}
		if budget.AccountID != 0 {
			if _, err = txRepo.Account.AdjustNetBalance(currentUserID, budget.AccountID, -1*budget.Budget); err != nil {
				log.Printf("failed to create project budget: %v", err)
				return ErrInternal
			}
		}
		return nil
	})
}

func (s *projectService) UpdateBudget(budget *models.ProjectBudget, currentUserID uint64) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		// 获取预算对应的项目 ID
		existing, err := txRepo.ProjectBudget.GetByID(budget.ID)
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
		//正式开始事务
		var oldBudget *models.ProjectBudget
		if oldBudget, err = txRepo.ProjectBudget.GetByID(budget.ID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			log.Printf("failed to get budget %d: %v", budget.ID, err)
			return ErrInternal
		}
		//当改变了账户时
		if budget.AccountID != oldBudget.AccountID {
			if oldBudget.AccountID != 0 {
				_, err = txRepo.Account.AdjustNetBalance(currentUserID, oldBudget.AccountID, oldBudget.Budget)
				if err != nil {
					if errors.Is(err, repository.ErrNotFound) {
						return ErrBudgetNotFound
					}
					log.Printf("failed to update project budget: %v", err)
					return ErrInternal
				}
			}
			if budget.AccountID != 0 {
				_, err = txRepo.Account.AdjustNetBalance(currentUserID, budget.AccountID, -1*budget.Budget)
				if err != nil {
					if errors.Is(err, repository.ErrNotFound) {
						return ErrBudgetNotFound
					}
					log.Printf("failed to update project budget: %v", err)
					return ErrInternal
				}
			}
		} else if budget.Budget != oldBudget.Budget && budget.AccountID != 0 { //否则，当改变了预算额时
			_, err = txRepo.Account.AdjustNetBalance(currentUserID, oldBudget.AccountID, oldBudget.Budget-budget.Budget)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return ErrBudgetNotFound
				}
				log.Printf("failed to update project budget: %v", err)
				return ErrInternal
			}
		}
		if err = txRepo.ProjectBudget.Update(budget); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			log.Printf("failed to update project budget %d: %v", budget.ID, err)
			return ErrInternal
		}
		return nil

	})
}

func (s *projectService) DeleteBudget(budgetID uint64, currentUserID uint64) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error { // 获取预算对应的项目ID
		existing, err := txRepo.ProjectBudget.GetByID(budgetID)
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
		budget, err := txRepo.ProjectBudget.GetByID(budgetID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			log.Printf("failed to get budget %d: %v", budgetID, err)
			return ErrInternal
		}
		if _, err = txRepo.Account.AdjustNetBalance(currentUserID, budget.AccountID, budget.Budget); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			log.Printf("failed to update project budget: %v", err)
			return ErrInternal
		}
		if err = txRepo.ProjectBudget.Delete(budgetID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			log.Printf("failed to delete project budget %d: %v", budgetID, err)
			return ErrInternal
		}
		return nil
	})

}
