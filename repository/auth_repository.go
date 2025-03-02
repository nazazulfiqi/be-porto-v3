package repository

import (
	"be-porto-v3/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthRepository menangani operasi database untuk user
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository membuat instance baru AuthRepository
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser menambahkan user baru ke database
func (r *AuthRepository) CreateUser(user *models.User) error {
	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Simpan user ke database
	return r.db.Create(user).Error
}

// GetUserByEmail mencari user berdasarkan email
func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
