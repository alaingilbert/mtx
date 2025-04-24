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
	if got := m.Load(); got != "new" {
		t.Errorf("expected %q, got %q", "new", got)
	}
}

func TestMtx_With(t *testing.T) {
	m := NewMtx("old")
	m.With(func(v *string) {
		*v = "new"
	})
	if got := m.Load(); got != "new" {
		t.Errorf("expected %q, got %q", "new", got)
	}
}

func TestMtx_RWith(t *testing.T) {
	m := NewMtx("old")
	m.RWith(func(v string) {
		if v != "old" {
			t.Errorf("expected %q, got %q", "old", v)
		}
	})
}

func TestMtx_Store(t *testing.T) {
	m := NewMtx("old")
	if got := m.Load(); got != "old" {
		t.Errorf("expected %q, got %q", "old", got)
	}
	m.Store("new")
	if got := m.Load(); got != "new" {
		t.Errorf("expected %q, got %q", "new", got)
	}
}

func TestMtx_Swap(t *testing.T) {
	m := NewMtx("old")
	old := m.Swap("new")
	if old != "old" {
		t.Errorf("expected %q, got %q", "old", old)
	}
	if got := m.Load(); got != "new" {
		t.Errorf("expected %q, got %q", "new", got)
	}
}

func TestMtx_GetPointer(t *testing.T) {
	someString := "old"
	orig := &someString
	m := NewMtx(orig)
	val := m.GetPointer()
	**val = "new"
	if someString != "new" {
		t.Errorf("expected %q, got %q", "new", someString)
	}
	if **val != "new" {
		t.Errorf("expected %q, got %q", "new", **val)
	}
	if *orig != "new" {
		t.Errorf("expected %q, got %q", "new", *orig)
	}
}

func TestMtxPtr_Swap(t *testing.T) {
	m := NewMtxPtr("old")
	old := m.Swap("new")
	if old != "old" {
		t.Errorf("expected %q, got %q", "old", old)
	}
	if got := m.Load(); got != "new" {
		t.Errorf("expected %q, got %q", "new", got)
	}
}

func TestMtx_RLockRUnlock(t *testing.T) {
	m := NewMtx("old")
	m.RLock()
	val := m.GetPointer()
	if *val != "old" {
		t.Errorf("expected %q, got %q", "old", *val)
	}
	m.RUnlock()
}

func TestRWMtx_RLockRUnlock(t *testing.T) {
	m := NewRWMtx("old")
	m.RLock()
	val := m.GetPointer()
	if *val != "old" {
		t.Errorf("expected %q, got %q", "old", *val)
	}
	m.RUnlock()
}

func TestRWMtx_Swap(t *testing.T) {
	m := NewRWMtx("old")
	old := m.Swap("new")
	if old != "old" {
		t.Errorf("expected %q, got %q", "old", old)
	}
	if got := m.Load(); got != "new" {
		t.Errorf("expected %q, got %q", "new", got)
	}
}

func TestRWMtxPtr_Swap(t *testing.T) {
	m := NewRWMtxPtr("old")
	old := m.Swap("new")
	if old != "old" {
		t.Errorf("expected %q, got %q", "old", old)
	}
	if got := m.Load(); got != "new" {
		t.Errorf("expected %q, got %q", "new", got)
	}
}

func TestRWMtx_Val(t *testing.T) {
	someString := "old"
	orig := &someString
	m := NewRWMtx(orig)
	val := m.GetPointer()
	**val = "new"
	if **val != "new" {
		t.Errorf("expected %q, got %q", "new", **val)
	}
	if *orig != "new" {
		t.Errorf("expected %q, got %q", "new", *orig)
	}
}

func TestMap_Get(t *testing.T) {
	m := NewMap[string, int](nil)
	_, ok := m.Get("a")
	if ok {
		t.Error("expected false, got true")
	}
	m.Insert("a", 1)
	el, ok := m.Get("a")
	if !ok {
		t.Error("expected true, got false")
	}
	if el != 1 {
		t.Errorf("expected 1, got %d", el)
	}
}

func TestMap_GetKeyValue(t *testing.T) {
	m := NewMap[string, int](nil)
	_, _, ok := m.GetKeyValue("a")
	if ok {
		t.Error("expected false, got true")
	}
	m.Insert("a", 1)
	key, value, ok := m.GetKeyValue("a")
	if !ok {
		t.Error("expected true, got false")
	}
	if key != "a" {
		t.Errorf("expected %q, got %q", "a", key)
	}
	if value != 1 {
		t.Errorf("expected 1, got %d", value)
	}
}

func TestMap_HasKey(t *testing.T) {
	m := NewMap[string, int](nil)
	if m.ContainsKey("a") {
		t.Error("expected false, got true")
	}
	m.Insert("a", 1)
	if !m.ContainsKey("a") {
		t.Error("expected true, got false")
	}
	m.Delete("a")
	if m.ContainsKey("a") {
		t.Error("expected false, got true")
	}
}

func TestMap_Remove(t *testing.T) {
	m := NewMap[string, int](nil)
	m.Insert("a", 1)
	m.Insert("b", 2)
	m.Insert("c", 3)
	if m.Len() != 3 {
		t.Errorf("expected 3, got %d", m.Len())
	}
	if !m.ContainsKey("b") {
		t.Error("expected true, got false")
	}
	val, ok := m.Remove("b")
	if !ok {
		t.Error("expected true, got false")
	}
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
	if m.ContainsKey("b") {
		t.Error("expected false, got true")
	}
	_, ok = m.Remove("b")
	if ok {
		t.Error("expected false, got true")
	}
	if m.Len() != 2 {
		t.Errorf("expected 2, got %d", m.Len())
	}
}

func TestMap_Delete(t *testing.T) {
	m := NewMap[string, int](nil)
	if m.Len() != 0 {
		t.Errorf("expected 0, got %d", m.Len())
	}
	m.Delete("a")
	m.Insert("a", 1)
	if m.Len() != 1 {
		t.Errorf("expected 1, got %d", m.Len())
	}
	m.Delete("a")
	if m.Len() != 0 {
		t.Errorf("expected 0, got %d", m.Len())
	}
}

