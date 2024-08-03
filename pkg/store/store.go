package store

import (
	"fmt"
	"sync"
)

type Key string

type Value []byte

type Storer interface {
	Set(key string, value []byte) error
	Get(key []byte) (Value, error)
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[Key]Value
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		mu:   sync.RWMutex{},
		data: make(map[Key]Value),
	}
}

func (ms *MemoryStore) Set(key string, val []byte) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, ok := ms.data[Key(key)]; ok {
		return fmt.Errorf("key %s already exists", key)
	}
	ms.data[Key(key)] = Value(val)

	return nil
}

func (ms *MemoryStore) Get(key []byte) (Value, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	val, ok := ms.data[Key(key)]
	if !ok {
		return nil, fmt.Errorf("key doesn't exist")
	}
	return val, nil
}
