package service

import (
	"context"
	"net/http"
	"peer-link-server/internal/model"
	"peer-link-server/internal/repository"
	apperr "peer-link-server/pkg/errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetByID(ctx context.Context, id uint) (*model.User, error)
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
	Create(ctx context.Context, req *CreateUserRequest) (*model.User, error)
	Update(ctx context.Context, id uint, req *UpdateUserRequest) (*model.User, error)
	Delete(ctx context.Context, id uint) error
}

type CreateUserRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	Name string `json:"name" binding:"omitempty,min=2,max=100"`
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{repo: repo, logger: logger}
}

func (s *userService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.List(ctx, (page-1)*pageSize, pageSize)
}

func (s *userService) Create(ctx context.Context, req *CreateUserRequest) (*model.User, error) {
	// 邮箱唯一性检查
	if _, err := s.repo.FindByEmail(ctx, req.Email); err == nil {
		return nil, apperr.New(http.StatusConflict, 40900, "email already registered")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperr.Wrap(http.StatusInternalServerError, 50001, "hash password failed", err)
	}
	user := &model.User{Name: req.Name, Email: req.Email, Password: string(hashed)}
	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Error("create user", zap.Error(err))
		return nil, apperr.ErrInternal
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, id uint, req *UpdateUserRequest) (*model.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Error("update user", zap.Uint("id", id), zap.Error(err))
		return nil, apperr.ErrInternal
	}
	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	if _, err := s.repo.FindByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
