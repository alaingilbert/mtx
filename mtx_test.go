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

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"sync"
	"testing"
)

func TestMtx_LockUnlock(t *testing.T) {
	m := NewMtx("old")
	m.Lock()
	val := m.GetPointer()
	*val = "new"
	m.Unlock()
	assert.Equal(t, "new", m.Load())
}

func TestMtx_With(t *testing.T) {
	m := NewMtx("old")
	m.With(func(v *string) {
		*v = "new"
	})
	assert.Equal(t, "new", m.Load())
}

func TestMtx_RWith(t *testing.T) {
	m := NewMtx("old")
	m.RWith(func(v string) {
		assert.Equal(t, "old", v)
	})
}

func TestMtx_Store(t *testing.T) {
	m := NewMtx("old")
	assert.Equal(t, "old", m.Load())
	m.Store("new")
	assert.Equal(t, "new", m.Load())
}

func TestMtx_Swap(t *testing.T) {
	m := NewMtx("old")
	old := m.Swap("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Load())
}

func TestMtx_GetPointer(t *testing.T) {
	someString := "old"
	orig := &someString
	m := NewMtx(orig)
	val := m.GetPointer()
	**val = "new"
	assert.Equal(t, "new", someString)
	assert.Equal(t, "new", **val)
	assert.Equal(t, "new", *orig)
}

func TestMtxPtr_Swap(t *testing.T) {
	m := NewMtxPtr("old")
	old := m.Swap("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Load())
}

func TestMtx_RLockRUnlock(t *testing.T) {
	m := NewMtx("old")
	m.RLock()
	val := m.GetPointer()
	assert.Equal(t, "old", *val)
	m.RUnlock()
}

func TestRWMtx_RLockRUnlock(t *testing.T) {
	m := NewRWMtx("old")
	m.RLock()
	val := m.GetPointer()
	assert.Equal(t, "old", *val)
	m.RUnlock()
}

func TestRWMtx_Swap(t *testing.T) {
	m := NewRWMtx("old")
	old := m.Swap("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Load())
}

func TestRWMtxPtr_Swap(t *testing.T) {
	m := NewRWMtxPtr("old")
	old := m.Swap("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Load())
}

func TestRWMtx_Val(t *testing.T) {
	someString := "old"
	orig := &someString
	m := NewRWMtx(orig)
	val := m.GetPointer()
	**val = "new"
	assert.Equal(t, "new", **val)
	assert.Equal(t, "new", *orig)
}

func TestMap_Get(t *testing.T) {
	m := NewMap[string, int](nil)
	_, ok := m.Get("a")
	assert.False(t, ok)
	m.Insert("a", 1)
	el, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, el)
}

func TestMap_GetKeyValue(t *testing.T) {
	m := NewMap[string, int](nil)
	_, _, ok := m.GetKeyValue("a")
	assert.False(t, ok)
	m.Insert("a", 1)
	key, value, ok := m.GetKeyValue("a")
	assert.True(t, ok)
	assert.Equal(t, "a", key)
	assert.Equal(t, 1, value)
}

func TestMap_HasKey(t *testing.T) {
	m := NewMap[string, int](nil)
	assert.False(t, m.ContainsKey("a"))
	m.Insert("a", 1)
	assert.True(t, m.ContainsKey("a"))
	m.Delete("a")
	assert.False(t, m.ContainsKey("a"))
}

func TestMap_Remove(t *testing.T) {
	m := NewMap[string, int](nil)
	m.Insert("a", 1)
	m.Insert("b", 2)
	m.Insert("c", 3)
	assert.Equal(t, 3, m.Len())
	assert.True(t, m.ContainsKey("b"))
	val, ok := m.Remove("b")
	assert.True(t, ok)
	assert.Equal(t, 2, val)
	assert.False(t, m.ContainsKey("b"))
	_, ok = m.Remove("b")
	assert.False(t, ok)
	assert.Equal(t, 2, m.Len())
}

func TestMap_Delete(t *testing.T) {
	m := NewMap[string, int](nil)
	assert.Equal(t, 0, m.Len())
	m.Delete("a")
	m.Insert("a", 1)
	assert.Equal(t, 1, m.Len())
	m.Delete("a")
	assert.Equal(t, 0, m.Len())
}

func TestMap_Values(t *testing.T) {
	m := NewMap[string, int](nil)
	assert.Equal(t, []int{}, m.Values())
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	values := m.Values()
	slices.Sort(values)
	assert.Equal(t, []int{1, 2, 3}, values)
}

func TestMap_Keys(t *testing.T) {
	m := NewMapPtr[string, int](nil)
	assert.Equal(t, []string{}, m.Keys())
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	keys := m.Keys()
	slices.Sort(keys)
	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestMap_Each(t *testing.T) {
	m := NewMap[string, int](nil)
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	arr := make([]string, 0)
	m.Each(func(k string, v int) {
		arr = append(arr, fmt.Sprintf("%s_%d", k, v))
	})
	slices.Sort(arr)
	assert.Equal(t, []string{"a_1", "b_2", "c_3"}, arr)
}

func TestMap_Clone(t *testing.T) {
	m := NewMap[string, int](nil)
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	clonedMap := m.Clone()
	assert.Equal(t, 1, clonedMap["a"])
}

func TestRWMap_Get(t *testing.T) {
	m := NewRWMap[string, int](nil)
	_, ok := m.Get("a")
	assert.False(t, ok)
	m.Insert("a", 1)
	el, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, el)
}

func TestRWMap_HasKey(t *testing.T) {
	m := NewRWMap[string, int](nil)
	assert.False(t, m.ContainsKey("a"))
	m.Insert("a", 1)
	assert.True(t, m.ContainsKey("a"))
	m.Delete("a")
	assert.False(t, m.ContainsKey("a"))
}

func TestRWMap_Remove(t *testing.T) {
	m := NewRWMap[string, int](nil)
	m.Insert("a", 1)
	m.Insert("b", 2)
	m.Insert("c", 3)
	assert.Equal(t, 3, m.Len())
	assert.True(t, m.ContainsKey("b"))
	val, ok := m.Remove("b")
	assert.True(t, ok)
	assert.Equal(t, 2, val)
	assert.False(t, m.ContainsKey("b"))
	_, ok = m.Remove("b")
	assert.False(t, ok)
	assert.Equal(t, 2, m.Len())
}

func TestRWMap_Delete(t *testing.T) {
	m := NewRWMap[string, int](nil)
	assert.Equal(t, 0, m.Len())
	m.Delete("a")
	m.Insert("a", 1)
	assert.Equal(t, 1, m.Len())
	m.Delete("a")
	assert.Equal(t, 0, m.Len())
}

func TestRWMap_Values(t *testing.T) {
	m := NewRWMap[string, int](nil)
	assert.Equal(t, []int{}, m.Values())
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	values := m.Values()
	slices.Sort(values)
	assert.Equal(t, []int{1, 2, 3}, values)
}

func TestRWMap_Keys(t *testing.T) {
	m := NewRWMapPtr[string, int](nil)
	assert.Equal(t, []string{}, m.Keys())
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	keys := m.Keys()
	slices.Sort(keys)
	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestRWMap_Each(t *testing.T) {
	m := NewRWMap[string, int](nil)
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	arr := make([]string, 0)
	m.Each(func(k string, v int) {
		arr = append(arr, fmt.Sprintf("%s_%d", k, v))
	})
	slices.Sort(arr)
	assert.Equal(t, []string{"a_1", "b_2", "c_3"}, arr)
}

func TestRWMap_InitialValue(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1, "b": 2, "c": 3})
	assert.Equal(t, 3, m.Len())
	assert.Equal(t, 2, first(m.Get("b")))
}

func TestRWMap_Load(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1, "b": 2, "c": 3})
	theMap := m.Load()
	m.Insert("a", 4)
	assert.Equal(t, 4, theMap["a"])
	assert.Equal(t, 4, first(m.Get("a")))
	theMap["a"] = 5
	assert.Equal(t, 5, theMap["a"])
	assert.Equal(t, 5, first(m.Get("a")))
}

