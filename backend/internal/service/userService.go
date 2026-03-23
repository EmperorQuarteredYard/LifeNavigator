package service

import (
	"LifeNavigator/internal/interfaces/Repository"
	"LifeNavigator/internal/interfaces/Service"
	"LifeNavigator/internal/models"
	"LifeNavigator/middleWare/jwt"
	"LifeNavigator/pkg/roles"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo Repository.UserRepository
}

func NewUserService(userRepo Repository.UserRepository) Service.UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) UpdateAvatar(userID uint64, avatarURL string) error {
	// 检查用户是否存在（可选，直接调用 repo 即可）
	if err := s.userRepo.UpdateAvatar(userID, avatarURL); err != nil {
		if errors.Is(err, Repository.ErrNotFound) {
			return Service.ErrUserNotFound
		}
		log.Printf("failed to update avatar for user %d: %v", userID, err)
		return Service.ErrInternal
	}
	return nil
}

func (s *userService) Register(user *models.User) error {
	if user.Role == "" {
		user.Role = roles.User
	}
	if user.Password == "" || user.Username == "" {
		return Service.ErrUserInfoIncomplete
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return Service.ErrInternal
	}
	user.Password = string(hashed)
	user.Role = "user" //TODO 这里之后要做成允许注册为管理员、开发者等
	if err := s.userRepo.Create(user); err != nil {
		if errors.Is(err, Repository.ErrRecordExist) {
			return Service.ErrUserNameExists
		}
		log.Printf("failed to create user: %v", err)
		return Service.ErrInternal
	}
	return nil
}

func (s *userService) Login(username, password string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	log.Println("用户登录:" + username + "/" + password)
	if err != nil {
		if errors.Is(err, Repository.ErrNotFound) {
			return nil, Service.ErrUserNotFound
		}
		log.Printf("failed to get user by username: %v", err)
		return nil, Service.ErrInternal
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, Service.ErrPasswordWrong
	}
	return user, nil
}

func (s *userService) GetByID(id uint64, currentUserID uint64) (*models.User, error) {
	// 只允许用户查询自己
	if id != currentUserID {
		return nil, Service.ErrForbidden
	}
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, Repository.ErrNotFound) {
			return nil, Service.ErrUserNotFound
		}
		log.Printf("failed to get user by id %d: %v", id, err)
		return nil, Service.ErrInternal
	}
	return user, nil
}

func (s *userService) GetByUsername(username string, currentUserID uint64) (*models.User, error) {
	// 假设只允许用户查询自己的用户名信息，但用户名是登录标识，通常允许公开？这里保守只允许查自己
	// 更合理：先查出用户，再比对 ID
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, Repository.ErrNotFound) {
			return nil, Service.ErrUserNotFound
		}
		log.Printf("failed to get user by username: %v", err)
		return nil, Service.ErrInternal
	}
	if user.ID != currentUserID {
		return nil, Service.ErrForbidden
	}
	return user, nil
}

func (s *userService) GetByEmail(email string, currentUserID uint64) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		if errors.Is(err, Repository.ErrNotFound) {
			return nil, Service.ErrUserNotFound
		}
		log.Printf("failed to get user by email: %v", err)
		return nil, Service.ErrInternal
	}
	if user.ID != currentUserID {
		return nil, Service.ErrForbidden
	}
	return user, nil
}

func (s *userService) Update(user *models.User, currentUserID uint64) error {
	// 只能更新自己的信息
	if user.ID != currentUserID {
		return Service.ErrForbidden
	}
	// 可选：不允许修改某些字段，如密码需单独处理，这里简化
	if err := s.userRepo.Update(user); err != nil {
		if errors.Is(err, Repository.ErrNotFound) {
			return Service.ErrUserNotFound
		}
		log.Printf("failed to update user %d: %v", user.ID, err)
		return Service.ErrInternal
	}
	return nil
}

func (s *userService) HardDeleteByID(id uint64, currentUserID uint64) error {
	if id != currentUserID {
		return Service.ErrForbidden
	}
	if err := s.userRepo.HardDeleteById(id); err != nil {
		if errors.Is(err, Repository.ErrNotFound) {
			return Service.ErrUserNotFound
		}
		log.Printf("failed to hard delete user %d: %v", id, err)
		return Service.ErrInternal
	}
	return nil
}

func (s *userService) SoftDeleteByID(id uint64, currentUserID uint64) error {
	if id != currentUserID {
		return Service.ErrForbidden
	}
	if err := s.userRepo.SoftDeleteById(id); err != nil {
		if errors.Is(err, Repository.ErrNotFound) {
			return Service.ErrUserNotFound
		}
		log.Printf("failed to soft delete user %d: %v", id, err)
		return Service.ErrInternal
	}
	return nil
}

func (s *userService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := jwt.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", "", Service.ErrInvalidToken
	}
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		if errors.Is(err, Repository.ErrNotFound) {
			return "", "", Service.ErrUserNotFound
		}
		log.Printf("failed to get user for refresh: %v", err)
		return "", "", Service.ErrInternal
	}
	accessToken, newRefreshToken, err := jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		return "", "", Service.ErrInternal
	}
	return accessToken, newRefreshToken, nil
}
