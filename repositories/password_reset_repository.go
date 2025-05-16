package repositories

import (
	"github.com/azizmasr1001/go-api-northwind/models"
	"gorm.io/gorm"
)

type PasswordResetRepository interface {
	Save(reset *models.PasswordReset) error
	FindByToken(token string) (*models.PasswordReset, error)
	MarkUsed(id int) error
}

type passwordResetRepo struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepo{db}
}

func (r *passwordResetRepo) Save(reset *models.PasswordReset) error {
	return r.db.Create(reset).Error
}

func (r *passwordResetRepo) FindByToken(token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.Where("token = ? AND used = false AND expires_at > GETDATE()", token).First(&reset).Error
	return &reset, err
}

func (r *passwordResetRepo) MarkUsed(id int) error {
	return r.db.Model(&models.PasswordReset{}).Where("id = ?", id).Update("used", true).Error
}