func TestRWMap_Clone(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1, "b": 2, "c": 3})
	clonedMap := m.Clone()
	m.Insert("a", 4)
	assert.Equal(t, 1, clonedMap["a"])
	assert.Equal(t, 4, first(m.Get("a")))
}

func TestMap_IsEmpty(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1})
	assert.False(t, m.IsEmpty())
	m.Delete("a")
	assert.True(t, m.IsEmpty())
}

func TestMap_Clear(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1, "b": 2, "c": 3})
	assert.False(t, m.IsEmpty())
	m.Clear()
	assert.True(t, m.IsEmpty())
	assert.Equal(t, map[string]int{}, m.Load())
}

func TestSlice(t *testing.T) {
	m := NewSlicePtr[int](nil)
	assert.Equal(t, 0, m.Len())
	m.Append(1, 2, 3)
	assert.Equal(t, 3, m.Len())
	assert.Equal(t, []int{1, 2, 3}, m.Load())
	val2 := m.Shift()
	assert.Equal(t, 1, val2)
	m.Unshift(4)
	assert.Equal(t, []int{4, 2, 3}, m.Load())
	val2 = m.Pop()
	assert.Equal(t, []int{4, 2}, m.Load())
	assert.Equal(t, 2, m.Remove(1))
	assert.Equal(t, []int{4}, m.Load())
	assert.Panics(t, func() { m.Remove(1) })
	m.Append(5, 6, 7)
	assert.Equal(t, []int{4, 5, 6, 7}, m.Load())
	assert.Equal(t, 6, m.Get(2))
	m.Insert(2, 8)
	assert.Equal(t, []int{4, 5, 8, 6, 7}, m.Load())
}

