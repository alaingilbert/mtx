package mtx

import "sync"

type baseMutex[T any] struct {
	m sync.Mutex
	v T
}

type baseRWMutex[T any] struct {
	m sync.RWMutex
	v T
}

type Mutex[T any] struct{ baseMutex[T] }

type RWMutex[T any] struct{ baseRWMutex[T] }

func (m *baseMutex[T]) Lock()                            { m.m.Lock() }
func (m *baseMutex[T]) Unlock()                          { m.m.Unlock() }
func (m *baseMutex[T]) RLock()                           { m.Lock() }
func (m *baseMutex[T]) RUnlock()                         { m.Unlock() }
func (m *baseMutex[T]) GetPointer() *T                   { return &m.v }
func (m *baseMutex[T]) WithE(clb func(v *T) error) error { return withE(m, clb) }
func (m *baseMutex[T]) With(clb func(v *T))              { with(m, clb) }
func (m *baseMutex[T]) RWithE(clb func(v T) error) error { return rWithE(m, clb) }
func (m *baseMutex[T]) RWith(clb func(v T))              { rWith(m, clb) }
func (m *baseMutex[T]) Load() (out T)                    { return load(m) }
func (m *baseMutex[T]) Store(newV T)                     { store(m, newV) }
func (m *baseMutex[T]) Swap(newVal T) (old T)            { return swap(m, newVal) }

func (m *baseRWMutex[T]) Lock()                            { m.m.Lock() }
func (m *baseRWMutex[T]) Unlock()                          { m.m.Unlock() }
func (m *baseRWMutex[T]) RLock()                           { m.m.RLock() }
func (m *baseRWMutex[T]) RUnlock()                         { m.m.RUnlock() }
func (m *baseRWMutex[T]) GetPointer() *T                   { return &m.v }
func (m *baseRWMutex[T]) WithE(clb func(v *T) error) error { return withE(m, clb) }
func (m *baseRWMutex[T]) With(clb func(v *T))              { with(m, clb) }
func (m *baseRWMutex[T]) RWithE(clb func(v T) error) error { return rWithE(m, clb) }
func (m *baseRWMutex[T]) RWith(clb func(v T))              { rWith(m, clb) }
func (m *baseRWMutex[T]) Load() (out T)                    { return load(m) }
func (m *baseRWMutex[T]) Store(newV T)                     { store(m, newV) }
func (m *baseRWMutex[T]) Swap(newVal T) (old T)            { return swap(m, newVal) }
