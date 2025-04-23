package mtx

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsage(t *testing.T) {
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
