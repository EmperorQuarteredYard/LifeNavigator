package service

import (
	"LifeNavigator/internal/interfaces/repositoryInte"
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/permission"
	"context"
	"errors"
	"log"
)

type ProjectService interface {
	Create(userID uint64, project *models.Project) (*dto.ProjectResponse, error)
	GetByID(userID, id uint64) (*dto.ProjectResponse, error)
	ListByUserID(userID uint64, offset, limit int) (*dto.ProjectListResponse, error)
	Update(userID uint64, project *models.Project) error
	Delete(userID, id uint64) error
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
		projectBase:       &projectBase{projectRepo: projectRepo},
		projectBudgetRepo: projectBudgetRepo,
		taskBudgetRepo:    taskBudgetRepo,
		taskRepo:          taskRepo,
	}
}

type projectService struct {
	transactor        repository.Transactor
	projectBudgetRepo repository.ProjectBudgetRepository
	taskBudgetRepo    repository.TaskBudgetRepository
	taskRepo          repository.TaskRepository
	*projectBase
}

func (s *projectService) Create(userID uint64, project *models.Project) (*dto.ProjectResponse, error) {
	if userID == 0 {
		return nil, ErrForbidden
	}
	project.Owner = userID
	err := s.projectRepo.Create(project, []uint64{userID})
	if err != nil {
		log.Printf("failed to create project: %v", err)
		return nil, ErrInternal
	}
	// 刷新调度交给 budgetService，这里不再操作调度器
	return &dto.ProjectResponse{
		ID:              project.ID,
		Name:            project.Name,
		OwnerID:         userID,
		Description:     project.Description,
		RefreshInterval: project.RefreshInterval,
		LastRefresh:     project.LastRefresh,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
		Permission:      project.Permission.String(),
	}, nil
}

func (s *projectService) GetByID(userID, id uint64) (*dto.ProjectResponse, error) {
	if err := s.checkProjectAccessibility(userID, id, permission.OpRead); err != nil {
		return nil, err
	}
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repositoryInte.ErrNotFound) {
			return nil, ErrProjectNotFound
		}
		log.Printf("failed to get project %d: %v", id, err)
		return nil, ErrInternal
	}

	budgets, _ := s.projectBudgetRepo.GetByProjectID(id) //TODO 这里不应该返回的(虽然这是我最近加发hh）
	dtoBudgets := make([]*dto.ProjectBudgetResponse, len(budgets))
	for i, b := range budgets {
		dtoBudgets[i] = &dto.ProjectBudgetResponse{
			ID:        b.ID,
			ProjectID: b.ProjectID,
			AccountID: b.AccountID,
			Type:      b.Unit,
			Budget:    b.Budget,
			Used:      b.Used,
		}
	}

	return &dto.ProjectResponse{
		ID:              project.ID,
		Name:            project.Name,
		Description:     project.Description,
		RefreshInterval: project.RefreshInterval,
		LastRefresh:     project.LastRefresh,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
		Budgets:         dtoBudgets,
		Permission:      project.Permission.String(),
	}, nil
}

func (s *projectService) ListByUserID(userID uint64, offset, limit int) (*dto.ProjectListResponse, error) {
	projects, err := s.projectRepo.ListByUserID(userID, offset, limit)
	if err != nil {
		log.Printf("failed to list projects for user %d: %v", userID, err)
		return nil, ErrInternal
	}

	items := make([]*dto.ProjectResponse, len(projects))
	for i, p := range projects {
		if p.Owner == userID {
			if !p.Permission.Has(permission.RoleOwner, permission.OpRead) {
				continue
			}
		} else {
			if !p.Permission.Has(permission.RoleWorkmate, permission.OpRead) {
				continue
			}
		}
		items[i] = &dto.ProjectResponse{
			ID:              p.ID,
			Name:            p.Name,
			Description:     p.Description,
			RefreshInterval: p.RefreshInterval,
			LastRefresh:     p.LastRefresh,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
			Permission:      p.Permission.String(),
		}
	}
	return &dto.ProjectListResponse{Items: items, Total: int64(len(items))}, nil
}

func (s *projectService) Update(userID uint64, project *models.Project) error {
	if err := s.checkProjectOwnership(project.ID, userID); err != nil {
		return err
	}
	if err := s.projectRepo.Update(project); err != nil {
		if errors.Is(err, repositoryInte.ErrNotFound) {
			return ErrProjectNotFound
		}
		log.Printf("failed to update project %d: %v", project.ID, err)
		return ErrInternal
	}
	// 调度更新已移至 budgetService，此处不再处理
	//todo 禁止清理自己的访问权
	return nil
}

func (s *projectService) Delete(userID, id uint64) error {
	if err := s.checkProjectOwnership(id, userID); err != nil {
		return err
	}

	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		budgets, err := txRepo.ProjectBudget.GetByProjectID(id) //TODO 鉴于Budget确实是依附于Project存在的，就不解耦了
		if err != nil {
			return ErrInternal
		}
		for _, budget := range budgets {
			if err = txRepo.ProjectBudget.Delete(budget.ID); err != nil {
				log.Printf("failed to delete project budget %d: %v", budget.ID, err)
				return ErrInternal
			}
		}

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

		if err := txRepo.Project.Delete(id); err != nil {
			if errors.Is(err, repositoryInte.ErrNotFound) {
				return ErrProjectNotFound
			}
			log.Printf("failed to delete project %d: %v", id, err)
			return ErrInternal
		}
		// TODO 删除调度由 budgetService 处理，但此处无直接调用，暂不处理
		return nil
	})
}
