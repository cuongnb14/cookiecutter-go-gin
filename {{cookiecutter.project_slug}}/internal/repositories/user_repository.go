package repositories

import (
	"{{ cookiecutter.project_slug }}/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	IRepository[models.User, uuid.UUID]
}

type userRepository struct {
	repository[models.User, uuid.UUID]
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		repository[models.User, uuid.UUID]{
			db: db,
		},
	}
}