func TestMap_Values(t *testing.T) {
	m := NewMap[string, int](nil)
	if len(m.Values()) != 0 {
		t.Errorf("expected empty slice, got %v", m.Values())
	}
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	values := m.Values()
	slices.Sort(values)
	expected := []int{1, 2, 3}
	if !slices.Equal(values, expected) {
		t.Errorf("expected %v, got %v", expected, values)
	}
}

func TestMap_Keys(t *testing.T) {
	m := NewMapPtr[string, int](nil)
	if len(m.Keys()) != 0 {
		t.Errorf("expected empty slice, got %v", m.Keys())
	}
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	keys := m.Keys()
	slices.Sort(keys)
	expected := []string{"a", "b", "c"}
	if !slices.Equal(keys, expected) {
		t.Errorf("expected %v, got %v", expected, keys)
	}
}

func TestMap_Each(t *testing.T) {
	m := NewMap[string, int](nil)
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	arr := make([]string, 0)
	m.Each(func(k string, v int) {
		arr = append(arr, fmt.Sprintf("%s_%d", k, v))
	})
	slices.Sort(arr)
	expected := []string{"a_1", "b_2", "c_3"}
	if !slices.Equal(arr, expected) {
		t.Errorf("expected %v, got %v", expected, arr)
	}
}

func TestMap_Clone(t *testing.T) {
	m := NewMap[string, int](nil)
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	clonedMap := m.Clone()
	if clonedMap["a"] != 1 {
		t.Errorf("expected 1, got %d", clonedMap["a"])
	}
}

func TestRWMap_Get(t *testing.T) {
	m := NewRWMap[string, int](nil)
	_, ok := m.Get("a")
	if ok {
		t.Error("expected false, got true")
	}
	m.Insert("a", 1)
	el, ok := m.Get("a")
	if !ok {
		t.Error("expected true, got false")
	}
	if el != 1 {
		t.Errorf("expected 1, got %d", el)
	}
}

func TestRWMap_HasKey(t *testing.T) {
	m := NewRWMap[string, int](nil)
	if m.ContainsKey("a") {
		t.Error("expected false, got true")
	}
	m.Insert("a", 1)
	if !m.ContainsKey("a") {
		t.Error("expected true, got false")
	}
	m.Delete("a")
	if m.ContainsKey("a") {
		t.Error("expected false, got true")
	}
}

func TestRWMap_Remove(t *testing.T) {
	m := NewRWMap[string, int](nil)
	m.Insert("a", 1)
	m.Insert("b", 2)
	m.Insert("c", 3)
	if m.Len() != 3 {
		t.Errorf("expected 3, got %d", m.Len())
	}
	if !m.ContainsKey("b") {
		t.Error("expected true, got false")
	}
	val, ok := m.Remove("b")
	if !ok {
		t.Error("expected true, got false")
	}
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
	if m.ContainsKey("b") {
		t.Error("expected false, got true")
	}
	_, ok = m.Remove("b")
	if ok {
		t.Error("expected false, got true")
	}
	if m.Len() != 2 {
		t.Errorf("expected 2, got %d", m.Len())
	}
}

func TestRWMap_Delete(t *testing.T) {
	m := NewRWMap[string, int](nil)
	if m.Len() != 0 {
		t.Errorf("expected 0, got %d", m.Len())
	}
	m.Delete("a")
	m.Insert("a", 1)
	if m.Len() != 1 {
		t.Errorf("expected 1, got %d", m.Len())
	}
	m.Delete("a")
	if m.Len() != 0 {
		t.Errorf("expected 0, got %d", m.Len())
	}
}

func TestRWMap_Values(t *testing.T) {
	m := NewRWMap[string, int](nil)
	if len(m.Values()) != 0 {
		t.Errorf("expected empty slice, got %v", m.Values())
	}
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	values := m.Values()
	slices.Sort(values)
	expected := []int{1, 2, 3}
	if !slices.Equal(values, expected) {
		t.Errorf("expected %v, got %v", expected, values)
	}
}

func TestRWMap_Keys(t *testing.T) {
	m := NewRWMapPtr[string, int](nil)
	if len(m.Keys()) != 0 {
		t.Errorf("expected empty slice, got %v", m.Keys())
	}
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	keys := m.Keys()
	slices.Sort(keys)
	expected := []string{"a", "b", "c"}
	if !slices.Equal(keys, expected) {
		t.Errorf("expected %v, got %v", expected, keys)
	}
}

func TestRWMap_Each(t *testing.T) {
	m := NewRWMap[string, int](nil)
	m.Store(map[string]int{"a": 1, "b": 2, "c": 3})
	arr := make([]string, 0)
	m.Each(func(k string, v int) {
		arr = append(arr, fmt.Sprintf("%s_%d", k, v))
	})
	slices.Sort(arr)
	expected := []string{"a_1", "b_2", "c_3"}
	if !slices.Equal(arr, expected) {
		t.Errorf("expected %v, got %v", expected, arr)
	}
}

func TestRWMap_InitialValue(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1, "b": 2, "c": 3})
	if m.Len() != 3 {
		t.Errorf("expected 3, got %d", m.Len())
	}
	if first(m.Get("b")) != 2 {
		t.Errorf("expected 2, got %d", first(m.Get("b")))
	}
}

func TestRWMap_Load(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1, "b": 2, "c": 3})
	theMap := m.Load()
	m.Insert("a", 4)
	if theMap["a"] != 4 {
		t.Errorf("expected 4, got %d", theMap["a"])
	}
	if first(m.Get("a")) != 4 {
		t.Errorf("expected 4, got %d", first(m.Get("a")))
	}
	theMap["a"] = 5
	if theMap["a"] != 5 {
		t.Errorf("expected 5, got %d", theMap["a"])
	}
	if first(m.Get("a")) != 5 {
		t.Errorf("expected 5, got %d", first(m.Get("a")))
	}
}

