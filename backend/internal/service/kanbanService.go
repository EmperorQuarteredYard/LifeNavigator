package service

import (
	"LifeNavigator/internal/interfaces/repositoryInte"
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"LifeNavigator/pkg/dto"
	"LifeNavigator/pkg/kanbanStatus"
	"LifeNavigator/pkg/permission"
	"errors"
	"log"
)

type KanbanService interface {
	Create(userID uint64, req *dto.CreateKanbanRequest) (*dto.KanbanResponse, error)
	GetByID(userID, id uint64) (*dto.KanbanResponse, error)
	ListByUserID(userID uint64) (*dto.KanbanListResponse, error)
	Update(userID uint64, id uint64, req *dto.UpdateKanbanRequest) error
	Delete(userID, id uint64) error
	GetKanbanTasks(userID, kanbanID uint64, status *uint8, page, pageSize int) (*dto.KanbanTaskListResponse, error)
	GetDefaultKanbanTasks(userID uint64, status *uint8, page, pageSize int) (*dto.KanbanTaskListResponse, error)
	SetDefaultKanban(userID, kanbanID uint64) error
	AddProjectInCheck(userID, projectID, kanbanID uint64) error
}

type kanbanService struct {
	kanbanRepo  repository.KanbanRepository
	taskRepo    repository.TaskRepository
	taskPayRepo repository.TaskBudgetRepository
	*projectBase
}

func (s *kanbanService) AddProjectInCheck(userID, projectID, kanbanID uint64) error {
	//TODO implement me
	panic("implement me")
}

func NewKanbanService(
	kanbanRepo repository.KanbanRepository,
	projectRepo repository.ProjectRepository,
	taskRepo repository.TaskRepository,
	taskPayRepo repository.TaskBudgetRepository,
) KanbanService {
	return &kanbanService{
		kanbanRepo:  kanbanRepo,
		projectBase: &projectBase{projectRepo: projectRepo},
		taskRepo:    taskRepo,
		taskPayRepo: taskPayRepo,
	}
}

func (s *kanbanService) checkKanbanOwnership(userID, kanbanID uint64) error {
	owned, err := s.kanbanRepo.CheckOwnership(userID, kanbanID)
	if err != nil {
		log.Printf("failed to check kanban ownership %d: %v", kanbanID, err)
		return ErrInternal
	}
	if !owned {
		return ErrForbidden
	}
	return nil
}

func (s *kanbanService) filterOwnedProjects(userID uint64, projectIDs []uint64) ([]uint64, error) {
	var ownedIDs []uint64
	for _, pid := range projectIDs {
		err := s.checkProjectAccessibility(pid, userID, permission.OpRead)
		if err != nil {
			if errors.Is(err, ErrForbidden) {
				continue
			}
			return nil, err
		}
	}
	return ownedIDs, nil
}

func (s *kanbanService) toKanbanResponse(kanban *models.Kanban, projectIDs []uint64) *dto.KanbanResponse {
	projects := make([]dto.KanbanProject, len(projectIDs))
	for i, pid := range projectIDs {
		project, err := s.projectRepo.GetByID(pid)
		if err == nil {
			projects[i] = dto.KanbanProject{
				ID:          project.ID,
				Name:        project.Name,
				Description: project.Description,
			}
		}
	}
	return &dto.KanbanResponse{
		ID:          kanban.ID,
		UserID:      kanban.UserID,
		Name:        kanban.Name,
		Description: kanban.Description,
		IsDefault:   kanban.IsDefault,
		SortOrder:   kanban.SortOrder,
		Projects:    projects,
		CreatedAt:   kanban.CreatedAt,
		UpdatedAt:   kanban.UpdatedAt,
	}
}

