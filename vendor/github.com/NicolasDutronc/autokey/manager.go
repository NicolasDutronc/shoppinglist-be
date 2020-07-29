package autokey

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// Manager is responsible for generating and rotating the application key that is used to sign JWTs
// The key is a base 64 encoded random string
type Manager struct {
	storage      Storage
	rotationTime time.Duration
	timer        *time.Timer
	stopChan     chan struct{}
	keySize      int
}

// NewManager is a constructor for Manager
func NewManager(storage Storage, keySize int, rotationTime time.Duration) *Manager {
	return &Manager{
		storage:      storage,
		rotationTime: rotationTime,
		timer:        nil,
		stopChan:     nil,
		keySize:      keySize,
	}
}

// generate returns a base 64 encoded random string
func (m *Manager) generate() (string, error) {
	key := make([]byte, m.keySize)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// rotate generates and stores a new key
func (m *Manager) rotate() error {
	// generate a new key
	newKey, err := m.generate()
	if err != nil {
		return err
	}

	// store the key
	if err := m.storage.Store(newKey); err != nil {
		return err
	}

	return nil
}

// Start starts the rotation loop
func (m *Manager) Start(interrupt chan struct{}) error {
	defer func() {
		interrupt <- struct{}{}
	}()
	// rotate a first time if needed
	stored, err := m.storage.IsStored()
	if err != nil {
		return err
	}
	if !stored {
		if err := m.rotate(); err != nil {
			return err
		}
	}

	// init stop channel
	m.stopChan = make(chan struct{})

	// init timer
	m.timer = time.NewTimer(m.rotationTime)

	for {
		select {
		case <-m.timer.C:
			if err := m.rotate(); err != nil {
				return err
			}
			m.timer.Reset(m.rotationTime)
		case <-m.stopChan:
			return nil
		}
	}
}

// Stop stops the rotation loop
func (m *Manager) Stop() {
	m.timer.Stop()
	close(m.stopChan)
}
