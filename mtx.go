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

// SyncMutex alias type
type SyncMutex = sync.Mutex

// SyncRWMutex alias type
type SyncRWMutex = sync.RWMutex

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
func (m *base[M, T]) WithE(clb func(v *T) error) error { return withE(m, clb) }

// With same as WithE but do return an error
func (m *base[M, T]) With(clb func(v *T)) { with(m, clb) }

// RWithE provide a callback scope where the wrapped value can be safely used for Read only purposes
func (m *base[M, T]) RWithE(clb func(v T) error) error { return rWithE(m, clb) }

// RWith same as RWithE but do not return an error
func (m *base[M, T]) RWith(clb func(v T)) { rWith(m, clb) }

// Load safely gets the wrapped value
func (m *base[M, T]) Load() (out T) { return load(m) }

// Store a new value
func (m *base[M, T]) Store(newV T) { store(m, newV) }

// Swap set a new value and return the old value
func (m *base[M, T]) Swap(newVal T) (old T) { return swap(m, newVal) }

//-----------------------------------------------------------------------------

// generic helpers for sync.Mutex/sync.RWMutex
type mtx[T any] struct{ *base[*SyncMutex, T] }
type rwMtx[T any] struct{ *base[*SyncRWMutex, T] }

// newMtxPtr/newRWMtxPtr creates a new mtx/rwMtx
func newMtxPtr[T any](v T) *mtx[T]     { return &mtx[T]{newBase(&SyncMutex{}, v)} }
func newRWMtxPtr[T any](v T) *rwMtx[T] { return &rwMtx[T]{newBase(&SyncRWMutex{}, v)} }

// RLock exposes the underlying sync.RWMutex RLock function
func (m *rwMtx[T]) RLock() { m.m.RLock() }

// RUnlock exposes the underlying sync.RWMutex RUnlock function
func (m *rwMtx[T]) RUnlock() { m.m.RUnlock() }

// RWithE provide a callback scope where the wrapped value can be safely used for Read only purposes
func (m *rwMtx[T]) RWithE(clb func(v T) error) error { return rWithE(m, clb) }

// RWith same as RWithE but do not return an error
func (m *rwMtx[T]) RWith(clb func(v T)) { rWith(m, clb) }

//-----------------------------------------------------------------------------
// Methods for Mtx

//-----------------------------------------------------------------------------
// Methods for Map

// Clear clears the map, removing all key-value pairs
func (m *Map[K, V]) Clear() { mapClear(m) }

// Insert inserts a key/value in the map
func (m *Map[K, V]) Insert(k K, v V) { mapInsert(m, k, v) }

// Get returns the value corresponding to the key
func (m *Map[K, V]) Get(k K) (out V, ok bool) { return mapGet(m, k) }

// GetKeyValue returns the key-value pair corresponding to the supplied key.
func (m *Map[K, V]) GetKeyValue(k K) (K, V, bool) { return getKeyValue(m, k) }

// ContainsKey returns true if the map contains a value for the specified key
func (m *Map[K, V]) ContainsKey(k K) bool { return containsKey(m, k) }

// Remove if the key exists, its value is returned to the caller and the key deleted from the map
func (m *Map[K, V]) Remove(k K) (V, bool) { return mapRemove(m, k) }

// Delete deletes a key from the map
func (m *Map[K, V]) Delete(k K) { mapDelete(m, k) }

// Len returns the length of the map
func (m *Map[K, V]) Len() int { return mapLen(m) }

// IsEmpty returns true if the map contains no elements.
func (m *Map[K, V]) IsEmpty() bool { return mapIsEmpty(m) }

// Each iterates each key/value of the map
func (m *Map[K, V]) Each(clb func(K, V)) { mapEach(m, clb) }

// Keys returns a slice of all keys
func (m *Map[K, V]) Keys() []K { return keys(m) }

// Values returns a slice of all values
func (m *Map[K, V]) Values() []V { return values(m) }

// Clone returns a clone of the map
func (m *Map[K, V]) Clone() map[K]V { return mapClone(m) }

//-----------------------------------------------------------------------------
// Methods for Slice

