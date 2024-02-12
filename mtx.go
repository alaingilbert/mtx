package mtx

import (
	"cmp"
	"sync"
)

// Mtx generic helper for sync.Mutex
type Mtx[T any] struct {
	sync.Mutex
	v T
}

// NewMtx creates a new Mtx
func NewMtx[T any](v T) Mtx[T] {
	return Mtx[T]{v: v}
}

// NewMtxPtr creates a new pointer to *Mtx
func NewMtxPtr[T any](v T) *Mtx[T] {
	return &Mtx[T]{v: v}
}

// Val gets the wrapped value by the mutex.
// WARNING: the caller must make sure the code that uses it is thread-safe
func (m *Mtx[T]) Val() *T {
	return &m.v
}

// Get safely gets the wrapped value
func (m *Mtx[T]) Get() T {
	m.Lock()
	defer m.Unlock()
	return m.v
}

// Set a new value
func (m *Mtx[T]) Set(v T) {
	m.Lock()
	defer m.Unlock()
	m.v = v
}

// WithE provide a callback scope where the wrapped value can be safely used
func (m *Mtx[T]) WithE(clb func(v *T) error) error {
	m.Lock()
	defer m.Unlock()
	return clb(&m.v)
}

// With same as WithE but do return an error
func (m *Mtx[T]) With(clb func(v *T)) {
	_ = m.WithE(func(tx *T) error {
		clb(tx)
		return nil
	})
}

//----------------------

// RWMtx generic helper for sync.RWMutex
type RWMtx[T any] struct {
	sync.RWMutex
	v T
}

// NewRWMtx creates a new RWMtx
func NewRWMtx[T any](v T) RWMtx[T] {
	return RWMtx[T]{v: v}
}

// NewRWMtxPtr creates a new pointer to *RWMtx
func NewRWMtxPtr[T any](v T) *RWMtx[T] {
	return &RWMtx[T]{v: v}
}

// Val gets the wrapped value by the mutex.
// WARNING: the caller must make sure the code that uses it is thread-safe
func (m *RWMtx[T]) Val() *T {
	return &m.v
}

// Get safely gets the wrapped value using the Read part of the Read-Write mutex
func (m *RWMtx[T]) Get() T {
	m.RLock()
	defer m.RUnlock()
	return m.v
}

// Set a new value using the Write part of the Read-Write mutex
func (m *RWMtx[T]) Set(v T) {
	m.Lock()
	defer m.Unlock()
	m.v = v
}

// RWithE provide a callback scope where the wrapped value can be safely used for Read only purposes
func (m *RWMtx[T]) RWithE(clb func(v T) error) error {
	m.RLock()
	defer m.RUnlock()
	return clb(m.v)
}

// WithE provide a callback scope where the wrapped value can be safely used
func (m *RWMtx[T]) WithE(clb func(v *T) error) error {
	m.Lock()
	defer m.Unlock()
	return clb(&m.v)
}

// RWith same as RWithE but do not return an error
func (m *RWMtx[T]) RWith(clb func(v T)) {
	_ = m.RWithE(func(tx T) error {
		clb(tx)
		return nil
	})
}

// With same as WithE but do return an error
func (m *RWMtx[T]) With(clb func(v *T)) {
	_ = m.WithE(func(tx *T) error {
		clb(tx)
		return nil
	})
}

// Replace set a new value and return the old value
func (m *RWMtx[T]) Replace(newVal T) (old T) {
	m.With(func(v *T) {
		old = *v
		*v = newVal
	})
	return
}

//----------------------

type RWMtxMap[K cmp.Ordered, V any] struct {
	RWMtx[map[K]V]
}

func NewMap[K cmp.Ordered, V any]() RWMtxMap[K, V] {
	return RWMtxMap[K, V]{RWMtx: NewRWMtx(make(map[K]V))}
}

func NewMapPtr[K cmp.Ordered, V any]() *RWMtxMap[K, V] {
	return &RWMtxMap[K, V]{RWMtx: NewRWMtx(make(map[K]V))}
}

func (m *RWMtxMap[K, V]) SetKey(k K, v V) {
	m.With(func(m *map[K]V) { (*m)[k] = v })
}

