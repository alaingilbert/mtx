// MIT License
//
// Copyright (c) 2024 Alain Gilbert
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package mtx

import "sync"

// Mutex alias type
type Mutex = sync.Mutex

// RWMutex alias type
type RWMutex = sync.RWMutex

func toPtr[T any](v T) *T { return &v }

func first[T any](a T, _ ...any) T { return a }

// returns a default empty map if v is nil
func defaultMap[K comparable, V any](v map[K]V) map[K]V {
	if v == nil {
		v = make(map[K]V)
	}
	return v
}

// returns a default empty slice if v is nil
func defaultSlice[T any](v []T) []T {
	if v == nil {
		v = make([]T, 0)
	}
	return v
}

//-----------------------------------------------------------------------------
// Interfaces

// Locker is the interface that each mtx types implements (Mtx/Map/Slice/Number)
type Locker[T any] interface {
	sync.Locker
	GetPointer() *T
	Load() T
	RLock()
	RUnlock()
	RWith(clb func(v T))
	RWithE(clb func(v T) error) error
	Store(v T)
	Swap(newVal T) (old T)
	With(clb func(v *T))
	WithE(clb func(v *T) error) error
}

// IMap is the interface that Map implements
type IMap[K comparable, V any] interface {
	Locker[map[K]V]
	Clear()
	Clone() (out map[K]V)
	ContainsKey(k K) (found bool)
	Delete(k K)
	Each(clb func(K, V))
	Get(k K) (out V, ok bool)
	GetKeyValue(k K) (key K, value V, ok bool)
	Insert(k K, v V)
	IsEmpty() bool
	Keys() (out []K)
	Len() (out int)
	Remove(k K) (out V, ok bool)
	Values() (out []V)
}

// ISlice is the interface that Slice implements
type ISlice[T any] interface {
	Locker[[]T]
	Append(els ...T)
	Clear()
	Clone() (out []T)
	Each(clb func(T))
	Filter(func(T) bool) []T
	Get(i int) (out T)
	Insert(i int, el T)
	IsEmpty() bool
	Len() (out int)
	Pop() (out T)
	Remove(i int) (out T)
	Shift() (out T)
	Unshift(el T)
}

// INumber all numbers
type INumber interface {
	~float32 | ~float64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~complex64 | ~complex128
}

//-----------------------------------------------------------------------------
// Types

// Mtx mutex protected value
type Mtx[T any] struct{ Locker[T] }

// Map mutex protected map
type Map[K comparable, V any] struct{ Locker[map[K]V] }

// Slice mutex protected slice
type Slice[V any] struct{ Locker[[]V] }

// Number mutex protected number
type Number[T INumber] struct{ Locker[T] }

type base[M sync.Locker, T any] struct {
	m M
	v T
}

// Compile time checks to ensure types satisfies interfaces
var _ Locker[any] = (*Mtx[any])(nil)
var _ Locker[int] = (*Number[int])(nil)
var _ IMap[int, int] = (*Map[int, int])(nil)
var _ ISlice[any] = (*Slice[any])(nil)
var _ Locker[any] = (*base[sync.Locker, any])(nil)

//-----------------------------------------------------------------------------
// Constructors

// NewMtx returns a new Mtx with a sync.Mutex as backend
func NewMtx[T any](v T) Mtx[T] { return Mtx[T]{newMtxPtr(v)} }

// NewRWMtx returns a new Mtx with a sync.RWMutex as backend
func NewRWMtx[T any](v T) Mtx[T] { return Mtx[T]{newRWMtxPtr(v)} }

// NewNumber returns a new Number with a sync.Mutex as backend
func NewNumber[T INumber](v T) Number[T] { return Number[T]{newMtxPtr(v)} }

// NewRWNumber returns a new Number with a sync.RWMutex as backend
func NewRWNumber[T INumber](v T) Number[T] { return Number[T]{newRWMtxPtr(v)} }

