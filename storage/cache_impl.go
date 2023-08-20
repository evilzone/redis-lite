package storage

import (
	"errors"
	"sync"
	"time"
)

const (
	DEFAULT = -1
)

type Item struct {
	Value    []byte
	ExpireAt int64
}

type InMemoryStorage struct {
	mu   sync.Mutex
	Data map[string]Item
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		Data: make(map[string]Item),
	}
}

var ErrKeyNotFound = errors.New("key not found")
var ErrKeyExpired = errors.New("key is expired")

func (i *InMemoryStorage) Set(key string, value []byte, expiryDuration time.Duration) {

	var expireAt int64
	expireAt = DEFAULT

	if expiryDuration > 0 {
		expireAt = time.Now().Add(expiryDuration).UnixNano()
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	i.Data[key] = Item{Value: value, ExpireAt: expireAt}
}

func (i *InMemoryStorage) Get(key string) (Item, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	item, ok := i.Data[key]

	if !ok {
		return Item{}, ErrKeyNotFound
	}

	if item.ExpireAt > 0 && item.ExpireAt < time.Now().UnixNano() {
		i.Delete([]string{key})
		return Item{}, ErrKeyExpired
	}

	return item, nil
}

func (i *InMemoryStorage) Delete(keys []string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	for j := 0; j < len(keys); j++ {
		delete(i.Data, keys[j])
	}
}

func (i *InMemoryStorage) Keys() []string {
	i.mu.Lock()
	defer i.mu.Unlock()

	keys := make([]string, 0, len(i.Data))

	for k, _ := range i.Data {
		keys = append(keys, k)
	}
	return keys
}
