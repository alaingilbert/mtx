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
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestMtx_Debug(t *testing.T) {
	debug = true
	defer func() { debug = false }()
	m := NewMap[string, int](nil)
	m.Len()
	m2 := NewRWMap[string, int](nil)
	m2.Len()
	assert.Equal(t, 1, 1)
}

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

func TestMtx_Replace(t *testing.T) {
	m := NewMtx("old")
	old := m.Replace("new")
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

func TestMtxPtr_Replace(t *testing.T) {
	m := NewMtxPtr("old")
	old := m.Replace("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Load())
}

func TestMtx_MarshalJSON(t *testing.T) {
	// Mtx ptr; works
	var tmp1 = struct{ Field *Mtx[[]int] }{Field: NewMtxPtr([]int{1, 2, 3})}
	out1, _ := json.Marshal(tmp1)
	assert.Equal(t, `{"Field":[1,2,3]}`, string(out1))
}

func TestRWMtx_RLockRUnlock(t *testing.T) {
	m := NewRWMtx("old")
	m.RLock()
	val := m.GetPointer()
	assert.Equal(t, "old", *val)
	m.RUnlock()
}

func TestRWMtx_Replace(t *testing.T) {
	m := NewRWMtx("old")
	old := m.Replace("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Load())
}

func TestRWMtxPtr_Replace(t *testing.T) {
	m := NewRWMtxPtr("old")
	old := m.Replace("new")
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
	m.SetKey("a", 1)
	el, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, el)
}

func TestMap_GetKeyValue(t *testing.T) {
	m := NewMap[string, int](nil)
	_, _, ok := m.GetKeyValue("a")
	assert.False(t, ok)
	m.SetKey("a", 1)
	key, value, ok := m.GetKeyValue("a")
	assert.True(t, ok)
	assert.Equal(t, "a", key)
	assert.Equal(t, 1, value)
}

func TestMap_HasKey(t *testing.T) {
	m := NewMap[string, int](nil)
	assert.False(t, m.ContainsKey("a"))
	m.SetKey("a", 1)
	assert.True(t, m.ContainsKey("a"))
	m.Remove("a")
	assert.False(t, m.ContainsKey("a"))
}

func TestMap_Take(t *testing.T) {
	m := NewMap[string, int](nil)
	m.SetKey("a", 1)
	m.SetKey("b", 2)
	m.SetKey("c", 3)
	assert.Equal(t, 3, m.Len())
	assert.True(t, m.ContainsKey("b"))
	val, ok := m.Take("b")
	assert.True(t, ok)
	assert.Equal(t, 2, val)
	assert.False(t, m.ContainsKey("b"))
	_, ok = m.Take("b")
	assert.False(t, ok)
	assert.Equal(t, 2, m.Len())
}

func TestMap_Remove(t *testing.T) {
	m := NewMap[string, int](nil)
	assert.Equal(t, 0, m.Len())
	m.Remove("a")
	m.SetKey("a", 1)
	assert.Equal(t, 1, m.Len())
	m.Remove("a")
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
	m.SetKey("a", 1)
	el, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, 1, el)
}

func TestRWMap_HasKey(t *testing.T) {
	m := NewRWMap[string, int](nil)
	assert.False(t, m.ContainsKey("a"))
	m.SetKey("a", 1)
	assert.True(t, m.ContainsKey("a"))
	m.Remove("a")
	assert.False(t, m.ContainsKey("a"))
}

func TestRWMap_Take(t *testing.T) {
	m := NewRWMap[string, int](nil)
	m.SetKey("a", 1)
	m.SetKey("b", 2)
	m.SetKey("c", 3)
	assert.Equal(t, 3, m.Len())
	assert.True(t, m.ContainsKey("b"))
	val, ok := m.Take("b")
	assert.True(t, ok)
	assert.Equal(t, 2, val)
	assert.False(t, m.ContainsKey("b"))
	_, ok = m.Take("b")
	assert.False(t, ok)
	assert.Equal(t, 2, m.Len())
}

func TestRWMap_Remove(t *testing.T) {
	m := NewRWMap[string, int](nil)
	assert.Equal(t, 0, m.Len())
	m.Remove("a")
	m.SetKey("a", 1)
	assert.Equal(t, 1, m.Len())
	m.Remove("a")
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
	m.SetKey("a", 4)
	assert.Equal(t, 4, theMap["a"])
	assert.Equal(t, 4, first(m.Get("a")))
	theMap["a"] = 5
	assert.Equal(t, 5, theMap["a"])
	assert.Equal(t, 5, first(m.Get("a")))
}

func TestRWMap_Clone(t *testing.T) {
	m := NewRWMap(map[string]int{"a": 1, "b": 2, "c": 3})
	clonedMap := m.Clone()
	m.SetKey("a", 4)
	assert.Equal(t, 1, clonedMap["a"])
	assert.Equal(t, 4, first(m.Get("a")))
}

func TestMap_MarshalJSON(t *testing.T) {
	// Mtx ptr; works
	var tmp1 = struct{ Field *Map[string, int] }{Field: NewRWMapPtr(map[string]int{"a": 1})}
	out1, _ := json.Marshal(tmp1)
	assert.Equal(t, `{"Field":{"a":1}}`, string(out1))

	// Mtx ptr and give addr to json.Marshal; works
	var tmp2 = struct{ Field *Map[string, int] }{Field: NewMapPtr(map[string]int{"a": 1})}
	out2, _ := json.Marshal(&tmp2)
	assert.Equal(t, `{"Field":{"a":1}}`, string(out2))

	// Mtx is not a pointer, but give address to struct to json.Marshal; works
	var tmp3 = struct{ Field Map[string, int] }{Field: NewRWMap(map[string]int{"a": 1})}
	out3, _ := json.Marshal(&tmp3)
	assert.Equal(t, `{"Field":{"a":1}}`, string(out3))

	// Somehow, the Locker field get marshalled if we do not provide an address to the struct
	// or if the mutex itself is not a pointer
	var tmp4 = struct{ Field Map[string, int] }{Field: NewRWMap(map[string]int{"a": 1})}
	out4, err := json.Marshal(tmp4)
	assert.NoError(t, err)
	assert.Equal(t, `{"Field":{"Locker":{"a":1}}}`, string(out4))

	var tmp5 = struct{ Field Map[string, int] }{Field: NewMap(map[string]int{"a": 1})}
	out5, _ := json.Marshal(&tmp5)
	assert.Equal(t, `{"Field":{"a":1}}`, string(out5))

	var tmp6 = struct{ Field *Map[string, int] }{Field: NewMapPtr(map[string]int{"a": 1})}
	out6, _ := json.Marshal(tmp6)
	assert.Equal(t, `{"Field":{"a":1}}`, string(out6))
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
	m.Remove(1)
	assert.Equal(t, []int{4}, m.Load())
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

func TestSlice_MarshalJSON(t *testing.T) {
	// Mtx ptr; works
	var tmp1 = struct{ Field *Slice[int] }{Field: NewRWSlicePtr([]int{1, 2, 3})}
	out1, _ := json.Marshal(tmp1)
	assert.Equal(t, `{"Field":[1,2,3]}`, string(out1))

	// Mtx ptr and give addr to json.Marshal; works
	var tmp2 = struct{ Field *Slice[int] }{Field: NewSlicePtr([]int{1, 2, 3})}
	out2, _ := json.Marshal(&tmp2)
	assert.Equal(t, `{"Field":[1,2,3]}`, string(out2))

	// Mtx is not a pointer, but give address to struct to json.Marshal; works
	var tmp3 = struct{ Field Slice[int] }{Field: NewRWSlice([]int{1, 2, 3})}
	out3, _ := json.Marshal(&tmp3)
	assert.Equal(t, `{"Field":[1,2,3]}`, string(out3))
}

func TestRWUInt64(t *testing.T) {
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