func TestRWMap_Clone(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1, "b": 2, "c": 3})
	clonedMap := m.Clone()
	m.Insert("a", 4)
	if clonedMap["a"] != 1 {
		t.Errorf("expected 1, got %d", clonedMap["a"])
	}
	if first(m.Get("a")) != 4 {
		t.Errorf("expected 4, got %d", first(m.Get("a")))
	}
}

func TestMap_IsEmpty(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1})
	if m.IsEmpty() {
		t.Error("expected false, got true")
	}
	m.Delete("a")
	if !m.IsEmpty() {
		t.Error("expected true, got false")
	}
}

func TestMap_Clear(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1, "b": 2, "c": 3})
	if m.IsEmpty() {
		t.Error("expected false, got true")
	}
	m.Clear()
	if !m.IsEmpty() {
		t.Error("expected true, got false")
	}
	if len(m.Load()) != 0 {
		t.Errorf("expected empty map, got %v", m.Load())
	}
}

func TestSlice(t *testing.T) {
	m := NewSlicePtr[int](nil)
	if m.Len() != 0 {
		t.Errorf("expected 0, got %d", m.Len())
	}
	m.Append(1, 2, 3)
	if m.Len() != 3 {
		t.Errorf("expected 3, got %d", m.Len())
	}
	if !slices.Equal(m.Load(), []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", m.Load())
	}
	val2 := m.Shift()
	if val2 != 1 {
		t.Errorf("expected 1, got %d", val2)
	}
	m.Unshift(4)
	if !slices.Equal(m.Load(), []int{4, 2, 3}) {
		t.Errorf("expected [4 2 3], got %v", m.Load())
	}
	val2 = m.Pop()
	if !slices.Equal(m.Load(), []int{4, 2}) {
		t.Errorf("expected [4 2], got %v", m.Load())
	}
	if val2 != 3 {
		t.Errorf("expected 3, got %d", val2)
	}
	val2 = m.Remove(1)
	if val2 != 2 {
		t.Errorf("expected 2, got %d", val2)
	}
	if !slices.Equal(m.Load(), []int{4}) {
		t.Errorf("expected [4], got %v", m.Load())
	}
	// Test panic
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic, got none")
			}
		}()
		m.Remove(1)
	}()
	m.Append(5, 6, 7)
	if !slices.Equal(m.Load(), []int{4, 5, 6, 7}) {
		t.Errorf("expected [4 5 6 7], got %v", m.Load())
	}
	if m.Get(2) != 6 {
		t.Errorf("expected 6, got %d", m.Get(2))
	}
	m.Insert(2, 8)
	if !slices.Equal(m.Load(), []int{4, 5, 8, 6, 7}) {
		t.Errorf("expected [4 5 8 6 7], got %v", m.Load())
	}
}

func TestRWSlice(t *testing.T) {
	m := NewRWSlice[int](nil)
	if m.Len() != 0 {
		t.Errorf("expected 0, got %d", m.Len())
	}
	m.Append(1, 2, 3)
	if m.Len() != 3 {
		t.Errorf("expected 3, got %d", m.Len())
	}
	if !slices.Equal(m.Load(), []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", m.Load())
	}
	val2 := m.Shift()
	if val2 != 1 {
		t.Errorf("expected 1, got %d", val2)
	}
	m.Unshift(4)
	if !slices.Equal(m.Load(), []int{4, 2, 3}) {
		t.Errorf("expected [4 2 3], got %v", m.Load())
	}
	val2 = m.Pop()
	if !slices.Equal(m.Load(), []int{4, 2}) {
		t.Errorf("expected [4 2], got %v", m.Load())
	}
	m.Remove(1)
	if !slices.Equal(m.Load(), []int{4}) {
		t.Errorf("expected [4], got %v", m.Load())
	}
	m.Append(5, 6, 7)
	if !slices.Equal(m.Load(), []int{4, 5, 6, 7}) {
		t.Errorf("expected [4 5 6 7], got %v", m.Load())
	}
	if m.Get(2) != 6 {
		t.Errorf("expected 6, got %d", m.Get(2))
	}
	m.Insert(2, 8)
	if !slices.Equal(m.Load(), []int{4, 5, 8, 6, 7}) {
		t.Errorf("expected [4 5 8 6 7], got %v", m.Load())
	}
}

func TestSlice_InitialValue(t *testing.T) {
	m := NewSlice([]int{1, 2, 3})
	if !slices.Equal(m.Load(), []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", m.Load())
	}
}

func TestRWSlice_Clone(t *testing.T) {
	m := NewRWSlice[int](nil)
	m.Store([]int{1, 2, 3})
	clonedSlice := m.Clone()
	if !slices.Equal(clonedSlice, []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", clonedSlice)
	}
}

func TestRWSlice_Each(t *testing.T) {
	m := NewRWSlicePtr[int](nil)
	m.Append(1, 2, 3)
	arr := make([]string, 0)
	m.Each(func(el int) {
		arr = append(arr, fmt.Sprintf("E%d", el))
	})
	expected := []string{"E1", "E2", "E3"}
	if !slices.Equal(arr, expected) {
		t.Errorf("expected %v, got %v", expected, arr)
	}
}

func TestRWSlice_Filter(t *testing.T) {
	m := NewRWSlicePtr([]int{1, 2, 3, 4, 5, 6})
	out := m.Filter(func(el int) bool { return el%2 == 0 })
	if len(out) != 3 {
		t.Errorf("expected 3, got %d", len(out))
	}
	if !slices.Equal(out, []int{2, 4, 6}) {
		t.Errorf("expected [2 4 6], got %v", out)
	}
	if !slices.Equal(m.Load(), []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("expected [1 2 3 4 5 6], got %v", m.Load())
	}
}

func TestSlice_IsEmpty(t *testing.T) {
	s := NewSlice([]int{1})
	if s.IsEmpty() {
		t.Error("expected false, got true")
	}
	s.Pop()
	if !s.IsEmpty() {
		t.Error("expected true, got false")
	}
}

