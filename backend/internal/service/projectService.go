package service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/refresh"
	"LifeNavigator/pkg/scheduler"
	"context"
	"errors"
	"log"
	"time"
)

// todo 任务扣除预算的行为应当放在任务DeadLine所在刷新周期内

type ProjectService interface {
	Create(userID uint64, project *models.Project) (*dto.ProjectResponse, error)
	GetByID(userID, id uint64) (*dto.ProjectResponse, error)
	ListByUserID(userID uint64, offset, limit int) (*dto.ProjectListResponse, error)
	Update(userID uint64, project *models.Project) error
	Delete(userID, id uint64) error
	AddBudget(userID, projectID uint64, budget *models.ProjectBudget) error
	UpdateBudget(userID uint64, budget *models.ProjectBudget) error
	DeleteBudget(userID, budgetID uint64) error
	RefreshBudget(projectID uint64) error
	StartAutoRefresh() error
	EndAutoRefresh()
}

func NewProjectService(
	transactor repository.Transactor,
	projectRepo repository.ProjectRepository,
	projectBudgetRepo repository.ProjectBudgetRepository,
	taskBudgetRepo repository.TaskBudgetRepository,
	taskRepo repository.TaskRepository,
) ProjectService {
	scheduleServ := scheduler.NewScheduleService(func(uint64) error { return nil }, func(int, int) (int64, []*scheduler.Schedule) { return 0, nil })
	return &projectService{
		transactor:        transactor,
		projectRepo:       projectRepo,
		projectBudgetRepo: projectBudgetRepo,
		taskBudgetRepo:    taskBudgetRepo,
		taskRepo:          taskRepo,
		scheduleService:   scheduleServ,
	}
}

type projectService struct {
	transactor        repository.Transactor
	projectRepo       repository.ProjectRepository
	projectBudgetRepo repository.ProjectBudgetRepository
	taskBudgetRepo    repository.TaskBudgetRepository
	taskRepo          repository.TaskRepository
	scheduleService   *scheduler.ScheduleService
}

func (s *projectService) StartAutoRefresh() error {
	var schExeFunc = func(ID uint64) error {
		return s.RefreshBudget(ID)
	}
	var schGetFunc = func(page, pageSize int) (int64, []*scheduler.Schedule) {
		sch, totalZ, errZ := s.projectRepo.GetRefreshInformation(page, pageSize)
		if errZ != nil {
			log.Printf("failed to start auto-refresh:%v", errZ)
		}
		return totalZ, sch
	}
	s.scheduleService = scheduler.NewScheduleService(schExeFunc, schGetFunc)
	s.scheduleService.Start()
	return nil
}
func (s *projectService) EndAutoRefresh() {
	s.scheduleService.Stop()
}

func (s *projectService) RefreshBudget(projectID uint64) error {
	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		project, err := s.projectRepo.GetByID(projectID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotFound):
				return ErrProjectNotFound
			default:
				log.Printf("fail to get project %d :%v", projectID, err)
				return ErrInternal
			}
		}
		project.LastRefresh = time.Now()
		if err = txRepo.Project.Update(project); err != nil {
			switch {
			case errors.Is(err, repository.ErrNotFound):
				return ErrProjectNotFound
			default:
				log.Printf("fail to update project %d :%v", projectID, err)
				return ErrInternal
			}
		}
		budgets, err := txRepo.ProjectBudget.GetByProjectID(projectID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotFound):
				return ErrBudgetNotFound
			}
		}
		for _, budget := range budgets {
			if _, err = txRepo.Account.AdjustNetBalance(budget.AccountID, budget.Used); err != nil {
				switch {
				case errors.Is(err, repository.ErrNotFound):
					return ErrAccountNotFound
				default:
					log.Printf("fail to add budget %d :%v", budget.ID, err)
					return ErrInternal
				}
			}
		}
		if _, strategy := refresh.ReduceRefreshInterval(project.RefreshInterval); strategy != refresh.TypeRefreshNever {
			s.scheduleService.AddRefreshSchedule(projectID, *refresh.GetNextRefreshTime(project.RefreshInterval, project.LastRefresh))
		}

		return nil
	})
}

func (s *projectService) checkProjectOwnership(projectID uint64, userID uint64) error {
	owned, err := s.projectRepo.CheckOwnership(userID, projectID)
	if err != nil {
		log.Printf("failed to check project ownership %d: %v", projectID, err)
		return ErrInternal
	}
	if !owned {
		return ErrForbidden
	}
	return nil
}

func (s *projectService) Create(userID uint64, project *models.Project) (*dto.ProjectResponse, error) {
	err := s.projectRepo.Create(project, []uint64{userID})
	if err != nil {
		log.Printf("failed to create project: %v", err)
		return nil, ErrInternal
	}
	if _, strategy := refresh.ReduceRefreshInterval(project.RefreshInterval); strategy != refresh.TypeRefreshNever {
		s.scheduleService.AddRefreshSchedule(project.ID, *refresh.GetNextRefreshTime(project.RefreshInterval, project.LastRefresh))
	}
	return &dto.ProjectResponse{
		ID:              project.ID,
		Name:            project.Name,
		Description:     project.Description,
		RefreshInterval: project.RefreshInterval,
		LastRefresh:     project.LastRefresh,
		MaxTaskID:       project.MaxTaskID,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
	}, nil
}

