package mtx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func TestMtx_Set(t *testing.T) {
	m := NewMtx("old")
	assert.Equal(t, "old", m.Get())
	m.Set("new")
	assert.Equal(t, "new", m.Get())
}

func TestMtx_Replace(t *testing.T) {
	m := NewMtx("old")
	old := m.Replace("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Get())
}

func TestMtx_Val(t *testing.T) {
	someString := "old"
	orig := &someString
	m := NewMtx(orig)
	val := m.Val()
	**val = "new"
	assert.Equal(t, "new", someString)
	assert.Equal(t, "new", **val)
	assert.Equal(t, "new", *orig)
}

func TestMtxPtr_Replace(t *testing.T) {
	m := NewMtxPtr("old")
	old := m.Replace("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Get())
}

func TestRWMtx_Replace(t *testing.T) {
	m := NewRWMtx("old")
	old := m.Replace("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Get())
}

func TestRWMtxPtr_Replace(t *testing.T) {
	m := NewRWMtxPtr("old")
	old := m.Replace("new")
	assert.Equal(t, "old", old)
	assert.Equal(t, "new", m.Get())
}

func TestRWMtx_Val(t *testing.T) {
	someString := "old"
	orig := &someString
	m := NewRWMtx(orig)
	val := m.Val()
	**val = "new"
	assert.Equal(t, "new", **val)
	assert.Equal(t, "new", *orig)
}

func TestRWMap_GetKey(t *testing.T) {
	m := NewMap[string, int]()
	_, ok := m.GetKey("a")
	assert.False(t, ok)
	m.SetKey("a", 1)
	el, ok := m.GetKey("a")
	assert.True(t, ok)
	assert.Equal(t, 1, el)
}

func TestRWMap_HasKey(t *testing.T) {
	m := NewMap[string, int]()
	assert.False(t, m.HasKey("a"))
	m.SetKey("a", 1)
	assert.True(t, m.HasKey("a"))
	m.DeleteKey("a")
	assert.False(t, m.HasKey("a"))
}

func TestRWMap_TakeKey(t *testing.T) {
	m := NewMap[string, int]()
	m.SetKey("a", 1)
	m.SetKey("b", 2)
	m.SetKey("c", 3)
	assert.Equal(t, 3, m.Len())
	assert.True(t, m.HasKey("b"))
	val, ok := m.TakeKey("b")
	assert.True(t, ok)
	assert.Equal(t, 2, val)
	assert.False(t, m.HasKey("b"))
	_, ok = m.TakeKey("b")
	assert.False(t, ok)
	assert.Equal(t, 2, m.Len())
}

func TestRWMap_DeleteKey(t *testing.T) {
	m := NewMap[string, int]()
	assert.Equal(t, 0, m.Len())
	m.DeleteKey("a")
	m.SetKey("a", 1)
	assert.Equal(t, 1, m.Len())
	m.DeleteKey("a")
	assert.Equal(t, 0, m.Len())
}

func TestRWMap_Values(t *testing.T) {
	m := NewMap[string, int]()
	assert.Equal(t, []int{}, m.Values())
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	values := m.Values()
	slices.Sort(values)
	assert.Equal(t, []int{1, 2, 3}, values)
}

func TestRWMap_Keys(t *testing.T) {
	m := NewMapPtr[string, int]()
	assert.Equal(t, []string{}, m.Keys())
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	keys := m.Keys()
	slices.Sort(keys)
	assert.Equal(t, []string{"a", "b", "c"}, keys)
}

func TestRWMap_Each(t *testing.T) {
	m := NewMap[string, int]()
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	arr := make([]string, 0)
	m.Each(func(k string, v int) {
		arr = append(arr, fmt.Sprintf("%s_%d", k, v))
	})
	slices.Sort(arr)
	assert.Equal(t, []string{"a_1", "b_2", "c_3"}, arr)
}

func TestRWMap_Clone(t *testing.T) {
	m := NewMap[string, int]()
	m.Set(map[string]int{"a": 1, "b": 2, "c": 3})
	clonedMap := m.Clone()
	assert.Equal(t, 1, clonedMap["a"])
}

func TestRWSlice(t *testing.T) {
	m := NewSlice[int]()
	assert.Equal(t, 0, m.Len())
	m.Append(1, 2, 3)
	assert.Equal(t, 3, m.Len())
	assert.Equal(t, []int{1, 2, 3}, m.Get())
	val2 := m.Shift()
	assert.Equal(t, 1, val2)
	m.Unshift(4)
	assert.Equal(t, []int{4, 2, 3}, m.Get())
	val2 = m.Pop()
	assert.Equal(t, []int{4, 2}, m.Get())
	m.DeleteIdx(1)
	assert.Equal(t, []int{4}, m.Get())
	m.Append(5, 6, 7)
	assert.Equal(t, []int{4, 5, 6, 7}, m.Get())
	assert.Equal(t, 6, m.GetIdx(2))
	m.Insert(2, 8)
	assert.Equal(t, []int{4, 5, 8, 6, 7}, m.Get())
}

func TestRWSlice_Clone(t *testing.T) {
	m := NewSlice[int]()
	m.Set([]int{1, 2, 3})
	clonedSlice := m.Clone()
	assert.Equal(t, []int{1, 2, 3}, clonedSlice)
}

func TestRWSlice_Each(t *testing.T) {
	m := NewSlicePtr[int]()
	m.Append(1, 2, 3)
	arr := make([]string, 0)
	m.Each(func(el int) {
		arr = append(arr, fmt.Sprintf("E%d", el))
	})
	assert.Equal(t, []string{"E1", "E2", "E3"}, arr)
}

func TestRWUInt64(t *testing.T) {
	var m RWUInt64[uint64]
	assert.Equal(t, uint64(0), m.Get())
	m.Incr(10)
	assert.Equal(t, uint64(10), m.Get())
	m.Decr(5)
	assert.Equal(t, uint64(5), m.Get())
}
