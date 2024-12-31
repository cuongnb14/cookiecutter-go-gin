package services

import (
	"{{ cookiecutter.project_slug }}/internal/core/models"
	"{{ cookiecutter.project_slug }}/internal/core/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	db             *gorm.DB
	userRepository repositories.IUserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db,
		repositories.NewUserRepository(db),
	}
}

func (s *UserService) Login(token string, agencyId uuid.UUID) (*models.User, error) {
	return nil, nil
}
