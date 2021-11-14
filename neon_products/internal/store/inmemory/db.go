package inmemory

import (
	"context"
	"fmt"
	"lectures-6/internal/models"
	"lectures-6/internal/store"
	"sync"
)

type DB struct {
	data map[int]*models.Neon

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		data: make(map[int]*models.Neon),
		mu:   new(sync.RWMutex),
	}
}

func (db *DB) Create(ctx context.Context, neon *models.Neon) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[neon.ID] = neon
	return nil
}

func (db *DB) All(ctx context.Context) ([]*models.Neon, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	neons := make([]*models.Neon, 0, len(db.data))
	for _, neon := range db.data {
		neons = append(neons, neon)
	}

	return neons, nil
}

func (db *DB) ByID(ctx context.Context, id int) (*models.Neon, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	neon, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("No neon with id %d", id)
	}

	return neon, nil
}

func (db *DB) Update(ctx context.Context, neon *models.Neon) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[neon.ID] = neon
	return nil
}

func (db *DB) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