func TestRWSlice(t *testing.T) {
	m := NewRWSlice[int](nil)
	assert.Equal(t, 0, m.Len())
	m.Append(1, 2, 3)
	assert.Equal(t, 3, m.Len())
	assert.Equal(t, []int{1, 2, 3}, m.Load())
	val2 := m.Shift()
	assert.Equal(t, 1, val2)
	m.Unshift(4)
	assert.Equal(t, []int{4, 2, 3}, m.Load())
	val2 = m.Pop()
	assert.Equal(t, []int{4, 2}, m.Load())
	m.Remove(1)
	assert.Equal(t, []int{4}, m.Load())
	m.Append(5, 6, 7)
	assert.Equal(t, []int{4, 5, 6, 7}, m.Load())
	assert.Equal(t, 6, m.Get(2))
	m.Insert(2, 8)
	assert.Equal(t, []int{4, 5, 8, 6, 7}, m.Load())
}

func TestSlice_InitialValue(t *testing.T) {
	m := NewSlice([]int{1, 2, 3})
	assert.Equal(t, []int{1, 2, 3}, m.Load())
}

func TestRWSlice_Clone(t *testing.T) {
	m := NewRWSlice[int](nil)
	m.Store([]int{1, 2, 3})
	clonedSlice := m.Clone()
	assert.Equal(t, []int{1, 2, 3}, clonedSlice)
}

func TestRWSlice_Each(t *testing.T) {
	m := NewRWSlicePtr[int](nil)
	m.Append(1, 2, 3)
	arr := make([]string, 0)
	m.Each(func(el int) {
		arr = append(arr, fmt.Sprintf("E%d", el))
	})
	assert.Equal(t, []string{"E1", "E2", "E3"}, arr)
}

func TestRWSlice_Filter(t *testing.T) {
	m := NewRWSlicePtr([]int{1, 2, 3, 4, 5, 6})
	out := m.Filter(func(el int) bool { return el%2 == 0 })
	assert.Equal(t, 3, len(out))
	assert.Equal(t, []int{2, 4, 6}, out)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, m.Load())
}

func TestSlice_IsEmpty(t *testing.T) {
	s := NewSlice([]int{1})
	assert.False(t, s.IsEmpty())
	s.Pop()
	assert.True(t, s.IsEmpty())
}

func TestSlice_Clear(t *testing.T) {
	s := NewSlice([]int{1, 2, 3})
	assert.False(t, s.IsEmpty())
	s.Clear()
	assert.True(t, s.IsEmpty())
	assert.Equal(t, []int{}, s.Load())
}

