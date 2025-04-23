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

type MapMutex[K comparable, V any] struct{ baseMutex[map[K]V] }

type MapRWMutex[K comparable, V any] struct{ baseRWMutex[map[K]V] }

type SliceMutex[T any] struct{ baseMutex[[]T] }

type SliceRWMutex[T any] struct{ baseRWMutex[[]T] }

type NumberMutex[T INumber] struct{ baseMutex[T] }

type NumberRWMutex[T INumber] struct{ baseRWMutex[T] }

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

func (s *SliceMutex[T]) Each(clb func(T))             { each(s, clb) }
func (s *SliceMutex[T]) Clear()                       { sliceClear(s) }
func (s *SliceMutex[T]) Append(els ...T)              { sliceAppend(s, els...) }
func (s *SliceMutex[T]) Unshift(el T)                 { unshift(s, el) }
func (s *SliceMutex[T]) Shift() T                     { return shift(s) }
func (s *SliceMutex[T]) Pop() T                       { return pop(s) }
func (s *SliceMutex[T]) Clone() []T                   { return clone(s) }
func (s *SliceMutex[T]) Len() int                     { return sliceLen(s) }
func (s *SliceMutex[T]) IsEmpty() bool                { return isEmpty(s) }
func (s *SliceMutex[T]) Get(i int) T                  { return get(s, i) }
func (s *SliceMutex[T]) Remove(i int) T               { return remove(s, i) }
func (s *SliceMutex[T]) Insert(i int, el T)           { insert(s, i, el) }
func (s *SliceMutex[T]) Filter(keep func(T) bool) []T { return filter(s, keep) }

func (s *SliceRWMutex[T]) Each(clb func(T))             { each(s, clb) }
func (s *SliceRWMutex[T]) Clear()                       { sliceClear(s) }
func (s *SliceRWMutex[T]) Append(els ...T)              { sliceAppend(s, els...) }
func (s *SliceRWMutex[T]) Unshift(el T)                 { unshift(s, el) }
func (s *SliceRWMutex[T]) Shift() T                     { return shift(s) }
func (s *SliceRWMutex[T]) Pop() T                       { return pop(s) }
func (s *SliceRWMutex[T]) Clone() []T                   { return clone(s) }
func (s *SliceRWMutex[T]) Len() int                     { return sliceLen(s) }
func (s *SliceRWMutex[T]) IsEmpty() bool                { return isEmpty(s) }
func (s *SliceRWMutex[T]) Get(i int) T                  { return get(s, i) }
func (s *SliceRWMutex[T]) Remove(i int) T               { return remove(s, i) }
func (s *SliceRWMutex[T]) Insert(i int, el T)           { insert(s, i, el) }
func (s *SliceRWMutex[T]) Filter(keep func(T) bool) []T { return filter(s, keep) }

func (m *MapMutex[K, V]) Clear()                       { mapClear(m) }
func (m *MapMutex[K, V]) Insert(k K, v V)              { mapInsert(m, k, v) }
func (m *MapMutex[K, V]) Get(k K) (V, bool)            { return mapGet(m, k) }
func (m *MapMutex[K, V]) GetKeyValue(k K) (K, V, bool) { return getKeyValue(m, k) }
func (m *MapMutex[K, V]) ContainsKey(k K) bool         { return containsKey(m, k) }
func (m *MapMutex[K, V]) Remove(k K) (out V, ok bool)  { return mapRemove(m, k) }
func (m *MapMutex[K, V]) Delete(k K)                   { mapDelete(m, k) }
func (m *MapMutex[K, V]) Len() int                     { return mapLen(m) }
func (m *MapMutex[K, V]) IsEmpty() (out bool)          { return mapIsEmpty(m) }
func (m *MapMutex[K, V]) Each(clb func(K, V))          { mapEach(m, clb) }
func (m *MapMutex[K, V]) Keys() []K                    { return keys(m) }
func (m *MapMutex[K, V]) Values() []V                  { return values(m) }
func (m *MapMutex[K, V]) Clone() map[K]V               { return mapClone(m) }

func (m *MapRWMutex[K, V]) Clear()                       { mapClear(m) }
func (m *MapRWMutex[K, V]) Insert(k K, v V)              { mapInsert(m, k, v) }
func (m *MapRWMutex[K, V]) Get(k K) (V, bool)            { return mapGet(m, k) }
func (m *MapRWMutex[K, V]) GetKeyValue(k K) (K, V, bool) { return getKeyValue(m, k) }
func (m *MapRWMutex[K, V]) ContainsKey(k K) bool         { return containsKey(m, k) }
func (m *MapRWMutex[K, V]) Remove(k K) (out V, ok bool)  { return mapRemove(m, k) }
func (m *MapRWMutex[K, V]) Delete(k K)                   { mapDelete(m, k) }
func (m *MapRWMutex[K, V]) Len() int                     { return mapLen(m) }
func (m *MapRWMutex[K, V]) IsEmpty() (out bool)          { return mapIsEmpty(m) }
func (m *MapRWMutex[K, V]) Each(clb func(K, V))          { mapEach(m, clb) }
func (m *MapRWMutex[K, V]) Keys() []K                    { return keys(m) }
func (m *MapRWMutex[K, V]) Values() []V                  { return values(m) }
func (m *MapRWMutex[K, V]) Clone() map[K]V               { return mapClone(m) }

func (m *NumberMutex[T]) Add(diff T) { add(m, diff) }
func (m *NumberMutex[T]) Sub(diff T) { sub(m, diff) }

func (m *NumberRWMutex[T]) Add(diff T) { add(m, diff) }
func (m *NumberRWMutex[T]) Sub(diff T) { sub(m, diff) }
