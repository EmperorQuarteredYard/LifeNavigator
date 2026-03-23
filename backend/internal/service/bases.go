package service

import (
	"LifeNavigator/internal/interfaces/Repository"
	"LifeNavigator/internal/interfaces/Service"
	"LifeNavigator/pkg/permission"
	"log"
)

type projectBase struct {
	projectRepo Repository.ProjectRepository
}

func (s *projectBase) checkProjectOwnership(projectID uint64, userID uint64) error {
	owned, err := s.projectRepo.CheckOwnership(userID, projectID)
	if err != nil {
		log.Printf("failed to check project ownership %d: %v", projectID, err)
		return Service.ErrInternal
	}
	if !owned {
		return Service.ErrForbidden
	}
	return nil
}

func (s *projectBase) checkProjectAccessibility(userID, projectID uint64, operation permission.Operation) error {
	result, err := s.projectRepo.CheckAccessibility(userID, projectID, operation, userID == 0)
	if err != nil {
		log.Printf("failed to check project accessibility %d: %v", projectID, err)
		return Service.ErrInternal
	}
	if !result {
		return Service.ErrForbidden
	}
	return nil
}
