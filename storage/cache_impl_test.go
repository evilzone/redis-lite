package storage

import (
	"errors"
	"testing"
	"time"
)

func TestNewInMemoryStorage(t *testing.T) {
	cache := NewInMemoryStorage()

	t.Run("get key doesn't exist returns error ", func(t *testing.T) {
		_, err := cache.Get("key1")

		if !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("Expected error to be ErrKeyNotFound but found %v", err)
		}
	})

	t.Run("delete key removes key ", func(t *testing.T) {
		cache.Set("test", []byte("val"), time.Duration(2000))
		cache.Delete([]string{"test"})
		_, err := cache.Get("test")

		if !errors.Is(err, ErrKeyNotFound) {
			t.Errorf("Expected error to be ErrKeyNotFound but found %v", err)
		}
	})
}