func (s *kanbanService) Create(userID uint64, req *dto.CreateKanbanRequest) (*dto.KanbanResponse, error) {
	ownedProjectIDs, err := s.filterOwnedProjects(userID, req.ProjectIDs)
	if err != nil {
		return nil, err
	}

	kanban := &models.Kanban{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.kanbanRepo.Create(kanban); err != nil {
		log.Printf("failed to create kanban: %v", err)
		return nil, ErrInternal
	}

	if len(ownedProjectIDs) > 0 {
		if err := s.kanbanRepo.SetProjects(kanban.ID, ownedProjectIDs); err != nil {
			log.Printf("failed to set kanban projects: %v", err)
			return nil, ErrInternal
		}
	}

	return s.toKanbanResponse(kanban, ownedProjectIDs), nil
}

func (s *kanbanService) GetByID(userID, id uint64) (*dto.KanbanResponse, error) {
	if err := s.checkKanbanOwnership(userID, id); err != nil {
		return nil, err
	}

	kanban, err := s.kanbanRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repositoryInte.ErrNotFound) {
			return nil, ErrKanbanNotFound
		}
		log.Printf("failed to get kanban %d: %v", id, err)
		return nil, ErrInternal
	}

	projectIDs, err := s.kanbanRepo.GetProjectIDs(id)
	if err != nil {
		log.Printf("failed to get kanban project IDs: %v", err)
		return nil, ErrInternal
	}
	projectIDs, err = s.filterOwnedProjects(userID, projectIDs)
	if err != nil {
		return nil, err
	}

	return s.toKanbanResponse(kanban, projectIDs), nil
}

func (s *kanbanService) ListByUserID(userID uint64) (*dto.KanbanListResponse, error) {
	kanbans, err := s.kanbanRepo.ListByUserID(userID)
	if err != nil {
		log.Printf("failed to list kanbans for user %d: %v", userID, err)
		return nil, ErrInternal
	}

	list := make([]dto.KanbanResponse, len(kanbans))
	for i, k := range kanbans {
		projectIDs, err := s.kanbanRepo.GetProjectIDs(k.ID)
		if err != nil {
			log.Printf("failed to get kanban project IDs: %v", err)
			projectIDs = []uint64{}
		}
		resp := s.toKanbanResponse(&k, projectIDs)
		list[i] = *resp
	}

	return &dto.KanbanListResponse{
		Total: int64(len(list)),
		List:  list,
	}, nil
}

func (s *kanbanService) Update(userID uint64, id uint64, req *dto.UpdateKanbanRequest) error {
	if err := s.checkKanbanOwnership(userID, id); err != nil {
		return err
	}

	kanban, err := s.kanbanRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repositoryInte.ErrNotFound) {
			return ErrKanbanNotFound
		}
		return ErrInternal
	}

	if req.Name != "" {
		kanban.Name = req.Name
	}
	if req.Description != "" {
		kanban.Description = req.Description
	}
	kanban.SortOrder = req.SortOrder

	if err := s.kanbanRepo.Update(kanban); err != nil {
		log.Printf("failed to update kanban %d: %v", id, err)
		return ErrInternal
	}

	if req.ProjectIDs != nil {
		ownedProjectIDs, err := s.filterOwnedProjects(userID, req.ProjectIDs)
		if err != nil {
			return err
		}
		if err := s.kanbanRepo.SetProjects(id, ownedProjectIDs); err != nil {
			log.Printf("failed to set kanban projects: %v", err)
			return ErrInternal
		}
	}

	return nil
}

func (s *kanbanService) Delete(userID, id uint64) error {
	if err := s.checkKanbanOwnership(userID, id); err != nil {
		return err
	}

	if err := s.kanbanRepo.Delete(id); err != nil {
		if errors.Is(err, repositoryInte.ErrNotFound) {
			return ErrKanbanNotFound
		}
		log.Printf("failed to delete kanban %d: %v", id, err)
		return ErrInternal
	}
	return nil
}

func (s *kanbanService) GetKanbanTasks(userID, kanbanID uint64, status *uint8, page, pageSize int) (*dto.KanbanTaskListResponse, error) {
	if err := s.checkKanbanOwnership(userID, kanbanID); err != nil {
		return nil, err
	}

	projectIDs, err := s.kanbanRepo.GetProjectIDs(kanbanID)
	if err != nil {
		log.Printf("failed to get kanban project IDs: %v", err)
		return nil, ErrInternal
	}

	ownedProjectIDs, err := s.filterOwnedProjects(userID, projectIDs)
	if err != nil {
		return nil, err
	}

	if len(ownedProjectIDs) == 0 {
		return &dto.KanbanTaskListResponse{Total: 0, List: []*dto.KanbanTaskResponse{}}, nil
	}

	return s.getTasksFromProjects(userID, ownedProjectIDs, status, page, pageSize)
}