func TestSlice_Clear(t *testing.T) {
	s := NewSlice([]int{1, 2, 3})
	if s.IsEmpty() {
		t.Error("expected false, got true")
	}
	s.Clear()
	if !s.IsEmpty() {
		t.Error("expected true, got false")
	}
	if !slices.Equal(s.Load(), []int{}) {
		t.Errorf("expected empty slice, got %v", s.Load())
	}
}

func TestNumber(t *testing.T) {
	n1 := NewNumber(uint64(0))
	if n1.Load() != 0 {
		t.Errorf("expected 0, got %d", n1.Load())
	}
	n1.Add(10)
	if n1.Load() != 10 {
		t.Errorf("expected 10, got %d", n1.Load())
	}
	n1.Sub(5)
	if n1.Load() != 5 {
		t.Errorf("expected 5, got %d", n1.Load())
	}

	n2 := NewNumberPtr(uint64(0))
	if n2.Load() != 0 {
		t.Errorf("expected 0, got %d", n2.Load())
	}
	n2.Add(10)
	if n2.Load() != 10 {
		t.Errorf("expected 10, got %d", n2.Load())
	}
	n2.Sub(5)
	if n2.Load() != 5 {
		t.Errorf("expected 5, got %d", n2.Load())
	}

	n3 := NewRWNumberPtr(uint64(0))
	if n3.Load() != 0 {
		t.Errorf("expected 0, got %d", n3.Load())
	}
	n3.Add(10)
	if n3.Load() != 10 {
		t.Errorf("expected 10, got %d", n3.Load())
	}
	n3.Sub(5)
	if n3.Load() != 5 {
		t.Errorf("expected 5, got %d", n3.Load())
	}

	n4 := NewRWNumber(uint64(0))
	if n4.Load() != 0 {
		t.Errorf("expected 0, got %d", n4.Load())
	}
	n4.Add(10)
	if n4.Load() != 10 {
		t.Errorf("expected 10, got %d", n4.Load())
	}
	n4.Sub(5)
	if n4.Load() != 5 {
		t.Errorf("expected 5, got %d", n4.Load())
	}
}

func TestValueUsage(t *testing.T) {
	type MyStruct struct {
		Value RWMutex[string]
	}
	m := MyStruct{}
	m.Value.Store("hello world")
	if m.Value.Load() != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", m.Value.Load())
	}
}

func TestBaseMutex_LockUnlock(t *testing.T) {
	m := NewMutex(42)
	m.Lock()
	*m.GetPointer() = 100
	m.Unlock()
	if m.Load() != 100 {
		t.Errorf("expected 100, got %d", m.Load())
	}
}

func TestBaseMutex_With(t *testing.T) {
	m := NewMutex("old")
	m.With(func(v *string) {
		*v = "new"
	})
	if m.Load() != "new" {
		t.Errorf("expected %q, got %q", "new", m.Load())
	}
}

func TestBaseMutex_RWith(t *testing.T) {
	m := NewMutex("old")
	m.RWith(func(v string) {
		if v != "old" {
			t.Errorf("expected %q, got %q", "old", v)
		}
	})
}

func TestBaseMutex_Store(t *testing.T) {
	m := NewMutex(42)
	m.Store(100)
	if m.Load() != 100 {
		t.Errorf("expected 100, got %d", m.Load())
	}
}

func TestBaseMutex_Swap(t *testing.T) {
	m := NewMutex("old")
	old := m.Swap("new")
	if old != "old" {
		t.Errorf("expected %q, got %q", "old", old)
	}
	if m.Load() != "new" {
		t.Errorf("expected %q, got %q", "new", m.Load())
	}
}

func TestBaseMutex_GetPointer(t *testing.T) {
	m := NewMutex(42)
	ptr := m.GetPointer()
	*ptr = 100
	if m.Load() != 100 {
		t.Errorf("expected 100, got %d", m.Load())
	}
}

func TestBaseMutex_RLockRUnlock(t *testing.T) {
	m := NewMutex("old")
	m.RLock()
	if *m.GetPointer() != "old" {
		t.Errorf("expected %q, got %q", "old", *m.GetPointer())
	}
	m.RUnlock()
}

func TestBaseRWMutex_LockUnlock(t *testing.T) {
	m := NewRWMutex(42)
	m.Lock()
	*m.GetPointer() = 100
	m.Unlock()
	if m.Load() != 100 {
		t.Errorf("expected 100, got %d", m.Load())
	}
}

func TestBaseRWMutex_RLockRUnlock(t *testing.T) {
	m := NewRWMutex("old")
	m.RLock()
	if *m.GetPointer() != "old" {
		t.Errorf("expected %q, got %q", "old", *m.GetPointer())
	}
	m.RUnlock()
}

func TestBaseRWMutex_With(t *testing.T) {
	m := NewRWMutex("old")
	m.With(func(v *string) {
		*v = "new"
	})
	if m.Load() != "new" {
		t.Errorf("expected %q, got %q", "new", m.Load())
	}
}

func TestBaseRWMutex_RWith(t *testing.T) {
	m := NewRWMutex("old")
	m.RWith(func(v string) {
		if v != "old" {
			t.Errorf("expected %q, got %q", "old", v)
		}
	})
}

func TestBaseRWMutex_Store(t *testing.T) {
	m := NewRWMutex(42)
	m.Store(100)
	if m.Load() != 100 {
		t.Errorf("expected 100, got %d", m.Load())
	}
}

func TestBaseRWMutex_Swap(t *testing.T) {
	m := NewRWMutex("old")
	old := m.Swap("new")
	if old != "old" {
		t.Errorf("expected %q, got %q", "old", old)
	}
	if m.Load() != "new" {
		t.Errorf("expected %q, got %q", "new", m.Load())
	}
}

func TestBaseRWMutex_GetPointer(t *testing.T) {
	m := NewRWMutex(42)
	ptr := m.GetPointer()
	*ptr = 100
	if m.Load() != 100 {
		t.Errorf("expected 100, got %d", m.Load())
	}
}

func TestSliceMutex_Append(t *testing.T) {
	s := NewMutexSlice([]int{1, 2})
	s.Append(3, 4)
	if !slices.Equal(s.Load(), []int{1, 2, 3, 4}) {
		t.Errorf("expected [1 2 3 4], got %v", s.Load())
	}
}

