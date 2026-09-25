package linkedlist

import "testing"

func newTestLinkedList(t *testing.T, values []int) *LinkedList {
	return nil
}

func TestNewListHasLengthZero(t *testing.T) {
	l := newTestLinkedList(t, []int{})
	if l.Len() != 0 {
		t.Errorf("New linked list length %d, want 0", l.Len())
	}
}

func TestPushPopRoundTrips(t *testing.T) {

}
