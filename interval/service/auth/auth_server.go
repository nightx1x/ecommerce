package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	models "github.com/nightx1x/ecommerce/interval/domain"
	repository "github.com/nightx1x/ecommerce/interval/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequset) (*AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	ValidateToken(tokenString string) (*Claims, error)
}
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
}
type LoginRequset struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type AuthResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *models.User `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}
type service struct {
	userRepo      repository.UserRepository
	jwtSecret     string
	tokenDuration time.Duration
}

func NewService(authRepo repository.UserRepository, jwtSecret string) AuthService {
	return &service{
		userRepo:      authRepo,
		jwtSecret:     jwtSecret,
		tokenDuration: 24 * time.Hour,
	}
}

// Login авторизація користувача
func (s *service) Login(ctx context.Context, req LoginRequset) (*AuthResponse, error) {

	if req.Email == "" {
		return nil, ErrEmailRequired
	}
	if req.Password == "" {
		return nil, ErrPasswordRequired
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.generateToken(user, s.tokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate acces token: %w", err)
	}
	refreshToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken оновлення JWT токену
func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {

	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	newAccessToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate acces token: %w", err)
	}
	newRefreshToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

// Register реєстрація користувача
func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {

	if req.Email == "" {
		return nil, ErrEmailRequired
	}
	if req.Password == "" {
		return nil, ErrPasswordRequired
	}
	if len(req.Password) < 8 {
		return nil, ErrWeakPassword
	}

	existUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existUser != nil {
		return nil, ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "customer",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	accessToken, err := s.generateToken(user, s.tokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate acces token: %w", err)
	}
	refreshToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// ValidateToken валідація токену
func (s *service) ValidateToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	return claims, nil
}

// generateToken генерація нового JWT токену
func (s *service) generateToken(user *models.User, duration time.Duration) (string, error) {

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // дата закінчення дії
			IssuedAt:  jwt.NewNumericDate(time.Now()),               // дата створення
			Issuer:    "eccomerce-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