func (s *projectService) GetByID(userID, id uint64) (*dto.ProjectResponse, error) {
	if err := s.checkProjectOwnership(id, userID); err != nil {
		return nil, err
	}

	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrProjectNotFound
		}
		log.Printf("failed to get project %d: %v", id, err)
		return nil, ErrInternal
	}

	budgets, _ := s.projectBudgetRepo.GetByProjectID(id)
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
		MaxTaskID:       project.MaxTaskID,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
		Budgets:         dtoBudgets,
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
		items[i] = &dto.ProjectResponse{
			ID:              p.ID,
			Name:            p.Name,
			Description:     p.Description,
			RefreshInterval: p.RefreshInterval,
			LastRefresh:     p.LastRefresh,
			MaxTaskID:       p.MaxTaskID,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
		}
	}
	return &dto.ProjectListResponse{Items: items, Total: int64(len(items))}, nil
}

func (s *projectService) Update(userID uint64, project *models.Project) error {
	if err := s.checkProjectOwnership(project.ID, userID); err != nil {
		return err
	}
	if err := s.projectRepo.Update(project); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrProjectNotFound
		}
		log.Printf("failed to update project %d: %v", project.ID, err)
		return ErrInternal
	}
	if _, strategy := refresh.ReduceRefreshInterval(project.RefreshInterval); strategy != refresh.TypeRefreshNever {
		s.scheduleService.UpdateRefreshSchedule(project.ID, *refresh.GetNextRefreshTime(project.RefreshInterval, project.LastRefresh))
	}
	return nil
}

func (s *projectService) Delete(userID, id uint64) error {
	if err := s.checkProjectOwnership(id, userID); err != nil {
		return err
	}

	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		budgets, err := txRepo.ProjectBudget.GetByProjectID(id)
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
			if errors.Is(err, repository.ErrNotFound) {
				return ErrProjectNotFound
			}
			log.Printf("failed to delete project %d: %v", id, err)
			return ErrInternal
		}
		s.scheduleService.DeleteRefreshSchedule(id)
		return nil
	})
}

func (s *projectService) AddBudget(userID, projectID uint64, budget *models.ProjectBudget) error {
	if err := s.checkProjectOwnership(projectID, userID); err != nil {
		return err
	}

	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		budget.ProjectID = projectID
		if err := txRepo.ProjectBudget.Create(budget); err != nil {
			log.Printf("failed to create project budget: %v", err)
			return ErrInternal
		}
		if budget.AccountID != 0 {
			if _, err := txRepo.Account.AdjustNetBalance(budget.AccountID, -1*budget.Budget); err != nil {
				log.Printf("failed to adjust net balance: %v", err)
				return ErrInternal
			}
		}
		return nil
	})
}

func (s *projectService) UpdateBudget(userID uint64, budget *models.ProjectBudget) error {
	existing, err := s.projectBudgetRepo.GetByID(budget.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to get budget %d: %v", budget.ID, err)
		return ErrInternal
	}

	if err := s.checkProjectOwnership(existing.ProjectID, userID); err != nil {
		return err
	}

	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		oldBudget, err := txRepo.ProjectBudget.GetByID(budget.ID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}

		if budget.AccountID != oldBudget.AccountID {
			if oldBudget.AccountID != 0 {
				if _, err := txRepo.Account.AdjustNetBalance(oldBudget.AccountID, oldBudget.Budget); err != nil {
					return ErrInternal
				}
			}
			if budget.AccountID != 0 {
				if _, err := txRepo.Account.AdjustNetBalance(budget.AccountID, -1*budget.Budget); err != nil {
					return ErrInternal
				}
			}
		} else if budget.Budget != oldBudget.Budget && budget.AccountID != 0 {
			if _, err := txRepo.Account.AdjustNetBalance(oldBudget.AccountID, oldBudget.Budget-budget.Budget); err != nil {
				return ErrInternal
			}
		}

		if err := txRepo.ProjectBudget.Update(budget); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}
		return nil
	})
}

func (s *projectService) DeleteBudget(userID, budgetID uint64) error {
	existing, err := s.projectBudgetRepo.GetByID(budgetID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrBudgetNotFound
		}
		log.Printf("failed to get budget %d: %v", budgetID, err)
		return ErrInternal
	}

	if err := s.checkProjectOwnership(existing.ProjectID, userID); err != nil {
		return err
	}

	return s.transactor.WithinTransaction(context.Background(), func(txRepo repository.TxRepositories) error {
		budget, err := txRepo.ProjectBudget.GetByID(budgetID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}

		if budget.AccountID != 0 {
			if _, err := txRepo.Account.AdjustNetBalance(budget.AccountID, budget.Budget); err != nil {
				return ErrInternal
			}
		}

		if err := txRepo.ProjectBudget.Delete(budgetID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return ErrBudgetNotFound
			}
			return ErrInternal
		}
		return nil
	})
}