func TestNumber(t *testing.T) {
	n1 := NewNumber(uint64(0))
	assert.Equal(t, uint64(0), n1.Load())
	n1.Add(10)
	assert.Equal(t, uint64(10), n1.Load())
	n1.Sub(5)
	assert.Equal(t, uint64(5), n1.Load())

	n2 := NewNumberPtr(uint64(0))
	assert.Equal(t, uint64(0), n2.Load())
	n2.Add(10)
	assert.Equal(t, uint64(10), n2.Load())
	n2.Sub(5)
	assert.Equal(t, uint64(5), n2.Load())

	n3 := NewRWNumberPtr(uint64(0))
	assert.Equal(t, uint64(0), n3.Load())
	n3.Add(10)
	assert.Equal(t, uint64(10), n3.Load())
	n3.Sub(5)
	assert.Equal(t, uint64(5), n3.Load())

	n4 := NewRWNumber(uint64(0))
	assert.Equal(t, uint64(0), n4.Load())
	n4.Add(10)
	assert.Equal(t, uint64(10), n4.Load())
	n4.Sub(5)
	assert.Equal(t, uint64(5), n4.Load())
}

func TestValueUsage(t *testing.T) {
	type MyStruct struct {
		Value RWMutex[string]
	}
	m := MyStruct{}
	m.Value.Store("hello world")
	assert.Equal(t, "hello world", m.Value.Load())
}

func TestBaseMutex_LockUnlock(t *testing.T) {
	m := &baseMutex[int]{v: 42}
	m.Lock()
	*m.GetPointer() = 100
	m.Unlock()
	assert.Equal(t, 100, m.Load())
}

func TestBaseMutex_With(t *testing.T) {
	m := &baseMutex[string]{v: "old"}
	m.With(func(v *string) {
		*v = "new"
	})
	assert.Equal(t, "new", m.Load())
}

func TestBaseMutex_RWith(t *testing.T) {
	m := &baseMutex[string]{v: "old"}
	m.RWith(func(v string) {
		assert.Equal(t, "old", v)
	})
}

func TestBaseMutex_Store(t *testing.T) {
	m := &baseMutex[int]{v: 42}
	m.Store(100)
	assert.Equal(t, 100, m.Load())
}

func TestBaseMutex_Swap(t *testing.T) {
	m := &baseMutex[string]{v: "old"}
	old := m.Swap("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Load())
}

func TestBaseMutex_GetPointer(t *testing.T) {
	m := &baseMutex[int]{v: 42}
	ptr := m.GetPointer()
	*ptr = 100
	assert.Equal(t, 100, m.Load())
}

func TestBaseMutex_RLockRUnlock(t *testing.T) {
	m := &baseMutex[string]{v: "old"}
	m.RLock()
	assert.Equal(t, "old", *m.GetPointer())
	m.RUnlock()
}

func TestBaseRWMutex_LockUnlock(t *testing.T) {
	m := &baseRWMutex[int]{v: 42}
	m.Lock()
	*m.GetPointer() = 100
	m.Unlock()
	assert.Equal(t, 100, m.Load())
}

func TestBaseRWMutex_RLockRUnlock(t *testing.T) {
	m := &baseRWMutex[string]{v: "old"}
	m.RLock()
	assert.Equal(t, "old", *m.GetPointer())
	m.RUnlock()
}

func TestBaseRWMutex_With(t *testing.T) {
	m := &baseRWMutex[string]{v: "old"}
	m.With(func(v *string) {
		*v = "new"
	})
	assert.Equal(t, "new", m.Load())
}

func TestBaseRWMutex_RWith(t *testing.T) {
	m := &baseRWMutex[string]{v: "old"}
	m.RWith(func(v string) {
		assert.Equal(t, "old", v)
	})
}

func TestBaseRWMutex_Store(t *testing.T) {
	m := &baseRWMutex[int]{v: 42}
	m.Store(100)
	assert.Equal(t, 100, m.Load())
}

func TestBaseRWMutex_Swap(t *testing.T) {
	m := &baseRWMutex[string]{v: "old"}
	old := m.Swap("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Load())
}

func TestBaseRWMutex_GetPointer(t *testing.T) {
	m := &baseRWMutex[int]{v: 42}
	ptr := m.GetPointer()
	*ptr = 100
	assert.Equal(t, 100, m.Load())
}

func TestSliceMutex_Append(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 2}}}
	s.Append(3, 4)
	assert.Equal(t, []int{1, 2, 3, 4}, s.Load())
}

