package main

import (
	"fmt"
	"log"
	"sync"
)

type Storer[K comparable, V any] interface {
	Put(K, V) error
	Get(K) (V, error)
	Update(K, V) error
	Delete(K) (V, error)
}

type KVStore[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	return &KVStore[K, V]{
		data: make(map[K]V),
	}
}

func (s *KVStore[K, V]) Put(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	return nil
}

func (s *KVStore[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return value, fmt.Errorf("the key (%v) does not exist", key)
	}

	return value, nil
}

// func StoreThings(s Storer[string, int]) error {
// 	return s.Put("foo", 1)
// }

func main() {
	store := NewKVStore[string, string]()

	if err := store.Put("King", "Aryan"); err != nil {
		log.Fatal(err)
	}

	value, err := store.Get("King")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(value)

	// StoreThings(kv)
}
