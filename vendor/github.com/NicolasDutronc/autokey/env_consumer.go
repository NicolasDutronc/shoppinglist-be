package autokey

import (
	"fmt"
	"os"
)

// EnvConsumer is key consumer that is able to retrieve the key from the environment
type EnvConsumer struct {
	envVar string
}

// NewEnvConsumer inits a new EnvConsumer
func NewEnvConsumer(envVar string) Consumer {
	return &EnvConsumer{
		envVar: envVar,
	}
}

// Get retrieves the key from the environment
// An error is returned instead is the environment key was not set
func (c *EnvConsumer) Get() (string, error) {
	key := os.Getenv(c.envVar)
	if len(key) == 0 {
		return "", fmt.Errorf("%s was not set in the environment", c.envVar)
	}

	return key, nil
}