func TestSliceMutex_Unshift(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 2}}}
	s.Unshift(0)
	assert.Equal(t, []int{0, 1, 2}, s.Load())
}

func TestSliceMutex_Shift(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 2}}}
	val := s.Shift()
	assert.Equal(t, 1, val)
	assert.Equal(t, []int{2}, s.Load())
}

func TestSliceMutex_Pop(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 2}}}
	val := s.Pop()
	assert.Equal(t, 2, val)
	assert.Equal(t, []int{1}, s.Load())
}

func TestSliceMutex_Clone(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 2}}}
	clone := s.Clone()
	assert.Equal(t, []int{1, 2}, clone)
}

func TestSliceMutex_Len(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 2, 3}}}
	assert.Equal(t, 3, s.Len())
}

func TestSliceMutex_IsEmpty(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{}}}
	assert.True(t, s.IsEmpty())
	s.Append(1)
	assert.False(t, s.IsEmpty())
}

func TestSliceMutex_Get(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 2, 3}}}
	assert.Equal(t, 2, s.Get(1))
}

func TestSliceMutex_Remove(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 2, 3}}}
	val := s.Remove(1)
	assert.Equal(t, 2, val)
	assert.Equal(t, []int{1, 3}, s.Load())
}

func TestSliceMutex_Insert(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 3}}}
	s.Insert(1, 2)
	assert.Equal(t, []int{1, 2, 3}, s.Load())
}

func TestSliceMutex_Filter(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 2, 3, 4}}}
	filtered := s.Filter(func(v int) bool { return v%2 == 0 })
	assert.Equal(t, []int{2, 4}, filtered)
	assert.Equal(t, []int{1, 2, 3, 4}, s.Load())
}

func TestMapMutex_Insert(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{}}}
	m.Insert("a", 1)
	assert.Equal(t, 1, m.Load()["a"])
}

func TestMapMutex_Get(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{"a": 1}}}
	val, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, val)
}

func TestMapMutex_Remove(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{"a": 1}}}
	val, ok := m.Remove("a")
	assert.True(t, ok)
	assert.Equal(t, 1, val)
	assert.False(t, m.ContainsKey("a"))
}

func TestMapMutex_Keys(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{"a": 1, "b": 2}}}
	keys := m.Keys()
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, "a")
	assert.Contains(t, keys, "b")
}

func TestNumberMutex_Add(t *testing.T) {
	n := &NumberMutex[int]{baseMutex[int]{v: 10}}
	n.Add(5)
	assert.Equal(t, 15, n.Load())
}

func TestNumberMutex_Sub(t *testing.T) {
	n := &NumberMutex[int]{baseMutex[int]{v: 10}}
	n.Sub(5)
	assert.Equal(t, 5, n.Load())
}

func TestSliceMutex_Each(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 2, 3}}}
	var sum int
	s.Each(func(v int) {
		sum += v
	})
	assert.Equal(t, 6, sum)
}

func TestSliceMutex_Clear(t *testing.T) {
	s := &SliceMutex[int]{baseMutex[[]int]{v: []int{1, 2, 3}}}
	s.Clear()
	assert.Equal(t, []int{}, s.Load())
	assert.True(t, s.IsEmpty())
}

func TestMapMutex_Clear(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{"a": 1}}}
	m.Clear()
	assert.Equal(t, map[string]int{}, m.Load())
	assert.True(t, m.IsEmpty())
}

func TestMapMutex_GetKeyValue(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{"a": 1}}}
	k, v, ok := m.GetKeyValue("a")
	assert.True(t, ok)
	assert.Equal(t, "a", k)
	assert.Equal(t, 1, v)

	_, _, ok = m.GetKeyValue("b")
	assert.False(t, ok)
}