// NewMap returns a new Map with a sync.Mutex as backend
func NewMap[K comparable, V any](v map[K]V) Map[K, V] { return Map[K, V]{newMtxPtr(defaultMap(v))} }

// NewRWMap returns a new Map with a sync.RWMutex as backend
func NewRWMap[K comparable, V any](v map[K]V) Map[K, V] { return Map[K, V]{newRWMtxPtr(defaultMap(v))} }

// NewSlice returns a new Slice with a sync.Mutex as backend
func NewSlice[T any](v []T) Slice[T] { return Slice[T]{newMtxPtr(defaultSlice(v))} }

// NewRWSlice returns a new Slice with a sync.RWMutex as backend
func NewRWSlice[T any](v []T) Slice[T] { return Slice[T]{newRWMtxPtr(defaultSlice(v))} }

// NewMtxPtr same as NewMtx, but as a pointer
func NewMtxPtr[T any](v T) *Mtx[T] { return toPtr(NewMtx(v)) }

// NewRWMtxPtr same as Mtx, but as a pointer
func NewRWMtxPtr[T any](v T) *Mtx[T] { return toPtr(NewRWMtx(v)) }

// NewNumberPtr same as NewNumber, but as a pointer
func NewNumberPtr[T INumber](v T) *Number[T] { return toPtr(NewNumber(v)) }

// NewRWNumberPtr same as NewRWNumber, but as a pointer
func NewRWNumberPtr[T INumber](v T) *Number[T] { return toPtr(NewRWNumber(v)) }

// NewMapPtr same as NewMap, but as a pointer
func NewMapPtr[K comparable, V any](v map[K]V) *Map[K, V] { return toPtr(NewMap(v)) }

// NewRWMapPtr same as NewRWMap, but as a pointer
func NewRWMapPtr[K comparable, V any](v map[K]V) *Map[K, V] { return toPtr(NewRWMap(v)) }

// NewSlicePtr same as NewSlice, but as a pointer
func NewSlicePtr[T any](v []T) *Slice[T] { return toPtr(NewSlice(v)) }

// NewRWSlicePtr same as NewRWSlice, but as a pointer
func NewRWSlicePtr[T any](v []T) *Slice[T] { return toPtr(NewRWSlice(v)) }

//-----------------------------------------------------------------------------
// Base implementation

func newBase[M sync.Locker, T any](m M, v T) *base[M, T] { return &base[M, T]{m, v} }

// Lock exposes the underlying sync.Mutex Lock function
func (m *base[M, T]) Lock() { m.m.Lock() }

// Unlock exposes the underlying sync.Mutex Unlock function
func (m *base[M, T]) Unlock() { m.m.Unlock() }

// RLock is a default implementation of RLock to satisfy Locker interface
func (m *base[M, T]) RLock() { m.Lock() }

// RUnlock is a default implementation of RUnlock to satisfy Locker interface
func (m *base[M, T]) RUnlock() { m.Unlock() }

// GetPointer returns a pointer to the protected value
// WARNING: the caller must make sure the code that uses the returned pointer is thread-safe
func (m *base[M, T]) GetPointer() *T { return &m.v }

// WithE provide a callback scope where the wrapped value can be safely used
func (m *base[M, T]) WithE(clb func(v *T) error) error {
	m.Lock()
	defer m.Unlock()
	return clb(&m.v)
}

// With same as WithE but do return an error
func (m *base[M, T]) With(clb func(v *T)) {
	_ = m.WithE(func(tx *T) error {
		clb(tx)
		return nil
	})
}

// RWithE provide a callback scope where the wrapped value can be safely used for Read only purposes
func (m *base[M, T]) RWithE(clb func(v T) error) error {
	return m.WithE(func(v *T) error {
		return clb(*v)
	})
}

// RWith same as RWithE but do not return an error
func (m *base[M, T]) RWith(clb func(v T)) {
	_ = m.RWithE(func(tx T) error {
		clb(tx)
		return nil
	})
}