func TestSliceMutex_Unshift(t *testing.T) {
	s := NewMutexSlice([]int{1, 2})
	s.Unshift(0)
	if !slices.Equal(s.Load(), []int{0, 1, 2}) {
		t.Errorf("expected [0 1 2], got %v", s.Load())
	}
}

func TestSliceMutex_Shift(t *testing.T) {
	s := NewMutexSlice([]int{1, 2})
	val := s.Shift()
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
	if !slices.Equal(s.Load(), []int{2}) {
		t.Errorf("expected [2], got %v", s.Load())
	}
}

func TestSliceMutex_Pop(t *testing.T) {
	s := NewMutexSlice([]int{1, 2})
	val := s.Pop()
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
	if !slices.Equal(s.Load(), []int{1}) {
		t.Errorf("expected [1], got %v", s.Load())
	}
}

func TestSliceMutex_Clone(t *testing.T) {
	s := NewMutexSlice([]int{1, 2})
	clone := s.Clone()
	if !slices.Equal(clone, []int{1, 2}) {
		t.Errorf("expected [1 2], got %v", clone)
	}
}

func TestSliceMutex_Len(t *testing.T) {
	s := NewMutexSlice([]int{1, 2, 3})
	if s.Len() != 3 {
		t.Errorf("expected 3, got %d", s.Len())
	}
}

func TestSliceMutex_IsEmpty(t *testing.T) {
	s := NewMutexSlice([]int{})
	if !s.IsEmpty() {
		t.Error("expected true, got false")
	}
	s.Append(1)
	if s.IsEmpty() {
		t.Error("expected false, got true")
	}
}

func TestSliceMutex_Get(t *testing.T) {
	s := NewMutexSlice([]int{1, 2, 3})
	if s.Get(1) != 2 {
		t.Errorf("expected 2, got %d", s.Get(1))
	}
}

func TestSliceMutex_Remove(t *testing.T) {
	s := NewMutexSlice([]int{1, 2, 3})
	val := s.Remove(1)
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
	if !slices.Equal(s.Load(), []int{1, 3}) {
		t.Errorf("expected [1 3], got %v", s.Load())
	}
}

func TestSliceMutex_Insert(t *testing.T) {
	s := NewMutexSlice([]int{1, 3})
	s.Insert(1, 2)
	if !slices.Equal(s.Load(), []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", s.Load())
	}
}

func TestSliceMutex_Filter(t *testing.T) {
	s := NewMutexSlice([]int{1, 2, 3, 4})
	filtered := s.Filter(func(v int) bool { return v%2 == 0 })
	if !slices.Equal(filtered, []int{2, 4}) {
		t.Errorf("expected [2 4], got %v", filtered)
	}
	if !slices.Equal(s.Load(), []int{1, 2, 3, 4}) {
		t.Errorf("expected [1 2 3 4], got %v", s.Load())
	}
}

func TestMapMutex_Insert(t *testing.T) {
	m := NewMutexMap(map[string]int{})
	m.Insert("a", 1)
	if m.Load()["a"] != 1 {
		t.Errorf("expected 1, got %d", m.Load()["a"])
	}
}

