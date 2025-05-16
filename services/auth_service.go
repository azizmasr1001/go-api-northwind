package services

import (
	"context"
	"fmt"
	"time"

	"github.com/aronipurwanto/go-api-northwind/internal/redis"
	"github.com/aronipurwanto/go-api-northwind/models"
	"github.com/aronipurwanto/go-api-northwind/repositories"
	"github.com/golang-jwt/jwt/v5"
	redislib "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, identifier, password string) (*models.User, error)
	Register(ctx context.Context, input *models.RegisterRequest) (*models.User, error)
	GenerateAccessToken(user *models.User) (string, error)
	GenerateRefreshToken(user *models.User) (string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	Logout(ctx context.Context, token string, exp time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	GetSecret() string

	SendResetOTP(ctx context.Context, email string) error
	VerifyResetOTP(ctx context.Context, email, otp string) (bool, error)
	ResetPassword(ctx context.Context, email, newPassword string) error
}

type AuthServiceImpl struct {
	repo   repositories.AuthRepository
	secret string
	redis  *redislib.Client
}

func NewAuthService(repo repositories.AuthRepository, secret string, redis *redislib.Client) AuthService {
	return &AuthServiceImpl{repo: repo, secret: secret, redis: redis}
}

func (s *AuthServiceImpl) Login(ctx context.Context, identifier, password string) (*models.User, error) {
	user, err := s.repo.FindByUsername(identifier)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return user, nil
}

func (s *AuthServiceImpl) Register(ctx context.Context, input *models.RegisterRequest) (*models.User, error) {
	if exists, _ := s.repo.CheckUsernameExists(input.Username); exists {
		return nil, fmt.Errorf("username already taken")
	}
	if exists, _ := s.repo.CheckEmailExists(input.Email); exists {
		return nil, fmt.Errorf("email already taken")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashed),
		IsActive:     true,
		Role:         "user",
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthServiceImpl) GenerateAccessToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.UserID,
		"firstName": user.Username,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *AuthServiceImpl) GenerateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.UserID,
		"firstName": user.Username,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid refresh token")
	}

	isBlacklisted, err := s.IsTokenBlacklisted(ctx, refreshToken)
	if err != nil {
		return "", err
	}
	if isBlacklisted {
		return "", fmt.Errorf("refresh token revoked")
	}

	user := &models.User{
		UserID:   int(claims["user_id"].(float64)),
		Username: claims["firstName"].(string),
	}

	return s.GenerateAccessToken(user)
}

func (s *AuthServiceImpl) Logout(ctx context.Context, token string, exp time.Duration) error {
	return redis.SetBlacklistToken(ctx, s.redis, token, exp)
}

func (s *AuthServiceImpl) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	return redis.IsTokenBlacklisted(ctx, s.redis, token)
}

func (s *AuthServiceImpl) GetSecret() string {
	return s.secret
}

// Reset Password OTP

func (s *AuthServiceImpl) SendResetOTP(ctx context.Context, email string) error {
	_, err := s.repo.FindByUsername(email)
	if err != nil {
		return fmt.Errorf("email not found")
	}
	otp := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	key := fmt.Sprintf("otp:%s", email)
	err = s.redis.Set(ctx, key, otp, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	fmt.Printf("[DEBUG] Send OTP to %s: %s\n", email, otp)
	return nil
}

func (s *AuthServiceImpl) VerifyResetOTP(ctx context.Context, email, otp string) (bool, error) {
	key := fmt.Sprintf("otp:%s", email)
	storedOTP, err := s.redis.Get(ctx, key).Result()
	if err == redislib.Nil {
		return false, fmt.Errorf("OTP expired or not found")
	} else if err != nil {
		return false, err
	}
	return storedOTP == otp, nil
}

func (s *AuthServiceImpl) ResetPassword(ctx context.Context, email, newPassword string) error {
	user, err := s.repo.FindByUsername(email)
	if err != nil {
		return fmt.Errorf("user not found")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}
	user.PasswordHash = string(hashed)
	return s.repo.UpdatePassword(user)
}
