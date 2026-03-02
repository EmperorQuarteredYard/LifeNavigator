package service

import (
	"LifeNavigator/backend/internal/models"
	"LifeNavigator/backend/internal/repository"
	"LifeNavigator/backend/pkg/jwt"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *models.User) (err error)
	Login(userName, password string) (user *models.User, err error)
	FindById(id uint64) (user *models.User, err error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	HardDeleteById(id uint64) (err error)
	SoftDeleteById(id uint64) (err error)
	RefreshToken(refreshToken string) (accessToken string, newRefreshToken string, err error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
func (s *userService) RefreshToken(refreshToken string) (accessToken string, newRefreshToken string, err error) {
	// 验证 refresh token
	claims, err := jwt.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	// 检查用户是否存在（防止用户已被删除）
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", ErrUserNotFound
		}
		return "", "", ErrInternal
	}

	accessToken, newRefreshToken, err = jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", ErrInternal
	}
	return accessToken, newRefreshToken, nil
}

func (s *userService) SoftDeleteById(id uint64) (err error) {
	err = s.userRepo.SoftDeleteById(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
	}
	return
}

func (s *userService) HardDeleteById(id uint64) (err error) {
	err = s.userRepo.HardDeleteById(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		log.Printf("failed to hard delete user : %v", err)
		return ErrInternal
	}
	return
}

func (s *userService) Register(user *models.User) (err error) {
	if user.Password == "" || user.Username == "" {
		return ErrUserInfoNotFound
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.Role = "user"
	err = s.userRepo.Create(user)
	if err != nil {
		if errors.Is(err, repository.ErrUserNameExists) {
			return ErrUserNameExists
		}
		log.Printf("failed to register user : %v", err)
		return ErrInternal
	}
	return nil
}

func (s *userService) Login(userName, password string) (user *models.User, err error) {
	user, err = s.userRepo.GetByUsername(userName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("failed to login user  %s: %v", userName, err)
		return nil, ErrInternal
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrPasswordWrong
	}
	return user, nil
}

func (s *userService) FindById(id uint64) (user *models.User, err error) {
	user, err = s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("failed to get user by id %d: %v", id, err)
		return nil, ErrInternal
	}
	return user, nil
}

func (s *userService) GetByUsername(username string) (user *models.User, err error) {
	user, err = s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("failed to get user by username %s: %v", username, err)
		return nil, ErrInternal
	}
	return user, nil
}
func (s *userService) GetByEmail(email string) (user *models.User, err error) {
	user, err = s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("failed to get user by email %s: %v", email, err)
		return nil, ErrInternal
	}
	return user, nil
}
