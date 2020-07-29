package autokey

import "os"

// EnvStorage stores the key in the environment
type EnvStorage struct {
	envVar string
}

// NewEnvStorage inits a new EnvStorage
func NewEnvStorage(envVar string) Storage {
	return &EnvStorage{
		envVar: envVar,
	}
}

// Store set the key in the environment
func (s *EnvStorage) Store(key string) error {
	return os.Setenv(s.envVar, key)
}

// IsStored checks if the environment variable is set
func (s *EnvStorage) IsStored() (bool, error) {
	key := os.Getenv(s.envVar)
	if len(key) == 0 {
		return false, nil
	}

	return true, nil
}
