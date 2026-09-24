package vector

import (
	"testing"
)

func TestNewVectorHasZeroLength(t *testing.T) {
	vector := New()
	if vector.length != 0 { // If new vector is not zero length
		t.Errorf("Expected length of newVector to be 0, but got %d", vector.length)
	}
}

func TestPush(t *testing.T) {
	vector := New()
	vector.Push(1)
	if vector.length != 1 {
		t.Errorf("Expected length of vector to be 1, but got %d", vector.length)
	}
	if value, _ := vector.Get(0); value != 1 {
		t.Errorf("Expected vector.Get(0) to be 1, but got %d", value)
	}
}

func TestCapacity(t *testing.T) {
	vector := NewWithCapacity(2)
	vector.Push(1)
	vector.Push(2)
	vector.Push(3)
	if value, _ := vector.Get(0); value != 1 { // Test capacity increase conserves order
		t.Errorf("Expected vector.Get(0) to be 1, but got %v", value)
	}
	if value, _ := vector.Get(1); value != 2 {
		t.Errorf("Expected vector.Get(1) to be 2, but got %v", value)
	}
	if vector.Cap() != 4 { // Check capacity doubles after push
		t.Errorf("Expected capacity of vector to be 4, but got %d", vector.Cap())
	}
}

func TestPop(t *testing.T) {
	vector := New()
	vector.Push(1)
	if value, _ := vector.Pop(); value != 1 {
		t.Errorf("Expected vector.Pop() to be 2, but got %d", value)
	}
	if vector.length != 0 {
		t.Errorf("Expected length of vector to be 1 after pop, but got %d", vector.length)
	}
	if _, ok := vector.Pop(); ok {
		t.Errorf("Expected vector.Pop() on empty vector to return pop, but got %t", ok)
	}
}

func TestGet(t *testing.T) {
	vector := New()
	vector.Push(1)
	if _, ok := vector.Get(-1); ok {
		t.Errorf("Expected vector.Get(-1) on vector with length 1 to return false, but got %t", ok)
	}
	if _, ok := vector.Get(1); ok {
		t.Errorf("Expected vector.Get(1) on vector with length 1 to return false, but got %t", ok)
	}
}

func TestSet(t *testing.T) {
	vector := New()
	vector.Push(1)
	vector.Push(2)
	vector.Set(0, 3)
	if value, _ := vector.Get(0); value != 3 {
		t.Errorf("Expected vector.Get(0) to be 3, but got %d", value)
	}
	if vector.length != 2 {
		t.Errorf("Expected length of vector to be 2 after set, but got %d", vector.length)
	}
	vector.Set(3, 4)
	if vector.Set(3, 4) == true {
		t.Errorf("Expected vector.Set(3, 4) on vector with length 2 to return false, but got %t", vector.Set(3, 4))
	}
}

func TestInsert(t *testing.T) {
	vector := New()
	vector.Push(1)
	vector.Push(2)
	vector.Insert(1, 3)
	if value, _ := vector.Get(1); value != 3 {
		t.Errorf("Expected vector.Get(1) to be 3, but got %d", value)
	}
	if vector.length != 3 {
		t.Errorf("Expected length of vector to be 3 after insert, but got %d", vector.length)
	}
	if value, _ := vector.Get(2); value != 2 {
		t.Errorf("Expected vector.Get(2) to be 2, but got %d", value)
	}
}

func TestInsertEdgeCases(t *testing.T) {
	vector := New()
	vector.Push(1)
	vector.Insert(0, 2)
	if value, _ := vector.Get(0); value != 2 { // Test insert to start
		t.Errorf("Expected vector.Get(0) to be 1, but got %d", value)
	}
	vector.Insert(2, 3)
	if value, _ := vector.Get(2); value != 3 { // Test insert to end
		t.Errorf("Expected vector.Get(2) to be 3, but got %d", value)
	}
	if vector.length != 3 {
		t.Errorf("Expected length of vector to be 3 after insert, but got %d", vector.length)
	}
}

func TestRemove(t *testing.T) {
	vector := New()
	vector.Push(1)
	vector.Push(2)
	vector.Push(3)
	vector.Remove(1)
	if value, _ := vector.Get(1); value != 3 {
		t.Errorf("Expected vector.Get(1) to be 3, but got %d", value)
	}
	if vector.length != 2 {
		t.Errorf("Expected length of vector to be 2 after remove, but got %d", vector.length)
	}
}

func TestCapacityGrowth(t *testing.T) {
	vector := New()
	for i := range 1000 {
		vector.Push(i)
		vector.Pop()
	}
	if vector.Cap() != 1 {
		t.Errorf("Expected capacity of vector to be 0 after 1000 pushes and pops, but got %d", vector.Cap())
	}
}
