package mtx

import (
	"errors"
	"testing"
)

func TestMutexStoreLoad(t *testing.T) {
	var m Mutex[int]
	m.Store(42)
	if got := m.Load(); got != 42 {
		t.Errorf("Load() = %d, want 42", got)
	}
}

func TestMutexSwap(t *testing.T) {
	var m Mutex[string]
	m.Store("hello")
	old := m.Swap("world")
	if old != "hello" {
		t.Errorf("Swap() old = %q, want %q", old, "hello")
	}
	if got := m.Load(); got != "world" {
		t.Errorf("Load() = %q, want %q", got, "world")
	}
}

func TestMutexWith(t *testing.T) {
	var m Mutex[int]
	m.Store(5)
	m.With(func(v *int) {
		*v += 3
	})
	if got := m.Load(); got != 8 {
		t.Errorf("With() result = %d, want 8", got)
	}
}

func TestMutexWithE(t *testing.T) {
	var m Mutex[int]
	err := m.WithE(func(v *int) error {
		*v = 123
		return nil
	})
	if err != nil {
		t.Errorf("WithE() returned unexpected error: %v", err)
	}
	if got := m.Load(); got != 123 {
		t.Errorf("WithE() failed to set value, got %d, want 123", got)
	}
	// Inject an error
	expectedErr := errors.New("some error")
	err = m.WithE(func(v *int) error {
		return expectedErr
	})
	if !errors.Is(expectedErr, err) {
		t.Errorf("WithE() error = %v, want %v", err, expectedErr)
	}
}

func TestMutexRWith(t *testing.T) {
	var m Mutex[int]
	m.Store(9)
	var read int
	m.RWith(func(v int) {
		read = v
	})
	if read != 9 {
		t.Errorf("RWith() = %d, want 9", read)
	}
}

func TestMutexRWithE(t *testing.T) {
	var m Mutex[string]
	m.Store("data")
	err := m.RWithE(func(v string) error {
		if v != "data" {
			return errors.New("unexpected value")
		}
		return nil
	})
	if err != nil {
		t.Errorf("RWithE() returned error: %v", err)
	}
	expectedErr := errors.New("triggered")
	err = m.RWithE(func(v string) error {
		return expectedErr
	})
	if !errors.Is(expectedErr, err) {
		t.Errorf("RWithE() error = %v, want %v", err, expectedErr)
	}
}

func TestMutexRLockUnlock(t *testing.T) {
	var m Mutex[int]
	// Use RLock and RUnlock (which are Lock and Unlock under the hood)
	m.RLock()
	v := m.GetPointer()
	// should read default 0
	if *v != 0 {
		t.Errorf("RWith inside RLock: got %d, want 0", v)
	}
	m.RUnlock()
	// Now try mutating with normal Lock/Unlock to confirm nothing was broken
	m.With(func(v *int) {
		*v = 10
	})
	if got := m.Load(); got != 10 {
		t.Errorf("Load() after RUnlock = %d, want 10", got)
	}
}

func TestRWMutexBasic(t *testing.T) {
	var m RWMutex[int]
	m.Store(100)

	if got := m.Load(); got != 100 {
		t.Errorf("RWMutex Load = %d, want 100", got)
	}
	old := m.Swap(200)
	if old != 100 {
		t.Errorf("RWMutex Swap old = %d, want 100", old)
	}
	if got := m.Load(); got != 200 {
		t.Errorf("RWMutex Load = %d, want 200", got)
	}
}

func TestRWMutexRLockUnlock(t *testing.T) {
	var m RWMutex[int]
	// Use RLock and RUnlock (which are Lock and Unlock under the hood)
	m.RLock()
	v := m.GetPointer()
	// should read default 0
	if *v != 0 {
		t.Errorf("RWith inside RLock: got %d, want 0", v)
	}
	m.RUnlock()
	// Now try mutating with normal Lock/Unlock to confirm nothing was broken
	m.With(func(v *int) {
		*v = 10
	})
	if got := m.Load(); got != 10 {
		t.Errorf("Load() after RUnlock = %d, want 10", got)
	}
}
