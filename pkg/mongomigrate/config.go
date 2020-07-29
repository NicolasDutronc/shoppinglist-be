package mongomigrate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// ErrKeyNotFound is an error indicating that a configuration field has not been found in the environment
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
	return fmt.Sprintf("%v was not found neither in the configuration nor the environment", e.Key)
}

// Config holds the configuration for mongomigrate
type Config struct {
	Database struct {
		Hostname string `mapstructure:"hostname"`
		Port     string `mapstructure:"port"`
		DB       string `mapstructure:"db"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"database"`
	Migrations struct {
		DB         string `mapstructure:"db"`
		Collection string `mapstructure:"collection"`
	} `mapstructure:"migrations"`
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

	v = v.Sub("dbctl")

	// fmt.Println("all settings :", v.AllSettings())
	// fmt.Println("all keys :", v.AllKeys())
	// fmt.Println("hostname :", v.GetString("database.hostname"))

	// Check for missing parameters
	t := reflect.TypeOf(Config{})
	for i := 0; i != t.NumField(); i++ {
		ti := t.Field(i).Type
		for j := 0; j != ti.NumField(); j++ {
			key := fmt.Sprintf("%v.%v", t.Field(i).Tag.Get("mapstructure"), ti.Field(j).Tag.Get("mapstructure"))
			// fmt.Printf("%v : %t\n", key, v.IsSet(key))
			if !v.IsSet(key) {
				return nil, NewErrKeyNotFound(key)
			}
		}
	}

	var c Config
	if err := v.UnmarshalExact(&c); err != nil {
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
		c.Database.DB,
	)
}
