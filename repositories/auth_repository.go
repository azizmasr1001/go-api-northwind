package repositories

import (
	"github.com/aronipurwanto/go-api-northwind/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	CheckUsernameExists(username string) (bool, error)
	CheckEmailExists(email string) (bool, error)
	UpdatePassword(user *models.User) error
}

type authRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepositoryImpl{db: db}
}

func (r *authRepositoryImpl) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.
		Where("Username = ? OR Email = ?", username, username).
		First(&user).Error
	return &user, err
}

func (r *authRepositoryImpl) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("Email = ?", email).First(&user).Error
	return &user, err
}

func (r *authRepositoryImpl) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *authRepositoryImpl) CheckUsernameExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("Username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *authRepositoryImpl) CheckEmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("Email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *authRepositoryImpl) UpdatePassword(user *models.User) error {
	return r.db.Model(user).Update("PasswordHash", user.PasswordHash).Error
}