func TestMapMutex_Delete(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{"a": 1}}}
	m.Delete("a")
	assert.False(t, m.ContainsKey("a"))
}

func TestMapMutex_Len(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{"a": 1, "b": 2}}}
	assert.Equal(t, 2, m.Len())
}

func TestMapMutex_IsEmpty(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{}}}
	assert.True(t, m.IsEmpty())
	m.Insert("a", 1)
	assert.False(t, m.IsEmpty())
}

func TestMapMutex_Each(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{"a": 1, "b": 2}}}
	var sum int
	m.Each(func(k string, v int) {
		sum += v
	})
	assert.Equal(t, 3, sum)
}

func TestMapMutex_Values(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{"a": 1, "b": 2}}}
	values := m.Values()
	assert.Len(t, values, 2)
	assert.Contains(t, values, 1)
	assert.Contains(t, values, 2)
}

func TestMapMutex_Clone(t *testing.T) {
	m := &MapMutex[string, int]{baseMutex[map[string]int]{v: map[string]int{"a": 1}}}
	clone := m.Clone()
	assert.Equal(t, 1, clone["a"])
	m.Insert("a", 2) // Original should be unaffected
	assert.Equal(t, 1, clone["a"])
}

func TestMapRWMutex_Clear(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{"a": 1}}}
	m.Clear()
	assert.Equal(t, map[string]int{}, m.Load())
	assert.True(t, m.IsEmpty())
}

func TestMapRWMutex_Insert(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{}}}
	m.Insert("a", 1)
	assert.Equal(t, 1, m.Load()["a"])
}

func TestMapRWMutex_Get(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{"a": 1}}}
	val, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, val)

	_, ok = m.Get("b")
	assert.False(t, ok)
}

func TestMapRWMutex_GetKeyValue(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{"a": 1}}}
	k, v, ok := m.GetKeyValue("a")
	assert.True(t, ok)
	assert.Equal(t, "a", k)
	assert.Equal(t, 1, v)

	_, _, ok = m.GetKeyValue("b")
	assert.False(t, ok)
}

func TestMapRWMutex_ContainsKey(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{"a": 1}}}
	assert.True(t, m.ContainsKey("a"))
	assert.False(t, m.ContainsKey("b"))
}

func TestMapRWMutex_Remove(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{"a": 1}}}
	val, ok := m.Remove("a")
	assert.True(t, ok)
	assert.Equal(t, 1, val)
	assert.False(t, m.ContainsKey("a"))

	_, ok = m.Remove("a")
	assert.False(t, ok)
}

func TestMapRWMutex_Delete(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{"a": 1}}}
	m.Delete("a")
	assert.False(t, m.ContainsKey("a"))
}

func TestMapRWMutex_Len(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{"a": 1, "b": 2}}}
	assert.Equal(t, 2, m.Len())
}

func TestMapRWMutex_IsEmpty(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{}}}
	assert.True(t, m.IsEmpty())
	m.Insert("a", 1)
	assert.False(t, m.IsEmpty())
}

func TestMapRWMutex_Each(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{"a": 1, "b": 2}}}
	var sum int
	m.Each(func(k string, v int) {
		sum += v
	})
	assert.Equal(t, 3, sum)
}

func TestMapRWMutex_Keys(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{"a": 1, "b": 2}}}
	keys := m.Keys()
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, "a")
	assert.Contains(t, keys, "b")
}

func TestMapRWMutex_Values(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{"a": 1, "b": 2}}}
	values := m.Values()
	assert.Len(t, values, 2)
	assert.Contains(t, values, 1)
	assert.Contains(t, values, 2)
}

func TestMapRWMutex_Clone(t *testing.T) {
	m := &MapRWMutex[string, int]{baseRWMutex[map[string]int]{v: map[string]int{"a": 1}}}
	clone := m.Clone()
	assert.Equal(t, 1, clone["a"])
	m.Insert("a", 2) // Original should be unaffected
	assert.Equal(t, 1, clone["a"])
}

func TestSliceRWMutex_Each(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 2, 3}}}
	var sum int
	s.Each(func(v int) {
		sum += v
	})
	assert.Equal(t, 6, sum)
}

