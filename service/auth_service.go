package service

import (
	"be-porto-v3/models"
	"be-porto-v3/repository"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService menangani logika bisnis untuk autentikasi
type AuthService struct {
	authRepo *repository.AuthRepository
}

// NewAuthService membuat instance baru AuthService
func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{authRepo}
}

// RegisterUser menangani proses registrasi user
func (s *AuthService) RegisterUser(username, email, password string) (*models.User, error) {
	// Cek apakah email sudah terdaftar
	existingUser, _ := s.authRepo.GetUserByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Buat user baru
	user := &models.User{
		Username: username,
		Email:    email,
		Password: password, // Password akan di-hash di repository
	}

	// Simpan ke database
	err := s.authRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) LoginUser(email, password string) (string, error) {
	user, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Membuat token JWT dengan jwt v5
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