func (m *RWMtxMap[K, V]) GetKey(k K) (out V, ok bool) {
	m.RWith(func(m map[K]V) { out, ok = m[k] })
	return
}

func (m *RWMtxMap[K, V]) HasKey(k K) (found bool) {
	m.RWith(func(m map[K]V) { _, found = m[k] })
	return
}

func (m *RWMtxMap[K, V]) TakeKey(k K) (out V, ok bool) {
	m.With(func(m *map[K]V) {
		out, ok = (*m)[k]
		if ok {
			delete(*m, k)
		}
	})
	return
}

func (m *RWMtxMap[K, V]) DeleteKey(k K) {
	m.With(func(m *map[K]V) { delete(*m, k) })
	return
}

func (m *RWMtxMap[K, V]) Len() (out int) {
	m.With(func(m *map[K]V) { out = len(*m) })
	return
}

func (m *RWMtxMap[K, V]) Each(clb func(K, V)) {
	m.RWith(func(m map[K]V) {
		for k, v := range m {
			clb(k, v)
		}
	})
}

func (m *RWMtxMap[K, V]) Keys() (out []K) {
	m.RWith(func(m map[K]V) {
		for k := range m {
			out = append(out, k)
		}
	})
	return
}

func (m *RWMtxMap[K, V]) Values() (out []V) {
	m.RWith(func(m map[K]V) {
		for _, v := range m {
			out = append(out, v)
		}
	})
	return
}

func (m *RWMtxMap[K, V]) Clone() (out map[K]V) {
	m.RWith(func(m map[K]V) {
		out = make(map[K]V, len(m))
		for k, v := range m {
			out[k] = v
		}
	})
	return
}

//----------------------

type RWMtxSlice[T any] struct {
	RWMtx[[]T]
}

func NewSlicePtr[T any]() *RWMtxSlice[T] {
	return &RWMtxSlice[T]{RWMtx: NewRWMtx(make([]T, 0))}
}

func (s *RWMtxSlice[T]) Each(clb func(T)) {
	s.RWith(func(v []T) {
		for _, e := range v {
			clb(e)
		}
	})
}

func (s *RWMtxSlice[T]) Append(els ...T) {
	s.With(func(v *[]T) { *v = append(*v, els...) })
}

// Unshift insert new element at beginning of the slice
func (s *RWMtxSlice[T]) Unshift(el T) {
	s.With(func(v *[]T) { *v = append([]T{el}, *v...) })
}

// Shift (pop front)
func (s *RWMtxSlice[T]) Shift() (out T) {
	s.With(func(v *[]T) { out, *v = (*v)[0], (*v)[1:] })
	return
}

func (s *RWMtxSlice[T]) Pop() (out T) {
	s.With(func(v *[]T) { out, *v = (*v)[len(*v)-1], (*v)[:len(*v)-1] })
	return
}

func (s *RWMtxSlice[T]) Clone() (out []T) {
	s.RWith(func(v []T) {
		out = make([]T, len(v))
		copy(out, v)
	})
	return
}

func (s *RWMtxSlice[T]) Len() (out int) {
	s.With(func(v *[]T) { out = len(*v) })
	return
}

func (s *RWMtxSlice[T]) GetIdx(i int) (out T) {
	s.With(func(v *[]T) { out = (*v)[i] })
	return
}

func (s *RWMtxSlice[T]) DeleteIdx(i int) {
	s.With(func(v *[]T) { *v = (*v)[:i+copy((*v)[i:], (*v)[i+1:])] })
}

func (s *RWMtxSlice[T]) Insert(i int, el T) {
	s.With(func(v *[]T) {
		var zero T
		*v = append(*v, zero)
		copy((*v)[i+1:], (*v)[i:])
		(*v)[i] = el
	})
}

//----------------------

type RWMtxUInt64[T ~uint64] struct {
	RWMtx[T]
}

func (s *RWMtxUInt64[T]) Incr(diff T) {
	s.With(func(v *T) { *v += diff })
}

func (s *RWMtxUInt64[T]) Decr(diff T) {
	s.With(func(v *T) { *v -= diff })
}
