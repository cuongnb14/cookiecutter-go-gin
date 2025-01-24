package services

import (
	"{{ cookiecutter.project_slug }}/internal/core/models"
	"{{ cookiecutter.project_slug }}/internal/core/repositories"

	"github.com/google/uuid"
)

type UserService struct {
	userRepository repositories.IUserRepository
}

func NewUserService(userRepo repositories.IUserRepository) *UserService {
	return &UserService{
		userRepo,
	}
}

func (s *UserService) Login(token string, agencyId uuid.UUID) (*models.User, error) {
	return nil, nil
}