// Each iterates each values of the slice
func (s *Slice[T]) Each(clb func(T)) { each(s, clb) }

// Clear clears the slice, removing all values
func (s *Slice[T]) Clear() { sliceClear(s) }

// Append appends elements at the end of the slice
func (s *Slice[T]) Append(els ...T) { sliceAppend(s, els...) }

// Unshift insert new element at beginning of the slice
func (s *Slice[T]) Unshift(el T) { unshift(s, el) }

// Shift (pop front) remove and return the first element from the slice
func (s *Slice[T]) Shift() T { return shift(s) }

// Pop remove and return the last element from the slice
func (s *Slice[T]) Pop() T { return pop(s) }

// Clone returns a clone of the slice
func (s *Slice[T]) Clone() []T { return clone(s) }

// Len returns the length of the slice
func (s *Slice[T]) Len() int { return sliceLen(s) }

// IsEmpty returns true if the map contains no elements.
func (s *Slice[T]) IsEmpty() bool { return isEmpty(s) }

// Get gets the element at index i
func (s *Slice[T]) Get(i int) T { return get(s, i) }

// Remove removes the element at position i within the slice,
// shifting all elements after it to the left
// Panics if index is out of bounds
func (s *Slice[T]) Remove(i int) T { return remove(s, i) }

// Insert insert a new element at index i
func (s *Slice[T]) Insert(i int, el T) { insert(s, i, el) }

// Filter returns a new slice of the elements that satisfy the "keep" predicate callback
func (s *Slice[T]) Filter(keep func(T) bool) []T { return filter(s, keep) }

//-----------------------------------------------------------------------------
// Methods for Number

// Add adds "diff" to the protected number
func (n *Number[T]) Add(diff T) { add(n, diff) }

// Sub subtract "diff" to the protected number
func (n *Number[T]) Sub(diff T) { sub(n, diff) }

//-----------------------------------------------------------------------------
// Generic functions

func withE[M Locker[T], T any](m M, clb func(v *T) error) error {
	m.Lock()
	defer m.Unlock()
	return clb(m.GetPointer())
}

func with[M Locker[T], T any](m M, clb func(v *T)) {
	_ = m.WithE(func(tx *T) error {
		clb(tx)
		return nil
	})
}

func rWithE[M Locker[T], T any](m M, clb func(v T) error) error {
	return m.WithE(func(v *T) error {
		return clb(*v)
	})
}

func rWith[M Locker[T], T any](m M, clb func(v T)) {
	_ = m.RWithE(func(tx T) error {
		clb(tx)
		return nil
	})
}

func load[M Locker[T], T any](m M) (out T) {
	m.RWith(func(v T) { out = v })
	return out
}

func store[M Locker[T], T any](m M, newV T) {
	m.With(func(v *T) { *v = newV })
}

func swap[M Locker[T], T any](m M, newVal T) (old T) {
	m.With(func(v *T) {
		old = *v
		*v = newVal
	})
	return
}

func each[M Locker[T], T []E, E any](m M, clb func(E)) {
	m.RWith(func(v T) {
		for _, e := range v {
			clb(e)
		}
	})
}

func sliceClear[M Locker[T], T []E, E any](m M) {
	m.With(func(v *T) { *v = nil; *v = make([]E, 0) })
}

func sliceAppend[M Locker[T], T []E, E any](m M, els ...E) {
	m.With(func(v *T) { *v = append(*v, els...) })
}

func unshift[M Locker[T], T []E, E any](m M, el E) {
	m.With(func(v *T) { *v = append([]E{el}, *v...) })
}

func shift[M Locker[T], T []E, E any](m M) (out E) {
	m.With(func(v *T) { out, *v = (*v)[0], (*v)[1:] })
	return
}

func pop[M Locker[T], T []E, E any](m M) (out E) {
	m.With(func(v *T) { out, *v = (*v)[len(*v)-1], (*v)[:len(*v)-1] })
	return
}

func clone[M Locker[T], T []E, E any](m M) (out []E) {
	m.RWith(func(v T) {
		out = make([]E, len(v))
		copy(out, v)
	})
	return
}

