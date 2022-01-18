package user

import (
	"context"
	"time"

	"github.com/NicolasDutronc/autokey"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// ServiceImpl is the concrete implementation of the user service interface
type ServiceImpl struct {
	repo        Repository
	keyConsumer autokey.Consumer
}

// NewService inits a new user service
func NewService(repo Repository, keyConsumer autokey.Consumer) Service {
	return &ServiceImpl{
		repo:        repo,
		keyConsumer: keyConsumer,
	}
}

// FindByID directly calls the repository
func (s *ServiceImpl) FindByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.FindByID(ctx, userID)
}

// FindByName directly calls the repository
func (s *ServiceImpl) FindByName(ctx context.Context, userName string) (*User, error) {
	return s.repo.FindByName(ctx, userName)
}

// Store hashes the password and then calls the repository
func (s *ServiceImpl) Store(ctx context.Context, name string, password string, permissions ...*Permission) (*User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.repo.Store(ctx, name, string(encryptedPassword), permissions...)
}

// UpdateName directly calls the repository
func (s *ServiceImpl) UpdateName(ctx context.Context, userID string, newName string) (int64, error) {
	return s.repo.UpdateName(ctx, userID, newName)
}

// UpdatePassword compares the given password with the one stored in the database. If this is a match, it hashes the new password before calling the repository
func (s *ServiceImpl) UpdatePassword(ctx context.Context, userID string, currentPassword string, newPassword string) (int64, error) {
	user, err := s.FindByID(ctx, userID)
	if err != nil {
		return -1, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return -1, err
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return -1, err
	}

	return s.repo.UpdatePassword(ctx, userID, string(encryptedPassword))
}

// Delete directly calls the repository
func (s *ServiceImpl) Delete(ctx context.Context, userID string) (int64, error) {
	return s.repo.Delete(ctx, userID)
}

// AddPermissions adds permissions to the user
func (s *ServiceImpl) AddPermissions(ctx context.Context, userID string, permissions ...*Permission) (int64, error) {
	return s.repo.AddPermissions(ctx, userID, permissions...)
}

// RemovePermissions removes permissions from the user
func (s *ServiceImpl) RemovePermissions(ctx context.Context, userID string, permissions ...*Permission) (int64, error) {
	return s.repo.RemovePermissions(ctx, userID, permissions...)
}

// Login takes in a user name and a password and returns the corresponding user along with a token.
// An error is returned instead if it cannot retrieve the user based on the given user name or if the given password does not match or if there is any issue retrieving the key and signing the token
func (s *ServiceImpl) Login(ctx context.Context, userName string, password string) (*User, string, error) {
	user, err := s.FindByName(ctx, userName)
	if err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", err
	}

	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		NotBefore: time.Now().Unix(),
		Subject:   user.ID.Hex(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	key, err := s.keyConsumer.Get()
	if err != nil {
		return nil, "", err
	}

	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, "", err
	}

	return user, ss, nil
}

// Authenticate retrieves a user from a token after validating it
func (s *ServiceImpl) Authenticate(ctx context.Context, token string) (*User, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, jwt.ErrSignatureInvalid
		}

		key, err := s.keyConsumer.Get()
		if err != nil {
			return nil, err
		}

		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, jwt.NewValidationError("Claims are not standard claims", jwt.ValidationErrorClaimsInvalid)
	}
	if err := claims.Valid(); err != nil {
		return nil, err
	}

	return s.FindByID(ctx, claims.Subject)
}
