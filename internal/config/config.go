package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// ErrKeyNotFound is an error indicating that a configuration field has not been configured
type ErrKeyNotFound struct {
	Key string
}

// NewErrKeyNotFound returns a new ErrNewErrKeyNotFound with the given key
func NewErrKeyNotFound(key string) *ErrKeyNotFound {
	return &ErrKeyNotFound{
		Key: key,
	}
}

// Error returns the error message
// It makes ErrKeyNotFound implements the Error interface
func (e *ErrKeyNotFound) Error() string {
	return fmt.Sprintf("%v was not configured", e.Key)
}

// Config contains server configuration
type Config struct {
	key       string
	KeyConfig struct {
		Size          int           `mapstructure:"size"`
		ValidDuration time.Duration `mapstructure:"validation_duration"`
	} `mapstructure:"key"`
	Server struct {
		Hostname  string `mapstructure:"hostname"`
		Port      string `mapstructure:"port"`
		ServerCRT string `mapstructure:"crt"`
		ServerKey string `mapstructure:"key"`
	} `mapstructure:"server"`
	Database struct {
		Username        string `mapstructure:"username"`
		Password        string `mapstructure:"password"`
		Hostname        string `mapstructure:"hostname"`
		Port            string `mapstructure:"port"`
		Name            string `mapstructure:"db"`
		ListsCollection string `mapstructure:"lists_collection"`
		UsersCollection string `mapstructure:"users_collection"`
	} `mapstructure:"database"`
}

// NewConfig reads the given configuration file and the environment and returns a newly created Config
// An error is returned if a field was not found in the configuration (file or env)
func NewConfig(filepath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(filepath)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v = v.Sub("app")

	t := reflect.TypeOf(Config{})
	for i := 0; i != t.NumField(); i++ {
		key := t.Field(i).Tag.Get("mapstructure")
		if len(key) > 0 && !v.IsSet(key) {
			return nil, NewErrKeyNotFound(key)
		}
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

// BuildMongoDBConnexionString returns the mongoDB connexion url
func (c *Config) BuildMongoDBConnexionString() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s",
		c.Database.Username,
		c.Database.Password,
		c.Database.Hostname,
		c.Database.Port,
		c.Database.Name,
	)
}

// BuildServerAdress returns the adress that is served by the server
func (c *Config) BuildServerAdress() string {
	return fmt.Sprintf("%s:%s", c.Server.Hostname, c.Server.Port)
}

// Store sets the app key. It satisfies the key storage interface
func (c *Config) Store(key string) error {
	c.key = key
	return nil
}

// IsStored checks if the key is not empty. It satisfies the key storage interface
func (c *Config) IsStored() (bool, error) {
	return c.key == "", nil
}

// Get returns the app key. It satisfies the key consumer interface
func (c *Config) Get() (string, error) {
	stored, err := c.IsStored()
	if err != nil {
		return "", err
	}

	if !stored {
		return "", errors.New("The key was not found")
	}

	return c.key, nil
}