// Load safely gets the wrapped value
func (m *base[M, T]) Load() (out T) {
	m.RWith(func(v T) { out = v })
	return out
}

// Store a new value
func (m *base[M, T]) Store(newV T) {
	m.With(func(v *T) { *v = newV })
}

// Swap set a new value and return the old value
func (m *base[M, T]) Swap(newVal T) (old T) {
	m.With(func(v *T) {
		old = *v
		*v = newVal
	})
	return
}

//-----------------------------------------------------------------------------

// generic helpers for sync.Mutex/sync.RWMutex
type mtx[T any] struct{ *base[*Mutex, T] }
type rwMtx[T any] struct{ *base[*RWMutex, T] }

// newMtxPtr/newRWMtxPtr creates a new mtx/rwMtx
func newMtxPtr[T any](v T) *mtx[T]     { return &mtx[T]{newBase(&Mutex{}, v)} }
func newRWMtxPtr[T any](v T) *rwMtx[T] { return &rwMtx[T]{newBase(&RWMutex{}, v)} }

// RLock exposes the underlying sync.RWMutex RLock function
func (m *rwMtx[T]) RLock() { m.m.RLock() }

// RUnlock exposes the underlying sync.RWMutex RUnlock function
func (m *rwMtx[T]) RUnlock() { m.m.RUnlock() }

// RWithE provide a callback scope where the wrapped value can be safely used for Read only purposes
func (m *rwMtx[T]) RWithE(clb func(v T) error) error {
	m.RLock()
	defer m.RUnlock()
	return clb(m.v)
}

// RWith same as RWithE but do not return an error
func (m *rwMtx[T]) RWith(clb func(v T)) {
	_ = m.RWithE(func(tx T) error {
		clb(tx)
		return nil
	})
}

//-----------------------------------------------------------------------------
// Methods for Mtx

//-----------------------------------------------------------------------------
// Methods for Map

// Clear clears the map, removing all key-value pairs
func (m *Map[K, V]) Clear() {
	m.With(func(m *map[K]V) { clear(*m) })
}

// Insert inserts a key/value in the map
func (m *Map[K, V]) Insert(k K, v V) {
	m.With(func(m *map[K]V) { (*m)[k] = v })
}

// Get returns the value corresponding to the key
func (m *Map[K, V]) Get(k K) (out V, ok bool) {
	m.RWith(func(mm map[K]V) { out, ok = mm[k] })
	return
}

// GetKeyValue returns the key-value pair corresponding to the supplied key.
func (m *Map[K, V]) GetKeyValue(k K) (key K, value V, ok bool) {
	m.RWith(func(mm map[K]V) { value, ok = mm[k] })
	if ok {
		return k, value, true
	}
	return
}

// ContainsKey returns true if the map contains a value for the specified key
func (m *Map[K, V]) ContainsKey(k K) (found bool) {
	m.RWith(func(mm map[K]V) { _, found = mm[k] })
	return
}

// Remove if the key exists, its value is returned to the caller and the key deleted from the map
func (m *Map[K, V]) Remove(k K) (out V, ok bool) {
	m.With(func(m *map[K]V) {
		out, ok = (*m)[k]
		if ok {
			delete(*m, k)
		}
	})
	return
}

// Delete deletes a key from the map
func (m *Map[K, V]) Delete(k K) {
	m.With(func(m *map[K]V) { delete(*m, k) })
	return
}

// Len returns the length of the map
func (m *Map[K, V]) Len() (out int) {
	m.RWith(func(mm map[K]V) { out = len(mm) })
	return
}

// IsEmpty returns true if the map contains no elements.
func (m *Map[K, V]) IsEmpty() (out bool) {
	m.RWith(func(mm map[K]V) { out = len(mm) == 0 })
	return
}

// Each iterates each key/value of the map
func (m *Map[K, V]) Each(clb func(K, V)) {
	m.RWith(func(mm map[K]V) {
		for k, v := range mm {
			clb(k, v)
		}
	})
}

