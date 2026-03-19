package service

import (
	"LifeNavigator/internal/models"
	"LifeNavigator/internal/repository"
	"LifeNavigator/middleWare/jwt"
	"LifeNavigator/pkg/roles"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *models.User) error                                                  // 注册无需权限
	Login(username, password string) (*models.User, error)                             // 登录无需权限
	GetByID(id uint64, currentUserID uint64) (*models.User, error)                     // 查询用户信息：只有自己或管理员可查（这里简化：仅自己可查）
	GetByUsername(username string, currentUserID uint64) (*models.User, error)         // 若允许用户查看他人用户名信息则另当别论，这里仅允许查看自己
	GetByEmail(email string, currentUserID uint64) (*models.User, error)               // 若允许用户查看他人用户名信息则另当别论，这里仅允许查看自己
	Update(user *models.User, currentUserID uint64) error                              // 更新用户信息：只有自己可更新
	HardDeleteByID(id uint64, currentUserID uint64) error                              // 硬删除：仅自己可执行（或管理员，这里简化）
	SoftDeleteByID(id uint64, currentUserID uint64) error                              // 软删除：仅自己可执行
	RefreshToken(refreshToken string) (accessToken, newRefreshToken string, err error) // 刷新令牌无需用户 ID
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(user *models.User) error {
	if user.Role == "" {
		user.Role = roles.User
	}
	if user.Password == "" || user.Username == "" {
		return ErrUserInfoIncomplete
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return ErrInternal
	}
	user.Password = string(hashed)
	user.Role = "user" //TODO 这里之后要做成允许注册为管理员、开发者等
	if err := s.userRepo.Create(user); err != nil {
		if errors.Is(err, repository.ErrRecordExist) {
			return ErrUserNameExists
		}
		log.Printf("failed to create user: %v", err)
		return ErrInternal
	}
	return nil
}

func (s *userService) Login(username, password string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	log.Println("用户登录:" + username + "/" + password)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("failed to get user by username: %v", err)
		return nil, ErrInternal
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrPasswordWrong
	}
	return user, nil
}

func (s *userService) GetByID(id uint64, currentUserID uint64) (*models.User, error) {
	// 只允许用户查询自己
	if id != currentUserID {
		return nil, ErrForbidden
	}
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("failed to get user by id %d: %v", id, err)
		return nil, ErrInternal
	}
	return user, nil
}

func (s *userService) GetByUsername(username string, currentUserID uint64) (*models.User, error) {
	// 假设只允许用户查询自己的用户名信息，但用户名是登录标识，通常允许公开？这里保守只允许查自己
	// 更合理：先查出用户，再比对 ID
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("failed to get user by username: %v", err)
		return nil, ErrInternal
	}
	if user.ID != currentUserID {
		return nil, ErrForbidden
	}
	return user, nil
}

func (s *userService) GetByEmail(email string, currentUserID uint64) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		log.Printf("failed to get user by email: %v", err)
		return nil, ErrInternal
	}
	if user.ID != currentUserID {
		return nil, ErrForbidden
	}
	return user, nil
}

func (s *userService) Update(user *models.User, currentUserID uint64) error {
	// 只能更新自己的信息
	if user.ID != currentUserID {
		return ErrForbidden
	}
	// 可选：不允许修改某些字段，如密码需单独处理，这里简化
	if err := s.userRepo.Update(user); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		log.Printf("failed to update user %d: %v", user.ID, err)
		return ErrInternal
	}
	return nil
}

func (s *userService) HardDeleteByID(id uint64, currentUserID uint64) error {
	if id != currentUserID {
		return ErrForbidden
	}
	if err := s.userRepo.HardDeleteById(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		log.Printf("failed to hard delete user %d: %v", id, err)
		return ErrInternal
	}
	return nil
}

func (s *userService) SoftDeleteByID(id uint64, currentUserID uint64) error {
	if id != currentUserID {
		return ErrForbidden
	}
	if err := s.userRepo.SoftDeleteById(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrUserNotFound
		}
		log.Printf("failed to soft delete user %d: %v", id, err)
		return ErrInternal
	}
	return nil
}

func (s *userService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := jwt.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", "", ErrInvalidToken
	}
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", ErrUserNotFound
		}
		log.Printf("failed to get user for refresh: %v", err)
		return "", "", ErrInternal
	}
	accessToken, newRefreshToken, err := jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		return "", "", ErrInternal
	}
	return accessToken, newRefreshToken, nil
}
