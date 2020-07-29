package mongomigrate

import (
	"context"
	"sort"
	"time"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// versionRecord is what is stored in the database to keep track of the version
type versionRecord struct {
	Version     uint64    `bson:"version"`
	Description string    `bson:"description"`
	Timestamp   time.Time `bson:"ts"`
}

// Mongomigrate is a struct holding the migrations
type Mongomigrate struct {
	db             *mongo.Database
	collectionName string
	migrations     []*Migration
	seeders        []*Seeder
}

// New is a constructor for Mongomigrate
func New(db *mongo.Database, collectionName string, migrations []*Migration, seeders []*Seeder) *Mongomigrate {
	return &Mongomigrate{
		db:             db,
		collectionName: collectionName,
		migrations:     migrations,
		seeders:        seeders,
	}
}

// sortMigrations is a helper for sorting the migrations based on their IDs
func (m *Mongomigrate) sortMigrations() {
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].ID < m.migrations[j].ID
	})
}

func (m *Mongomigrate) getMigrationCollection() *mongo.Collection {
	return m.db.Collection(m.collectionName)
}

// getVersion returns the last version in the migrations collection
func (m *Mongomigrate) getVersion(ctx context.Context) (*versionRecord, error) {
	var rec versionRecord
	if err := m.getMigrationCollection().FindOne(ctx, bson.M{}, &options.FindOneOptions{
		Sort: bson.D{
			{
				Key:   "version",
				Value: -1,
			},
		},
	}).Decode(&rec); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &versionRecord{
				Version:     0,
				Description: "Database is fresh",
			}, nil
		}
		return nil, err
	}

	return &rec, nil
}

// SetVersion inserts a new version record in the migrations collection
// If the given version is 0, nothing happens
func (m *Mongomigrate) SetVersion(ctx context.Context, migration *Migration) error {
	if m.IsApplied(ctx, migration) {
		// rollback
		if _, err := m.getMigrationCollection().DeleteOne(ctx, bson.M{
			"version": migration.ID,
		}); err != nil {
			return err
		}

		return nil
	}

	// migrate
	_, err := m.getMigrationCollection().InsertOne(ctx, &versionRecord{
		Version:     migration.ID,
		Description: migration.Name,
		Timestamp:   time.Now().UTC(),
	})

	return err
}

// IsApplied returns true if the given migration has already been applied
func (m *Mongomigrate) IsApplied(ctx context.Context, migration *Migration) bool {
	if m.getMigrationCollection().FindOne(ctx, bson.M{
		"version": migration.ID,
	}).Err() != nil {
		return false
	}

	return true
}
