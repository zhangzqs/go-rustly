package gorustly

import (
	"sync"
)

type Mutex[T any] struct {
	// data is the data protected by the mutex
	data *T
	// lock is the lock for the mutex
	lock sync.Mutex
}

// MutexNew creates a new mutex
func MutexNew[T any](data *T) Mutex[T] {
	return Mutex[T]{data: data}
}

// MutexLock locks the mutex
func (m *Mutex[T]) Lock() *T {
	m.lock.Lock()
	return m.data
}

// MutexUnlock unlocks the mutex
func (m *Mutex[T]) Drop() {
	m.lock.Unlock()
}

// MutexTryLock tries to lock the mutex
func (m *Mutex[T]) TryLock() Option[*T] {
	if m.lock.TryLock() {
		return Some[*T](m.data)
	}
	return None[*T]()
}
