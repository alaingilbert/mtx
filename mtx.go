package mtx

import (
	"cmp"
	"os"
	"strconv"
	"sync"
)

var debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

func toPtr[T any](v T) *T { return &v }

type Locker[T any] interface {
	sync.Locker
	Get() T
	Set(v T)
	Val() *T
	With(clb func(v *T))
	WithE(clb func(v *T) error) error
	RWith(clb func(v T))
	RWithE(clb func(v T) error) error
}

type Base[M sync.Locker, T any] struct {
	m M
	v T
}

func NewBase[M sync.Locker, T any](m M, v T) *Base[M, T] {
	return &Base[M, T]{m: m, v: v}
}

// Lock exposes the underlying sync.Mutex Lock function
func (m *Base[M, T]) Lock() { m.m.Lock() }

// Unlock exposes the underlying sync.Mutex Unlock function
func (m *Base[M, T]) Unlock() { m.m.Unlock() }

// Val gets the wrapped value by the mutex.
// WARNING: the caller must make sure the code that uses it is thread-safe
func (m *Base[M, T]) Val() *T {
	return &m.v
}

// WithE provide a callback scope where the wrapped value can be safely used
func (m *Base[M, T]) WithE(clb func(v *T) error) error {
	m.m.Lock()
	defer m.m.Unlock()
	return clb(&m.v)
}

// With same as WithE but do return an error
func (m *Base[M, T]) With(clb func(v *T)) {
	_ = m.WithE(func(tx *T) error {
		clb(tx)
		return nil
	})
}

// RWithE provide a callback scope where the wrapped value can be safely used for Read only purposes
func (m *Base[M, T]) RWithE(clb func(v T) error) error {
	if debug {
		println("Base RWithE")
	}
	return m.WithE(func(v *T) error {
		return clb(*v)
	})
}

// RWith same as RWithE but do not return an error
func (m *Base[M, T]) RWith(clb func(v T)) {
	_ = m.RWithE(func(tx T) error {
		clb(tx)
		return nil
	})
}

// Get safely gets the wrapped value
func (m *Base[M, T]) Get() (out T) {
	m.RWith(func(v T) { out = v })
	return out
}

// Set a new value
func (m *Base[M, T]) Set(newV T) {
	m.With(func(v *T) { *v = newV })
}

// Replace set a new value and return the old value
func (m *Base[M, T]) Replace(newVal T) (old T) {
	m.With(func(v *T) {
		old = *v
		*v = newVal
	})
	return
}

//----------------------

// Mtx generic helper for sync.Mutex
type Mtx[T any] struct {
	*Base[*sync.Mutex, T]
}

// NewMtx creates a new Mtx
func NewMtx[T any](v T) Mtx[T] {
	return Mtx[T]{NewBase[*sync.Mutex, T](&sync.Mutex{}, v)}
}

// NewMtxPtr creates a new pointer to *Mtx
func NewMtxPtr[T any](v T) *Mtx[T] { return toPtr(NewMtx(v)) }

//----------------------

// RWMtx generic helper for sync.RWMutex
type RWMtx[T any] struct {
	*Base[*sync.RWMutex, T]
}

// NewRWMtx creates a new RWMtx
func NewRWMtx[T any](v T) RWMtx[T] {
	return RWMtx[T]{NewBase[*sync.RWMutex, T](&sync.RWMutex{}, v)}
}

// NewRWMtxPtr creates a new pointer to *RWMtx
func NewRWMtxPtr[T any](v T) *RWMtx[T] { return toPtr(NewRWMtx(v)) }

// RWithE provide a callback scope where the wrapped value can be safely used for Read only purposes
func (m *RWMtx[T]) RWithE(clb func(v T) error) error {
	if debug {
		println("RWMtx RWithE")
	}
	m.m.RLock()
	defer m.m.RUnlock()
	return clb(m.v)
}

// RWith same as RWithE but do not return an error
func (m *RWMtx[T]) RWith(clb func(v T)) {
	_ = m.RWithE(func(tx T) error {
		clb(tx)
		return nil
	})
}

// RLock exposes the underlying sync.RWMutex RLock function
func (m *RWMtx[T]) RLock() { m.m.RLock() }

// RUnlock exposes the underlying sync.RWMutex RUnlock function
func (m *RWMtx[T]) RUnlock() { m.m.RUnlock() }

//----------------------

func NewMap[K cmp.Ordered, V any]() Map[K, V] {
	return Map[K, V]{newBaseMapPtr[K, V](NewMtxPtr(make(map[K]V)))}
}

func NewMapPtr[K cmp.Ordered, V any]() *Map[K, V] { return toPtr(NewMap[K, V]()) }

//----------------------

func NewRWMap[K cmp.Ordered, V any]() Map[K, V] {
	return Map[K, V]{newBaseMapPtr[K, V](NewRWMtxPtr(make(map[K]V)))}
}

func NewRWMapPtr[K cmp.Ordered, V any]() *Map[K, V] { return toPtr(NewRWMap[K, V]()) }

//----------------------

