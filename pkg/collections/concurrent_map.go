package collections

import (
	"errors"
	"sync"
)

var (
	ErrEntityAlreadyExists = errors.New("entity already exists")
)

type ConcurrentMap[K comparable, V any] interface {
	Get(key K) (V, bool)
	Remove(key K)
	Exists(key K) bool
	SetIfNotExists(key K, val V) error
}

type concurrentMap[K comparable, V any] struct {
	mu sync.Mutex
	m  map[K]V
}

func NewConcurrentMap[K comparable, V any](size int) ConcurrentMap[K, V] {
	return &concurrentMap[K, V]{
		m: make(map[K]V, size),
	}
}

func (c *concurrentMap[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.m[key]
	return v, ok
}

func (c *concurrentMap[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, key)
}

func (c *concurrentMap[K, V]) Exists(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.m[key]
	return ok
}

func (c *concurrentMap[K, V]) SetIfNotExists(key K, val V) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.m[key]
	if ok {
		return ErrEntityAlreadyExists
	}

	c.m[key] = val
	return nil
}
