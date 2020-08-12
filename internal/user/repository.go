package user

import "context"

// FinderByID is a single method interface for finding a user from the database based on its id
type FinderByID interface {
	FindByID(ctx context.Context, userID string) (*User, error)
}

// FinderByName is a single method interface for finding a user from the database based on its name
type FinderByName interface {
	FindByName(ctx context.Context, userName string) (*User, error)
}

// Storer is a single method interface for storing a user in the database
type Storer interface {
	Store(ctx context.Context, name string, password string, permissions ...*Permission) (*User, error)
}

// NameUpdater is a single method interface for updating the name field of the user given by its id
type NameUpdater interface {
	UpdateName(ctx context.Context, userID string, newName string) (int64, error)
}

// PasswordUpdater is a single method interface for updating the password field of a user
type PasswordUpdater interface {
	UpdatePassword(ctx context.Context, userID string, currentPassword string, newPassword string) (int64, error)
}

// Deleter is a single method interface for deleting a user from the database
type Deleter interface {
	Delete(ctx context.Context, userID string) (int64, error)
}

// Repository defines all possible actions on users database
type Repository interface {
	FinderByID
	FinderByName
	Storer
	NameUpdater
	PasswordUpdater
	Deleter
}