// Keys returns a slice of all keys
func (m *Map[K, V]) Keys() (out []K) {
	out = make([]K, 0)
	m.RWith(func(mm map[K]V) {
		for k := range mm {
			out = append(out, k)
		}
	})
	return
}

// Values returns a slice of all values
func (m *Map[K, V]) Values() (out []V) {
	out = make([]V, 0)
	m.RWith(func(mm map[K]V) {
		for _, v := range mm {
			out = append(out, v)
		}
	})
	return
}

// Clone returns a clone of the map
func (m *Map[K, V]) Clone() (out map[K]V) {
	m.RWith(func(mm map[K]V) {
		out = make(map[K]V, len(mm))
		for k, v := range mm {
			out[k] = v
		}
	})
	return
}

//-----------------------------------------------------------------------------
// Methods for Slice

// Each iterates each values of the slice
func (s *Slice[T]) Each(clb func(T)) {
	s.RWith(func(v []T) {
		for _, e := range v {
			clb(e)
		}
	})
}

// Clear clears the slice, removing all values
func (s *Slice[T]) Clear() {
	s.With(func(v *[]T) { *v = nil; *v = make([]T, 0) })
}

// Append appends elements at the end of the slice
func (s *Slice[T]) Append(els ...T) {
	s.With(func(v *[]T) { *v = append(*v, els...) })
}

// Unshift insert new element at beginning of the slice
func (s *Slice[T]) Unshift(el T) {
	s.With(func(v *[]T) { *v = append([]T{el}, *v...) })
}

// Shift (pop front) remove and return the first element from the slice
func (s *Slice[T]) Shift() (out T) {
	s.With(func(v *[]T) { out, *v = (*v)[0], (*v)[1:] })
	return
}

// Pop remove and return the last element from the slice
func (s *Slice[T]) Pop() (out T) {
	s.With(func(v *[]T) { out, *v = (*v)[len(*v)-1], (*v)[:len(*v)-1] })
	return
}

// Clone returns a clone of the slice
func (s *Slice[T]) Clone() (out []T) {
	s.RWith(func(v []T) {
		out = make([]T, len(v))
		copy(out, v)
	})
	return
}

// Len returns the length of the slice
func (s *Slice[T]) Len() (out int) {
	s.RWith(func(v []T) { out = len(v) })
	return
}

// IsEmpty returns true if the map contains no elements.
func (s *Slice[T]) IsEmpty() (out bool) {
	s.RWith(func(v []T) { out = len(v) == 0 })
	return
}

// Get gets the element at index i
func (s *Slice[T]) Get(i int) (out T) {
	s.RWith(func(v []T) { out = (v)[i] })
	return
}

// Remove removes the element at position i within the slice,
// shifting all elements after it to the left
// Panics if index is out of bounds
func (s *Slice[T]) Remove(i int) (out T) {
	s.With(func(v *[]T) { out, *v = (*v)[i], (*v)[:i+copy((*v)[i:], (*v)[i+1:])] })
	return
}

// Insert insert a new element at index i
func (s *Slice[T]) Insert(i int, el T) {
	s.With(func(v *[]T) {
		var zero T
		*v = append(*v, zero)
		copy((*v)[i+1:], (*v)[i:])
		(*v)[i] = el
	})
}

// Filter returns a new slice of the elements that satisfy the "keep" predicate callback
func (s *Slice[T]) Filter(keep func(el T) bool) (out []T) {
	s.RWith(func(v []T) {
		out = make([]T, 0)
		for _, x := range v {
			if keep(x) {
				out = append(out, x)
			}
		}
	})
	return
}

//-----------------------------------------------------------------------------
// Methods for Number

// Add adds "diff" to the protected number
func (n *Number[T]) Add(diff T) { n.With(func(v *T) { *v += diff }) }

// Sub subtract "diff" to the protected number
func (n *Number[T]) Sub(diff T) { n.With(func(v *T) { *v -= diff }) }
