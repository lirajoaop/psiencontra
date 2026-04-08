package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joaop/psiencontra/api/repository"
	"github.com/joaop/psiencontra/api/schemas"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenLifetime = 7 * 24 * time.Hour
	bcryptCost    = 12
)

var (
	ErrEmailAlreadyUsed   = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordRequired   = errors.New("password required")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrEmailRequired      = errors.New("email required")
)

// UserRepository is the subset of repository.UserRepo that AuthService needs.
// Defined as an interface so tests can substitute an in-memory fake without
// touching a real database.
type UserRepository interface {
	Create(*schemas.User) error
	FindByID(uuid.UUID) (*schemas.User, error)
	FindByEmail(string) (*schemas.User, error)
	FindByGoogleID(string) (*schemas.User, error)
	Update(*schemas.User) error
}

type AuthService struct {
	userRepo  UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

type Claims struct {
	UserID uuid.UUID `json:"uid"`
	jwt.RegisteredClaims
}

// Register creates a new email/password user.
func (s *AuthService) Register(email, password, name string) (*schemas.User, string, error) {
	email = normalizeEmail(email)
	if email == "" {
		return nil, "", ErrEmailRequired
	}
	if len(password) < 8 {
		return nil, "", ErrPasswordTooShort
	}

	if existing, err := s.userRepo.FindByEmail(email); err == nil && existing != nil {
		return nil, "", ErrEmailAlreadyUsed
	} else if err != nil && !repository.IsNotFound(err) {
		return nil, "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	user := &schemas.User{
		Email:        email,
		Name:         strings.TrimSpace(name),
		PasswordHash: string(hash),
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	token, err := s.IssueToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// Login validates email/password and returns a JWT.
func (s *AuthService) Login(email, password string) (*schemas.User, string, error) {
	email = normalizeEmail(email)
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if repository.IsNotFound(err) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}
	if user.PasswordHash == "" {
		// User registered via Google — no password set.
		return nil, "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}
	token, err := s.IssueToken(user.ID)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil
}

// UpsertGoogleUser links a Google profile to an existing user (by email or
// google_id) or creates a new one.
func (s *AuthService) UpsertGoogleUser(googleID, email, name, avatarURL string) (*schemas.User, string, error) {
	email = normalizeEmail(email)

	if existing, err := s.userRepo.FindByGoogleID(googleID); err == nil && existing != nil {
		token, err := s.IssueToken(existing.ID)
		return existing, token, err
	} else if err != nil && !repository.IsNotFound(err) {
		return nil, "", err
	}

	if existing, err := s.userRepo.FindByEmail(email); err == nil && existing != nil {
		existing.GoogleID = googleID
		if existing.Name == "" {
			existing.Name = name
		}
		if existing.AvatarURL == "" {
			existing.AvatarURL = avatarURL
		}
		if err := s.userRepo.Update(existing); err != nil {
			return nil, "", err
		}
		token, err := s.IssueToken(existing.ID)
		return existing, token, err
	} else if err != nil && !repository.IsNotFound(err) {
		return nil, "", err
	}

	user := &schemas.User{
		Email:     email,
		Name:      name,
		GoogleID:  googleID,
		AvatarURL: avatarURL,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}
	token, err := s.IssueToken(user.ID)
	return user, token, err
}

func (s *AuthService) IssueToken(userID uuid.UUID) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "psiencontra",
			Subject:   userID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *AuthService) GetUser(id uuid.UUID) (*schemas.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *AuthService) TokenLifetime() time.Duration {
	return tokenLifetime
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
