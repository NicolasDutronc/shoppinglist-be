package mongomigrate

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// SeedFunc represents a operation that seeds the database
type SeedFunc func(ctx context.Context, db *mongo.Database) error

// Seeder holds a seed function and a name
type Seeder struct {
	Name string
	Seed SeedFunc
}

// SeedTo runs the first n seeders
func (m *Mongomigrate) SeedTo(ctx context.Context, n int) error {
	if n <= 0 || n > len(m.seeders) {
		n = len(m.seeders)
	}

	for i := 0; i != len(m.seeders); i++ {
		if err := m.seeders[i].Seed(ctx, m.db); err != nil {
			return fmt.Errorf("Error seeding %d/%d with name %s : %w", i, len(m.seeders), m.seeders[i].Name, err)
		}
	}

	return nil
}

// SeedByRegex runs the seeders whose names contains the given regex
func (m *Mongomigrate) SeedByRegex(ctx context.Context, regex string) error {
	matches := 0
	for i, seeder := range m.seeders {
		if !strings.Contains(seeder.Name, regex) {
			continue
		}

		if err := seeder.Seed(ctx, m.db); err != nil {
			return fmt.Errorf("Error seeding %d/%d with name %s : %w", i, len(m.seeders), seeder.Name, err)
		}

		fmt.Printf("Seeding %d/%d : %s", i, len(m.seeders), seeder.Name)

		matches++
	}

	if matches == 0 {
		return fmt.Errorf("No match found with %s", regex)
	}

	return nil
}
