// internal/repository/store_repository.go
package repository

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"kirana-club-assignment/backend/internal/models"
)

var (
	ErrStoreNotFound = errors.New("store not found")
)

type StoreRepository interface {
	GetByID(id string) (models.Store, error)
	LoadStoresFromFile(filePath string) error
}

type InMemoryStoreRepository struct {
	mu     sync.RWMutex
	stores map[string]models.Store
}

func NewInMemoryStoreRepository() *InMemoryStoreRepository {
	return &InMemoryStoreRepository{
		stores: make(map[string]models.Store),
	}
}

func (r *InMemoryStoreRepository) GetByID(id string) (models.Store, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	store, exists := r.stores[id]
	if !exists {
		return models.Store{}, ErrStoreNotFound
	}
	return store, nil
}

func (r *InMemoryStoreRepository) LoadStoresFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var stores []models.Store
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&stores); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	for _, store := range stores {
		r.stores[store.StoreID] = store
	}
	return nil
}