func (s *kanbanService) GetDefaultKanbanTasks(userID uint64, status *uint8, page, pageSize int) (*dto.KanbanTaskListResponse, error) {
	kanban, err := s.kanbanRepo.GetDefaultKanban(userID)
	if err != nil {
		if errors.Is(err, repositoryInte.ErrNotFound) {
			kanbans, listErr := s.kanbanRepo.ListByUserID(userID)
			if listErr != nil || len(kanbans) == 0 {
				return &dto.KanbanTaskListResponse{Total: 0, List: []*dto.KanbanTaskResponse{}}, nil
			}
			kanban = &kanbans[0]
		} else {
			log.Printf("failed to get default kanban: %v", err)
			return nil, ErrInternal
		}
	}

	return s.GetKanbanTasks(userID, kanban.ID, status, page, pageSize)
}

func (s *kanbanService) getTasksFromProjects(userID uint64, projectIDs []uint64, status *uint8, page, pageSize int) (*dto.KanbanTaskListResponse, error) {
	var allTasks []models.Task
	var total int64

	for _, pid := range projectIDs {
		var tasks []models.Task
		var count int64
		var err error

		if status != nil {
			tasks, err = s.taskRepo.GetByStatus(pid, *status)
			count = int64(len(tasks))
		} else {
			tasks, count, err = s.taskRepo.ListByProjectID(pid, page, pageSize)
		}

		if err != nil {
			log.Printf("failed to get tasks for project %d: %v", pid, err)
			continue
		}
		allTasks = append(allTasks, tasks...)
		total += count
	}

	projectMap := make(map[uint64]string)
	for _, pid := range projectIDs {
		project, err := s.projectRepo.GetByID(pid)
		if err == nil {
			projectMap[pid] = project.Name
		}
	}

	taskIDs := make([]uint64, len(allTasks))
	for i, t := range allTasks {
		taskIDs[i] = t.ID
	}

	allPrereqs := make(map[uint64][]models.TaskDependency)
	for _, tid := range taskIDs {
		prereqs, err := s.taskRepo.GetPrerequisites(tid)
		if err == nil {
			allPrereqs[tid] = prereqs
		}
	}

	list := make([]*dto.KanbanTaskResponse, len(allTasks))
	for i, task := range allTasks {
		statusText, _ := kanbanStatus.GetStatusName(int(task.Status))

		payments, _ := s.taskPayRepo.GetByTaskID(task.ID)
		paymentDtos := make([]*dto.TaskPaymentResponse, len(payments))
		for j, p := range payments {
			paymentDtos[j] = &dto.TaskPaymentResponse{
				ID:       p.ID,
				TaskID:   p.TaskID,
				BudgetID: p.BudgetID,
				Amount:   p.Amount,
			}
		}

		var prereqInfos []*dto.TaskDependencyInfo
		if prereqs, ok := allPrereqs[task.ID]; ok {
			prereqInfos = make([]*dto.TaskDependencyInfo, len(prereqs))
			for j, prereq := range prereqs {
				prereqTask, err := s.taskRepo.GetByID(prereq.PrerequisiteID)
				if err == nil {
					prereqStatusText, _ := kanbanStatus.GetStatusName(int(prereqTask.Status))
					prereqInfos[j] = &dto.TaskDependencyInfo{
						TaskID:         prereq.PrerequisiteID,
						TaskName:       prereqTask.Name,
						TaskStatus:     prereqTask.Status,
						TaskStatusText: prereqStatusText,
					}
				}
			}
		}

		list[i] = &dto.KanbanTaskResponse{
			ID:            task.ID,
			ProjectID:     task.ProjectID,
			ProjectName:   projectMap[task.ProjectID],
			Name:          task.Name,
			Description:   task.Description,
			Type:          task.Type,
			Status:        task.Status,
			StatusText:    statusText,
			Category:      task.Category,
			Deadline:      task.Deadline,
			CompletedAt:   task.CompletedAt,
			CreatedAt:     task.CreatedAt,
			UpdatedAt:     task.UpdatedAt,
			Prerequisites: prereqInfos,
			Payments:      paymentDtos,
		}
	}

	return &dto.KanbanTaskListResponse{
		Total: total,
		List:  list,
	}, nil
}

func (s *kanbanService) SetDefaultKanban(userID, kanbanID uint64) error {
	if err := s.checkKanbanOwnership(userID, kanbanID); err != nil {
		return err
	}

	if err := s.kanbanRepo.SetDefaultKanban(userID, kanbanID); err != nil {
		if errors.Is(err, repositoryInte.ErrNotFound) {
			return ErrKanbanNotFound
		}
		log.Printf("failed to set default kanban: %v", err)
		return ErrInternal
	}
	return nil
}
