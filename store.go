package main

import (
	"fmt"
	"sync"
)

type Key string

type Value []byte

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

func (ms *MemoryStore) Set(key string, val string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.data[Key(key)] = Value(val)

	return nil
}

func (ms *MemoryStore) Get(key string) (Value, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	val, ok := ms.data[Key(key)]
	if !ok {
		return nil, fmt.Errorf("key doesn't exist")
	}

	return val, nil
}