func TestMapMutex_Get(t *testing.T) {
	m := NewMutexMap(map[string]int{"a": 1})
	val, ok := m.Get("a")
	if !ok {
		t.Error("expected true, got false")
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
}

func TestMapMutex_Remove(t *testing.T) {
	m := NewMutexMap(map[string]int{"a": 1})
	val, ok := m.Remove("a")
	if !ok {
		t.Error("expected true, got false")
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
	if m.ContainsKey("a") {
		t.Error("expected false, got true")
	}
}

func TestMapMutex_Keys(t *testing.T) {
	m := NewMutexMap(map[string]int{"a": 1, "b": 2})
	keys := m.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2, got %d", len(keys))
	}
	if !slices.Contains(keys, "a") || !slices.Contains(keys, "b") {
		t.Errorf("expected keys to contain 'a' and 'b', got %v", keys)
	}
}

func TestNumberMutex_Add(t *testing.T) {
	n := NewMutexNumber(10)
	n.Add(5)
	if n.Load() != 15 {
		t.Errorf("expected 15, got %d", n.Load())
	}
}

func TestNumberMutex_Sub(t *testing.T) {
	n := NewMutexNumber(10)
	n.Sub(5)
	if n.Load() != 5 {
		t.Errorf("expected 5, got %d", n.Load())
	}
}

func TestSliceMutex_Each(t *testing.T) {
	s := NewMutexSlice([]int{1, 2, 3})
	var sum int
	s.Each(func(v int) {
		sum += v
	})
	if sum != 6 {
		t.Errorf("expected 6, got %d", sum)
	}
}

func TestSliceMutex_Clear(t *testing.T) {
	s := NewMutexSlice([]int{1, 2, 3})
	s.Clear()
	if !slices.Equal(s.Load(), []int{}) {
		t.Errorf("expected empty slice, got %v", s.Load())
	}
	if !s.IsEmpty() {
		t.Error("expected true, got false")
	}
}

func TestMapMutex_Clear(t *testing.T) {
	m := NewMutexMap(map[string]int{"a": 1, "b": 2})
	m.Clear()
	if len(m.Load()) != 0 {
		t.Errorf("expected empty map, got %v", m.Load())
	}
	if !m.IsEmpty() {
		t.Error("expected true, got false")
	}
}

func TestMapMutex_GetKeyValue(t *testing.T) {
	m := NewMutexMap(map[string]int{"a": 1})
	k, v, ok := m.GetKeyValue("a")
	if !ok {
		t.Error("expected true, got false")
	}
	if k != "a" {
		t.Errorf("expected 'a', got %q", k)
	}
	if v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	_, _, ok = m.GetKeyValue("b")
	if ok {
		t.Error("expected false, got true")
	}
}

func TestMapMutex_Delete(t *testing.T) {
	m := NewMutexMap(map[string]int{"a": 1})
	m.Delete("a")
	if m.ContainsKey("a") {
		t.Error("expected false, got true")
	}
}

func TestMapMutex_Len(t *testing.T) {
	m := NewMutexMap(map[string]int{"a": 1, "b": 2})
	if m.Len() != 2 {
		t.Errorf("expected 2, got %d", m.Len())
	}
}

func TestMapMutex_IsEmpty(t *testing.T) {
	m := NewMutexMap(map[string]int{})
	if !m.IsEmpty() {
		t.Error("expected true, got false")
	}
	m.Insert("a", 1)
	if m.IsEmpty() {
		t.Error("expected false, got true")
	}
}

func TestMapMutex_Each(t *testing.T) {
	m := NewMutexMap(map[string]int{"a": 1, "b": 2})
	var sum int
	m.Each(func(k string, v int) {
		sum += v
	})
	if sum != 3 {
		t.Errorf("expected 3, got %d", sum)
	}
}

func TestMapMutex_Values(t *testing.T) {
	m := NewMutexMap(map[string]int{"a": 1, "b": 2})
	values := m.Values()
	if len(values) != 2 {
		t.Errorf("expected 2, got %d", len(values))
	}
	if !slices.Contains(values, 1) || !slices.Contains(values, 2) {
		t.Errorf("expected values to contain 1 and 2, got %v", values)
	}
}

func TestMapMutex_Clone(t *testing.T) {
	m := NewMutexMap(map[string]int{"a": 1})
	clone := m.Clone()
	if clone["a"] != 1 {
		t.Errorf("expected 1, got %d", clone["a"])
	}
	m.Insert("a", 2)
	if clone["a"] != 1 {
		t.Errorf("expected 1, got %d", clone["a"])
	}
}

func TestMapRWMutex_Clear(t *testing.T) {
	m := NewRWMutexMap(map[string]int{"a": 1})
	m.Clear()
	if len(m.Load()) != 0 {
		t.Errorf("expected empty map, got %v", m.Load())
	}
	if !m.IsEmpty() {
		t.Error("expected true, got false")
	}
}

func TestMapRWMutex_Insert(t *testing.T) {
	m := NewRWMutexMap(map[string]int{})
	m.Insert("a", 1)
	if m.Load()["a"] != 1 {
		t.Errorf("expected 1, got %d", m.Load()["a"])
	}
}

func TestMapRWMutex_Get(t *testing.T) {
	m := NewRWMutexMap(map[string]int{"a": 1})
	val, ok := m.Get("a")
	if !ok {
		t.Error("expected true, got false")
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	_, ok = m.Get("b")
	if ok {
		t.Error("expected false, got true")
	}
}

func TestMapRWMutex_GetKeyValue(t *testing.T) {
	m := NewRWMutexMap(map[string]int{"a": 1})
	k, v, ok := m.GetKeyValue("a")
	if !ok {
		t.Error("expected true, got false")
	}
	if k != "a" {
		t.Errorf("expected 'a', got %q", k)
	}
	if v != 1 {
		t.Errorf("expected 1, got %d", v)
	}

	_, _, ok = m.GetKeyValue("b")
	if ok {
		t.Error("expected false, got true")
	}
}

func TestMapRWMutex_ContainsKey(t *testing.T) {
	m := NewRWMutexMap(map[string]int{"a": 1})
	if !m.ContainsKey("a") {
		t.Error("expected true, got false")
	}
	if m.ContainsKey("b") {
		t.Error("expected false, got true")
	}
}

func TestMapRWMutex_Remove(t *testing.T) {
	m := NewRWMutexMap(map[string]int{"a": 1})
	val, ok := m.Remove("a")
	if !ok {
		t.Error("expected true, got false")
	}
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
	if m.ContainsKey("a") {
		t.Error("expected false, got true")
	}

	_, ok = m.Remove("a")
	if ok {
		t.Error("expected false, got true")
	}
}

func TestMapRWMutex_Delete(t *testing.T) {
	m := NewRWMutexMap(map[string]int{"a": 1})
	m.Delete("a")
	if m.ContainsKey("a") {
		t.Error("expected false, got true")
	}
}

func TestMapRWMutex_Len(t *testing.T) {
	m := NewRWMutexMap(map[string]int{"a": 1, "b": 2})
	if m.Len() != 2 {
		t.Errorf("expected 2, got %d", m.Len())
	}
}

func TestMapRWMutex_IsEmpty(t *testing.T) {
	m := NewRWMutexMap(map[string]int{})
	if !m.IsEmpty() {
		t.Error("expected true, got false")
	}
	m.Insert("a", 1)
	if m.IsEmpty() {
		t.Error("expected false, got true")
	}
}

func TestMapRWMutex_Each(t *testing.T) {
	m := NewRWMutexMap(map[string]int{"a": 1, "b": 2})
	var sum int
	m.Each(func(k string, v int) {
		sum += v
	})
	if sum != 3 {
		t.Errorf("expected 3, got %d", sum)
	}
}

func TestMapRWMutex_Keys(t *testing.T) {
	m := NewRWMutexMap(map[string]int{"a": 1, "b": 2})
	keys := m.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2, got %d", len(keys))
	}
	if !slices.Contains(keys, "a") || !slices.Contains(keys, "b") {
		t.Errorf("expected keys to contain 'a' and 'b', got %v", keys)
	}
}

func TestMapRWMutex_Values(t *testing.T) {
	m := NewRWMutexMap(map[string]int{"a": 1, "b": 2})
	values := m.Values()
	if len(values) != 2 {
		t.Errorf("expected 2, got %d", len(values))
	}
	if !slices.Contains(values, 1) || !slices.Contains(values, 2) {
		t.Errorf("expected values to contain 1 and 2, got %v", values)
	}
}

func TestMapRWMutex_Clone(t *testing.T) {
	m := NewRWMutexMap(map[string]int{"a": 1})
	clone := m.Clone()
	if clone["a"] != 1 {
		t.Errorf("expected 1, got %d", clone["a"])
	}
	m.Insert("a", 2)
	if clone["a"] != 1 {
		t.Errorf("expected 1, got %d", clone["a"])
	}
}

func TestSliceRWMutex_Each(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 2, 3})
	var sum int
	s.Each(func(v int) {
		sum += v
	})
	if sum != 6 {
		t.Errorf("expected 6, got %d", sum)
	}
}

