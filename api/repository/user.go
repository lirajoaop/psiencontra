package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/schemas"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(u *schemas.User) error {
	return r.DB.Create(u).Error
}

func (r *UserRepo) FindByID(id uuid.UUID) (*schemas.User, error) {
	var u schemas.User
	err := r.DB.First(&u, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByEmail(email string) (*schemas.User, error) {
	var u schemas.User
	err := r.DB.First(&u, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) FindByGoogleID(googleID string) (*schemas.User, error) {
	var u schemas.User
	err := r.DB.First(&u, "google_id = ?", googleID).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) Update(u *schemas.User) error {
	return r.DB.Save(u).Error
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