func sliceLen[M Locker[T], T []E, E any](m M) (out int) {
	m.RWith(func(v T) { out = len(v) })
	return
}

func isEmpty[M Locker[T], T []E, E any](m M) (out bool) {
	m.RWith(func(v T) { out = len(v) == 0 })
	return
}

func get[M Locker[T], T []E, E any](m M, i int) (out E) {
	m.RWith(func(v T) { out = (v)[i] })
	return
}

func remove[M Locker[T], T []E, E any](m M, i int) (out E) {
	m.With(func(v *T) {
		out = (*v)[i]
		*v = (*v)[:i+copy((*v)[i:], (*v)[i+1:])]
	})
	return
}

func insert[M Locker[T], T []E, E any](m M, i int, el E) {
	m.With(func(v *T) {
		var zero E
		*v = append(*v, zero)
		copy((*v)[i+1:], (*v)[i:])
		(*v)[i] = el
	})
}

func filter[M Locker[T], T []E, E any](m M, keep func(el E) bool) (out []E) {
	m.RWith(func(v T) {
		out = make([]E, 0)
		for _, x := range v {
			if keep(x) {
				out = append(out, x)
			}
		}
	})
	return
}

func mapClear[M Locker[T], T map[K]V, K comparable, V any](m M) {
	m.With(func(m *T) { clear(*m) })
}

func mapInsert[M Locker[T], T map[K]V, K comparable, V any](m M, k K, v V) {
	m.With(func(m *T) { (*m)[k] = v })
}

func mapGet[M Locker[T], T map[K]V, K comparable, V any](m M, k K) (out V, ok bool) {
	m.RWith(func(mm T) { out, ok = mm[k] })
	return
}

func getKeyValue[M Locker[T], T map[K]V, K comparable, V any](m M, k K) (key K, value V, ok bool) {
	m.RWith(func(mm T) { value, ok = mm[k] })
	if ok {
		return k, value, true
	}
	return
}

func containsKey[M Locker[T], T map[K]V, K comparable, V any](m M, k K) (found bool) {
	m.RWith(func(mm T) { _, found = mm[k] })
	return
}

func mapRemove[M Locker[T], T map[K]V, K comparable, V any](m M, k K) (out V, ok bool) {
	m.With(func(m *T) {
		out, ok = (*m)[k]
		if ok {
			delete(*m, k)
		}
	})
	return
}

func mapDelete[M Locker[T], T map[K]V, K comparable, V any](m M, k K) {
	m.With(func(m *T) { delete(*m, k) })
	return
}

func mapLen[M Locker[T], T map[K]V, K comparable, V any](m M) (out int) {
	m.RWith(func(mm T) { out = len(mm) })
	return
}

func mapIsEmpty[M Locker[T], T map[K]V, K comparable, V any](m M) (out bool) {
	m.RWith(func(mm T) { out = len(mm) == 0 })
	return
}

func mapEach[M Locker[T], T map[K]V, K comparable, V any](m M, clb func(K, V)) {
	m.RWith(func(mm T) {
		for k, v := range mm {
			clb(k, v)
		}
	})
}

func keys[M Locker[T], T map[K]V, K comparable, V any](m M) (out []K) {
	out = make([]K, 0)
	m.RWith(func(mm T) {
		for k := range mm {
			out = append(out, k)
		}
	})
	return
}

func values[M Locker[T], T map[K]V, K comparable, V any](m M) (out []V) {
	out = make([]V, 0)
	m.RWith(func(mm T) {
		for _, v := range mm {
			out = append(out, v)
		}
	})
	return
}

func mapClone[M Locker[T], T map[K]V, K comparable, V any](m M) (out map[K]V) {
	m.RWith(func(mm T) {
		out = make(map[K]V, len(mm))
		for k, v := range mm {
			out[k] = v
		}
	})
	return
}

func add[M Locker[T], T INumber](m M, diff T) { m.With(func(v *T) { *v += diff }) }

func sub[M Locker[T], T INumber](m M, diff T) { m.With(func(v *T) { *v -= diff }) }