func TestSliceRWMutex_Clear(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 2, 3})
	s.Clear()
	if !slices.Equal(s.Load(), []int{}) {
		t.Errorf("expected empty slice, got %v", s.Load())
	}
	if !s.IsEmpty() {
		t.Error("expected true, got false")
	}
}

func TestSliceRWMutex_Append(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 2})
	s.Append(3, 4)
	if !slices.Equal(s.Load(), []int{1, 2, 3, 4}) {
		t.Errorf("expected [1 2 3 4], got %v", s.Load())
	}
}

func TestSliceRWMutex_Unshift(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 2})
	s.Unshift(0)
	if !slices.Equal(s.Load(), []int{0, 1, 2}) {
		t.Errorf("expected [0 1 2], got %v", s.Load())
	}
}

func TestSliceRWMutex_Shift(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 2})
	val := s.Shift()
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
	if !slices.Equal(s.Load(), []int{2}) {
		t.Errorf("expected [2], got %v", s.Load())
	}
}

func TestSliceRWMutex_Pop(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 2})
	val := s.Pop()
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
	if !slices.Equal(s.Load(), []int{1}) {
		t.Errorf("expected [1], got %v", s.Load())
	}
}

func TestSliceRWMutex_Clone(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 2})
	clone := s.Clone()
	if !slices.Equal(clone, []int{1, 2}) {
		t.Errorf("expected [1 2], got %v", clone)
	}
	s.Append(3)
	if !slices.Equal(clone, []int{1, 2}) {
		t.Errorf("expected [1 2], got %v", clone)
	}
}

func TestSliceRWMutex_Len(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 2, 3})
	if s.Len() != 3 {
		t.Errorf("expected 3, got %d", s.Len())
	}
}

func TestSliceRWMutex_IsEmpty(t *testing.T) {
	s := NewRWMutexSlice([]int{})
	if !s.IsEmpty() {
		t.Error("expected true, got false")
	}
	s.Append(1)
	if s.IsEmpty() {
		t.Error("expected false, got true")
	}
}

func TestSliceRWMutex_Get(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 2, 3})
	if s.Get(1) != 2 {
		t.Errorf("expected 2, got %d", s.Get(1))
	}
}

func TestSliceRWMutex_Remove(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 2, 3})
	val := s.Remove(1)
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
	if !slices.Equal(s.Load(), []int{1, 3}) {
		t.Errorf("expected [1 3], got %v", s.Load())
	}
}

func TestSliceRWMutex_Insert(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 3})
	s.Insert(1, 2)
	if !slices.Equal(s.Load(), []int{1, 2, 3}) {
		t.Errorf("expected [1 2 3], got %v", s.Load())
	}
}

func TestSliceRWMutex_Filter(t *testing.T) {
	s := NewRWMutexSlice([]int{1, 2, 3, 4})
	filtered := s.Filter(func(v int) bool { return v%2 == 0 })
	if !slices.Equal(filtered, []int{2, 4}) {
		t.Errorf("expected [2 4], got %v", filtered)
	}
	if !slices.Equal(s.Load(), []int{1, 2, 3, 4}) {
		t.Errorf("expected [1 2 3 4], got %v", s.Load())
	}
}

func TestNumberRWMutex_Add(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		n := NewRWMutexNumber[int](10)
		n.Add(5)
		if n.Load() != 15 {
			t.Errorf("expected 15, got %d", n.Load())
		}
	})

	t.Run("float64", func(t *testing.T) {
		n := NewRWMutexNumber[float64](10.5)
		n.Add(2.5)
		if n.Load() != 13.0 {
			t.Errorf("expected 13.0, got %f", n.Load())
		}
	})

	t.Run("uint", func(t *testing.T) {
		n := NewRWMutexNumber[uint](10)
		n.Add(5)
		if n.Load() != 15 {
			t.Errorf("expected 15, got %d", n.Load())
		}
	})
}

func TestNumberRWMutex_Sub(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		n := NewRWMutexNumber[int](10)
		n.Sub(3)
		if n.Load() != 7 {
			t.Errorf("expected 7, got %d", n.Load())
		}
	})

	t.Run("float64", func(t *testing.T) {
		n := NewRWMutexNumber[float64](10.5)
		n.Sub(2.5)
		if n.Load() != 8.0 {
			t.Errorf("expected 8.0, got %f", n.Load())
		}
	})

	t.Run("uint", func(t *testing.T) {
		n := NewRWMutexNumber[uint](10)
		n.Sub(3)
		if n.Load() != 7 {
			t.Errorf("expected 7, got %d", n.Load())
		}
	})
}

func TestNumberRWMutex_ConcurrentOperations(t *testing.T) {
	n := NewRWMutexNumber(0)
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
	if n.Load() != 0 {
		t.Errorf("expected 0, got %d", n.Load())
	}
}

func TestNewMapMutex(t *testing.T) {
	t.Run("creates new MutexMap with given map", func(t *testing.T) {
		input := map[string]int{"a": 1, "b": 2}
		m := NewMutexMap(input)

		// Verify the internal map matches input
		m.RWith(func(v map[string]int) {
			if len(v) != 2 {
				t.Errorf("expected length 2, got %d", len(v))
			}
			if v["a"] != 1 || v["b"] != 2 {
				t.Errorf("map contents don't match input")
			}
		})
	})
}

func TestNewMapRWMutex(t *testing.T) {
	t.Run("creates new RWMutexMap with given map", func(t *testing.T) {
		input := map[int]string{1: "one", 2: "two"}
		m := NewRWMutexMap(input)

		m.RWith(func(v map[int]string) {
			if len(v) != 2 {
				t.Errorf("expected length 2, got %d", len(v))
			}
			if v[1] != "one" || v[2] != "two" {
				t.Errorf("map contents don't match input")
			}
		})
	})
}