func TestSliceRWMutex_Clear(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 2, 3}}}
	s.Clear()
	assert.Equal(t, []int{}, s.Load())
	assert.True(t, s.IsEmpty())
}

func TestSliceRWMutex_Append(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 2}}}
	s.Append(3, 4)
	assert.Equal(t, []int{1, 2, 3, 4}, s.Load())
}

func TestSliceRWMutex_Unshift(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 2}}}
	s.Unshift(0)
	assert.Equal(t, []int{0, 1, 2}, s.Load())
}

func TestSliceRWMutex_Shift(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 2}}}
	val := s.Shift()
	assert.Equal(t, 1, val)
	assert.Equal(t, []int{2}, s.Load())
}

func TestSliceRWMutex_Pop(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 2}}}
	val := s.Pop()
	assert.Equal(t, 2, val)
	assert.Equal(t, []int{1}, s.Load())
}

func TestSliceRWMutex_Clone(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 2}}}
	clone := s.Clone()
	assert.Equal(t, []int{1, 2}, clone)
	s.Append(3) // Original should be unaffected
	assert.Equal(t, []int{1, 2}, clone)
}

func TestSliceRWMutex_Len(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 2, 3}}}
	assert.Equal(t, 3, s.Len())
}

func TestSliceRWMutex_IsEmpty(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{}}}
	assert.True(t, s.IsEmpty())
	s.Append(1)
	assert.False(t, s.IsEmpty())
}

func TestSliceRWMutex_Get(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 2, 3}}}
	assert.Equal(t, 2, s.Get(1))
}

func TestSliceRWMutex_Remove(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 2, 3}}}
	val := s.Remove(1)
	assert.Equal(t, 2, val)
	assert.Equal(t, []int{1, 3}, s.Load())
}

func TestSliceRWMutex_Insert(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 3}}}
	s.Insert(1, 2)
	assert.Equal(t, []int{1, 2, 3}, s.Load())
}

func TestSliceRWMutex_Filter(t *testing.T) {
	s := &SliceRWMutex[int]{baseRWMutex[[]int]{v: []int{1, 2, 3, 4}}}
	filtered := s.Filter(func(v int) bool { return v%2 == 0 })
	assert.Equal(t, []int{2, 4}, filtered)
	assert.Equal(t, []int{1, 2, 3, 4}, s.Load()) // Original unchanged
}

func TestNumberRWMutex_Add(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		n := &NumberRWMutex[int]{baseRWMutex[int]{v: 10}}
		n.Add(5)
		assert.Equal(t, 15, n.Load())
	})

	t.Run("float64", func(t *testing.T) {
		n := &NumberRWMutex[float64]{baseRWMutex[float64]{v: 10.5}}
		n.Add(2.5)
		assert.Equal(t, 13.0, n.Load())
	})

	t.Run("uint", func(t *testing.T) {
		n := &NumberRWMutex[uint]{baseRWMutex[uint]{v: 10}}
		n.Add(5)
		assert.Equal(t, uint(15), n.Load())
	})
}

func TestNumberRWMutex_Sub(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		n := &NumberRWMutex[int]{baseRWMutex[int]{v: 10}}
		n.Sub(3)
		assert.Equal(t, 7, n.Load())
	})

	t.Run("float64", func(t *testing.T) {
		n := &NumberRWMutex[float64]{baseRWMutex[float64]{v: 10.5}}
		n.Sub(2.5)
		assert.Equal(t, 8.0, n.Load())
	})

	t.Run("uint", func(t *testing.T) {
		n := &NumberRWMutex[uint]{baseRWMutex[uint]{v: 10}}
		n.Sub(3)
		assert.Equal(t, uint(7), n.Load())
	})
}

func TestNumberRWMutex_ConcurrentOperations(t *testing.T) {
	n := &NumberRWMutex[int]{baseRWMutex[int]{v: 0}}
	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(2)

	// Concurrent adder
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			n.Add(1)
		}
	}()

	// Concurrent subtractor
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			n.Sub(1)
		}
	}()

	wg.Wait()
	assert.Equal(t, 0, n.Load())
}
