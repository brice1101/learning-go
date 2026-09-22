package vector

import (
	"fmt"
	"testing"
)

func TestNewVectorHasZeroLength(t *testing.T) {

	// Test new vector has length of 0
	newVector := Vector{}
	if newVector.length != 0 {
		t.Errorf("Expected length of newVector to be 0, but got %d", newVector.length)
	}
}

func TestPushGrowsCapacity(t *testing.T) {
	// Test vectors increase capacity when necessary
	pushVector := Vector{}
	pushVector.Push(1)
	pushVector.Push(2)
	pushVector.Push(3)
	if pushVector.Get(0) != 1 {
		t.Errorf("Expected pushVector.Get(0) to be 1, but got %s", fmt.Sprint(pushVector))
	}
	if pushVector.Get(1) != 2 {
		t.Errorf("Expected pushVector.Get(1) to be 2, but got %s", fmt.Sprint(pushVector))
	}
	if pushVector.Get(2) != 3 {
		t.Errorf("Expected pushVector.Get(2) to be 3, but got %s", fmt.Sprint(pushVector))
	}
	if pushVector.capacity != 4 { // doubling growth policy
		t.Errorf("Expected capacity of pushVector to be 4, but got %d", pushVector.capacity)
	}
}

func TestPop(t *testing.T) {
	popVector := Vector{}
	popVector.Push(1)
	popVector.Pop()
	if popVector.length != 0 {
		t.Errorf("Expected popVector.length to be 0, but got %d", popVector.length)
	}
	if popVector.capacity != 1 {
		t.Errorf("Expected popVector.capacity to be 1, but got %d", popVector.capacity)
	}
	if _, ok := popVector.Pop(); ok { // Test popVector.Pop() returns nil when vector is empty
		t.Errorf("Expected popVector.Pop() to return false, but got %s", fmt.Sprint(ok))
	}
}
