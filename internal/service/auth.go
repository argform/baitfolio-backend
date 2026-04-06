package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/argform/baitfolio-backend/internal/auth"
	"github.com/argform/baitfolio-backend/internal/domain"
	"github.com/argform/baitfolio-backend/internal/repository"
)

type AuthService struct {
	users repository.UserRepository
	jwt *auth.JWTManager
}

type RegisterInput struct {
	Username string
	Email string
	Password string
	FirstName *string
	LastName *string
	About *string
}

type LoginInput struct {
	Email string
	Password string
}

func NewAuthService(users repository.UserRepository, jwt *auth.JWTManager) *AuthService {
	return &AuthService{
		users: users,
		jwt: jwt,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*domain.User, error) {
	if err := auth.ValidatePassword(input.Password); err != nil {
		return nil, fmt.Errorf("validate password: %w", err)
	}

	existingUser, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("check existing user: %w", err)
	}

	if existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	passwordHash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &domain.User{
		Username: input.Username,
		Email: input.Email,
		PasswordHash: passwordHash,
		FirstName: input.FirstName,
		LastName: input.LastName,
		About: input.About,
	}

	createdUser, err := s.users.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (string, error) {
	existingUser, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", fmt.Errorf("password or email incorrect")
		}
		return "", fmt.Errorf("get user by email: %w", err)
	}

	samePassword := auth.ComparePassword(existingUser.PasswordHash, input.Password)
	if !samePassword {
		return "", fmt.Errorf("password or email incorrect")
	}

	token, err := s.jwt.GenerateToken(existingUser.UserID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) GetMe(ctx context.Context, userID uint64) (*domain.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}