func newBaseMapPtr[K cmp.Ordered, V any](m Locker[map[K]V]) *Map[K, V] {
	return &Map[K, V]{m}
}

type Map[K cmp.Ordered, V any] struct {
	Locker[map[K]V]
}

func (m *Map[K, V]) SetKey(k K, v V) {
	m.With(func(m *map[K]V) { (*m)[k] = v })
}

func (m *Map[K, V]) GetKey(k K) (out V, ok bool) {
	m.RWith(func(mm map[K]V) { out, ok = mm[k] })
	return
}

func (m *Map[K, V]) HasKey(k K) (found bool) {
	m.RWith(func(mm map[K]V) { _, found = mm[k] })
	return
}

func (m *Map[K, V]) TakeKey(k K) (out V, ok bool) {
	m.With(func(m *map[K]V) {
		out, ok = (*m)[k]
		if ok {
			delete(*m, k)
		}
	})
	return
}

func (m *Map[K, V]) DeleteKey(k K) {
	m.With(func(m *map[K]V) { delete(*m, k) })
	return
}

func (m *Map[K, V]) Len() (out int) {
	m.RWith(func(mm map[K]V) { out = len(mm) })
	return
}

func (m *Map[K, V]) Each(clb func(K, V)) {
	m.RWith(func(mm map[K]V) {
		for k, v := range mm {
			clb(k, v)
		}
	})
}

func (m *Map[K, V]) Keys() (out []K) {
	out = make([]K, 0)
	m.RWith(func(mm map[K]V) {
		for k := range mm {
			out = append(out, k)
		}
	})
	return
}

func (m *Map[K, V]) Values() (out []V) {
	out = make([]V, 0)
	m.RWith(func(mm map[K]V) {
		for _, v := range mm {
			out = append(out, v)
		}
	})
	return
}

func (m *Map[K, V]) Clone() (out map[K]V) {
	m.RWith(func(mm map[K]V) {
		out = make(map[K]V, len(mm))
		for k, v := range mm {
			out[k] = v
		}
	})
	return
}

//----------------------

func NewSlice[V any]() Slice[V] {
	return Slice[V]{NewBaseSlicePtr[V](NewMtxPtr(make([]V, 0)))}
}

func NewSlicePtr[V any]() *Slice[V] { return toPtr(NewSlice[V]()) }

//----------------------

func NewRWSlice[V any]() Slice[V] {
	return Slice[V]{NewBaseSlicePtr[V](NewRWMtxPtr(make([]V, 0)))}
}

func NewRWSlicePtr[V any]() *Slice[V] { return toPtr(NewRWSlice[V]()) }

//----------------------

type Slice[V any] struct {
	Locker[[]V]
}

func NewBaseSlicePtr[V any](m Locker[[]V]) *Slice[V] {
	return &Slice[V]{m}
}

func (s *Slice[T]) Each(clb func(T)) {
	s.RWith(func(v []T) {
		for _, e := range v {
			clb(e)
		}
	})
}

func (s *Slice[T]) Append(els ...T) {
	s.With(func(v *[]T) { *v = append(*v, els...) })
}

// Unshift insert new element at beginning of the slice
func (s *Slice[T]) Unshift(el T) {
	s.With(func(v *[]T) { *v = append([]T{el}, *v...) })
}

// Shift (pop front)
func (s *Slice[T]) Shift() (out T) {
	s.With(func(v *[]T) { out, *v = (*v)[0], (*v)[1:] })
	return
}

func (s *Slice[T]) Pop() (out T) {
	s.With(func(v *[]T) { out, *v = (*v)[len(*v)-1], (*v)[:len(*v)-1] })
	return
}

func (s *Slice[T]) Clone() (out []T) {
	s.RWith(func(v []T) {
		out = make([]T, len(v))
		copy(out, v)
	})
	return
}

func (s *Slice[T]) Len() (out int) {
	s.RWith(func(v []T) { out = len(v) })
	return
}

func (s *Slice[T]) GetIdx(i int) (out T) {
	s.RWith(func(v []T) { out = (v)[i] })
	return
}

func (s *Slice[T]) DeleteIdx(i int) {
	s.With(func(v *[]T) { *v = (*v)[:i+copy((*v)[i:], (*v)[i+1:])] })
}

func (s *Slice[T]) Insert(i int, el T) {
	s.With(func(v *[]T) {
		var zero T
		*v = append(*v, zero)
		copy((*v)[i+1:], (*v)[i:])
		(*v)[i] = el
	})
}

//----------------------

type RWUInt64[T ~uint64] struct {
	*RWMtx[T]
}

func NewRWUInt64[T ~uint64]() RWUInt64[T] {
	return RWUInt64[T]{NewRWMtxPtr[T](0)}
}

func NewRWUInt64Ptr[T ~uint64]() *RWUInt64[T] {
	return &RWUInt64[T]{NewRWMtxPtr[T](0)}
}

func (s *RWUInt64[T]) Incr(diff T) {
	s.With(func(v *T) { *v += diff })
}

func (s *RWUInt64[T]) Decr(diff T) {
	s.With(func(v *T) { *v -= diff })
}
