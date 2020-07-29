package mongomigrate

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

// MigrateFunc is a migrate function
type MigrateFunc func(ctx context.Context, db *mongo.Database) error

// RollbackFunc is a rollback function
type RollbackFunc func(ctx context.Context, db *mongo.Database) error

// Migration represents a migration
type Migration struct {
	ID       uint64
	Name     string
	Migrate  MigrateFunc
	Rollback RollbackFunc
}

// MigrateTo runs the first migrationIndex migrations
func (m *Mongomigrate) MigrateTo(ctx context.Context, migrationIndex int) error {
	if migrationIndex <= 0 || migrationIndex > len(m.migrations) {
		migrationIndex = len(m.migrations)
	}

	m.sortMigrations()

	for i := 0; i < migrationIndex; i++ {
		if m.IsApplied(ctx, m.migrations[i]) || m.migrations[i].Migrate == nil {
			continue
		}
		if err := m.migrations[i].Migrate(ctx, m.db); err != nil {
			return fmt.Errorf("Error during migration n°%d with ID %d : %w", i, m.migrations[i].ID, err)
		}

		log.Printf("Migrating %d/%d : %d --> %s", i+1, migrationIndex, m.migrations[i].ID, m.migrations[i].Name)

		if err := m.SetVersion(ctx, m.migrations[i]); err != nil {
			return fmt.Errorf("Error setting version to %d after migrating migration n°%d with ID %d : %w", m.migrations[i].ID, i, m.migrations[i].ID, err)
		}
	}

	return nil
}

// RollbackTo rolls back the last migrationIndex migrations
func (m *Mongomigrate) RollbackTo(ctx context.Context, migrationIndex int) error {
	if migrationIndex <= 0 || migrationIndex > len(m.migrations) {
		migrationIndex = 0
	}

	m.sortMigrations()

	for i := len(m.migrations) - 1; i >= migrationIndex; i-- {
		if !m.IsApplied(ctx, m.migrations[i]) || m.migrations[i].Rollback == nil {
			continue
		}
		if err := m.migrations[i].Rollback(ctx, m.db); err != nil {
			return fmt.Errorf("Error during rollback of migration n°%d with ID %d : %w", i, m.migrations[i].ID, err)
		}

		log.Printf("Rollbacking %d/%d : %d --> %s", len(m.migrations)-migrationIndex-i, len(m.migrations)-migrationIndex, m.migrations[i].ID, m.migrations[i].Name)

		if err := m.SetVersion(ctx, m.migrations[i]); err != nil {
			return fmt.Errorf("Error setting version to %d after rollbacking migration n°%d with ID %d : %w", m.migrations[i].ID, i, m.migrations[i].ID, err)
		}
		i++
	}

	return nil
}
