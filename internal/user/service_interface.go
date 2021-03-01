package user

import context "context"

// LoginService defines the login interface
type LoginService interface {
	Login(ctx context.Context, userName string, password string) (*User, string, error)
}

// AuthenticateService defines the authenticate interface
type AuthenticateService interface {
	Authenticate(ctx context.Context, token string) (*User, error)
}

// Service defines the user service
type Service interface {
	LoginService
	AuthenticateService

	FindByID(ctx context.Context, userID string) (*User, error)

	FindByName(ctx context.Context, userName string) (*User, error)

	Store(ctx context.Context, name string, password string, permissions ...*Permission) (*User, error)

	UpdateName(ctx context.Context, userID string, newName string) (int64, error)

	UpdatePassword(ctx context.Context, userID string, currentPassword string, newPassword string) (int64, error)

	Delete(ctx context.Context, userID string) (int64, error)

	AddPermissions(ctx context.Context, userID string, permissions ...*Permission) (int64, error)

	RemovePermissions(ctx context.Context, userID string, permissions ...*Permission) (int64, error)
}
