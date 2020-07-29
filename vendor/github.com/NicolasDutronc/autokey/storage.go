package autokey

// Storage defines how a key is stored
type Storage interface {
	Store(key string) error
	IsStored() (bool, error)
}