func TestNewRWMutex(t *testing.T) {
	t.Run("creates new RWMutex with given value", func(t *testing.T) {
		initialValue := "test"
		rwMutex := NewRWMutex(initialValue)

		// Verify the initial value is stored correctly
		rwMutex.RWith(func(v string) {
			if v != initialValue {
				t.Errorf("expected %q, got %q", initialValue, v)
			}
		})

		// Test mutex operations
		newValue := "updated"
		rwMutex.With(func(v *string) {
			*v = newValue
		})

		rwMutex.RWith(func(v string) {
			if v != newValue {
				t.Errorf("expected %q, got %q", newValue, v)
			}
		})
	})
}

func TestNewSliceMutex(t *testing.T) {
	t.Run("creates new MutexSlice with given slice", func(t *testing.T) {
		initialSlice := []int{1, 2, 3}
		sliceMutex := NewMutexSlice(initialSlice)

		// Verify the initial slice is stored correctly
		sliceMutex.RWith(func(v []int) {
			if len(v) != len(initialSlice) {
				t.Errorf("expected length %d, got %d", len(initialSlice), len(v))
			}
			for i, val := range v {
				if val != initialSlice[i] {
					t.Errorf("expected %d at index %d, got %d", initialSlice[i], i, val)
				}
			}
		})

		// Test mutex operations
		sliceMutex.With(func(v *[]int) {
			*v = append(*v, 4)
		})

		sliceMutex.RWith(func(v []int) {
			if len(v) != 4 {
				t.Errorf("expected length 4, got %d", len(v))
			}
			if v[3] != 4 {
				t.Errorf("expected 4 at index 3, got %d", v[3])
			}
		})
	})
}

func TestNewSliceRWMutex(t *testing.T) {
	t.Run("creates new RWMutexSlice with given slice", func(t *testing.T) {
		initialSlice := []string{"a", "b", "c"}
		sliceRWMutex := NewRWMutexSlice(initialSlice)

		// Verify the initial slice is stored correctly
		sliceRWMutex.RWith(func(v []string) {
			if len(v) != len(initialSlice) {
				t.Errorf("expected length %d, got %d", len(initialSlice), len(v))
			}
			for i, val := range v {
				if val != initialSlice[i] {
					t.Errorf("expected %q at index %d, got %q", initialSlice[i], i, val)
				}
			}
		})

		// Test mutex operations
		sliceRWMutex.With(func(v *[]string) {
			*v = append(*v, "d")
		})

		sliceRWMutex.RWith(func(v []string) {
			if len(v) != 4 {
				t.Errorf("expected length 4, got %d", len(v))
			}
			if v[3] != "d" {
				t.Errorf("expected \"d\" at index 3, got %q", v[3])
			}
		})
	})
}

func TestNewNumberMutex(t *testing.T) {
	t.Run("creates new MutexNumber with given number", func(t *testing.T) {
		initialNumber := 42
		numberMutex := NewMutexNumber(initialNumber)

		// Verify the initial number is stored correctly
		numberMutex.RWith(func(v int) {
			if v != initialNumber {
				t.Errorf("expected %d, got %d", initialNumber, v)
			}
		})

		// Test mutex operations
		numberMutex.With(func(v *int) {
			*v += 10
		})

		numberMutex.RWith(func(v int) {
			if v != 52 {
				t.Errorf("expected 52, got %d", v)
			}
		})
	})
}

func TestNewNumberRWMutex(t *testing.T) {
	t.Run("creates new RWMutexNumber with given number", func(t *testing.T) {
		initialNumber := 3.14
		numberRWMutex := NewRWMutexNumber(initialNumber)

		// Verify the initial number is stored correctly
		numberRWMutex.RWith(func(v float64) {
			if v != initialNumber {
				t.Errorf("expected %f, got %f", initialNumber, v)
			}
		})

		// Test mutex operations
		numberRWMutex.With(func(v *float64) {
			*v *= 2
		})

		numberRWMutex.RWith(func(v float64) {
			if v != 6.28 {
				t.Errorf("expected 6.28, got %f", v)
			}
		})
	})
}

func TestNewMutex(t *testing.T) {
	t.Run("creates new Mutex with given value", func(t *testing.T) {
		initialValue := "initial"
		mutex := NewMutex(initialValue)

		// Verify the initial value is stored correctly
		mutex.RWith(func(v string) {
			if v != initialValue {
				t.Errorf("expected %q, got %q", initialValue, v)
			}
		})

		// Test mutex operations
		newValue := "updated"
		mutex.With(func(v *string) {
			*v = newValue
		})

		mutex.RWith(func(v string) {
			if v != newValue {
				t.Errorf("expected %q, got %q", newValue, v)
			}
		})
	})

	t.Run("works with numeric types", func(t *testing.T) {
		initialValue := 42
		mutex := NewMutex(initialValue)

		// Verify the initial value is stored correctly
		mutex.RWith(func(v int) {
			if v != initialValue {
				t.Errorf("expected %d, got %d", initialValue, v)
			}
		})

		// Test mutex operations
		mutex.With(func(v *int) {
			*v += 10
		})

		mutex.RWith(func(v int) {
			if v != 52 {
				t.Errorf("expected 52, got %d", v)
			}
		})
	})

	t.Run("works with struct types", func(t *testing.T) {
		type testStruct struct {
			Field1 string
			Field2 int
		}

		initialValue := testStruct{"hello", 42}
		mutex := NewMutex(initialValue)

		// Verify the initial value is stored correctly
		mutex.RWith(func(v testStruct) {
			if v.Field1 != "hello" || v.Field2 != 42 {
				t.Errorf("expected {hello 42}, got %v", v)
			}
		})

		// Test mutex operations
		mutex.With(func(v *testStruct) {
			v.Field1 = "world"
			v.Field2 = 100
		})

		mutex.RWith(func(v testStruct) {
			if v.Field1 != "world" || v.Field2 != 100 {
				t.Errorf("expected {world 100}, got %v", v)
			}
		})
	})
}
