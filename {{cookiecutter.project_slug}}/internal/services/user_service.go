package services

import (
	"{{ cookiecutter.project_slug }}/configs"
	"{{ cookiecutter.project_slug }}/internal/models"
	"{{ cookiecutter.project_slug }}/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var logger = configs.GetLogger()

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

func (s *UserService) Login(token string, agencyId uuid.UUID) (*models.User, string, error) {
	return nil, nil
}
