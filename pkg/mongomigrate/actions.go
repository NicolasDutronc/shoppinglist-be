package mongomigrate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

func upAction(c *cli.Context, m *Mongomigrate) error {

	if c.NArg() == 0 {
		return errors.New("Expected an argument, found nothing")
	}

	arg := c.Args().First()
	if arg == "all" {
		return m.MigrateTo(c.Context, -1)
	}

	n, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("Expected a int, found %s : %w", arg, err)
	}

	return m.MigrateTo(c.Context, n)

}

func downAction(c *cli.Context, m *Mongomigrate) error {
	if c.NArg() == 0 {
		return errors.New("Expected an argument, found nothing")
	}

	arg := c.Args().First()
	if arg == "all" {
		return m.RollbackTo(c.Context, -1)
	}

	n, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("Expected a int, found %s : %w", arg, err)
	}

	return m.RollbackTo(c.Context, n)
}

func getVersionAction(c *cli.Context, m *Mongomigrate) error {
	version, err := m.getVersion(c.Context)
	if err != nil {
		return err
	}

	if version.Version == 0 {
		fmt.Printf("database version : %d --> %s\n", version.Version, version.Description)
	} else {
		fmt.Printf("database version : %d --> %s --> %v\n", version.Version, version.Description, version.Timestamp)
	}

	return nil
}

func seedAction(c *cli.Context, m *Mongomigrate) error {
	nameFlag := c.String("name")

	if nameFlag == "" {
		// seed to n
		if c.NArg() == 0 {
			return errors.New("Expected an argument, found nothing")
		}

		arg := c.Args().First()
		if arg == "all" {
			m.SeedTo(c.Context, -1)
			return nil
		}

		n, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("Expected a int, found %s : %w", arg, err)
		}

		return m.SeedTo(c.Context, n)
	}

	// seed by regex
	return m.SeedByRegex(c.Context, nameFlag)
}

func search(c *cli.Context, m *Mongomigrate) error {
	if c.NArg() == 0 {
		return errors.New("Expected an argument, found nothing")
	}

	arg := c.Args().First()

	migrationMatches := 0
	for _, migration := range m.migrations {
		if strings.Contains(migration.Name, arg) {
			fmt.Printf("Migration n°%d : %s\n", migration.ID, migration.Name)
			migrationMatches++
		}
	}

	seederMatches := 0
	for _, seeder := range m.seeders {
		if strings.Contains(seeder.Name, arg) {
			fmt.Printf("Seeder : %s\n", seeder.Name)
			seederMatches++
		}
	}

	fmt.Printf("Matched %d migrations and %d seeders", migrationMatches, seederMatches)

	return nil
}
