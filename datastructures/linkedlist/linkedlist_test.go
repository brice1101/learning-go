package linkedlist

import (
	"slices"
	"testing"
)

func newTestLinkedlist(t *testing.T, values []int) Linkedlist {
	t.Helper()
	l := Linkedlist{newTestNode(t, values), nil}
	if l.head != nil {
		if l.head.next == nil {
			l.tail = l.head
		} else {
			n := l.head
			for n.next != nil {
				n = n.next
			}
			l.tail = n
		}
	}
	return l
}

func newTestNode(t *testing.T, values []int) *Node {
	t.Helper()
	if len(values) == 0 {
		return nil
	}
	if len(values) == 1 {
		return &Node{values[0], nil}

	}
	return &Node{values[0], newTestNode(t, values[1:])}
}

func assertValues(t *testing.T, l Linkedlist, want []int) {
	var got []int
	n := l.head
	for n.next != nil {
		got = append(got, n.value)
		n = n.next
	}
	got = append(got, n.value)
	if !slices.Equal(got, want) {
		t.Errorf("Contents %v, want %v", got, want)
	}
}

func TestNewListHasLengthZero(t *testing.T) {
	l := newTestLinkedlist(t, []int{})
	if l.Len() != 0 {
		t.Errorf("New list got len %d, want 0", l.Len())
	}
}

func TestPushPopRoundTrip(t *testing.T) {
	l := newTestLinkedlist(t, []int{1})
	l.PushFront(2)
	l.PopFront()
	assertValues(t, l, []int{1})
}

func TestPushFrontReverseOrder(t *testing.T) {
	l := newTestLinkedlist(t, []int{1})
	for _, value := range []int{2, 3, 4} {
		l.PushFront(value)
	}
	assertValues(t, l, []int{4, 3, 2, 1})
}

func TestPushBackOrder(t *testing.T) {
	l := newTestLinkedlist(t, []int{1})
	for _, value := range []int{2, 3, 4} {
		l.PushBack(value)
	}
	assertValues(t, l, []int{1, 2, 3, 4})
}

func PopFrontOnEmptyErrors(t *testing.T) {
	l := newTestLinkedlist(t, []int{})
	if _, ok := l.PopFront(); ok {
		t.Errorf("Pop on empty got %t, want false", ok)
	}
	assertValues(t, l, []int{})
}

func TestLinkedlist_Remove(t *testing.T) {
	tests := []struct {
		name        string
		values      []int
		removeValue int
		want        []int
		wantOk      bool
	}{
		{name: "head", values: []int{1, 2, 3}, removeValue: 1, want: []int{2, 3}, wantOk: true},
		{name: "tail", values: []int{1, 2, 3}, removeValue: 3, want: []int{1, 2}, wantOk: true},
		{name: "middle", values: []int{1, 2, 3}, removeValue: 2, want: []int{1, 3}, wantOk: true},
		{name: "not in list", values: []int{1, 2, 3}, removeValue: 4, want: []int{1, 2, 3}, wantOk: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := newTestLinkedlist(t, tc.values)
			ok := l.Remove(tc.removeValue)
			if ok != tc.wantOk {
				t.Errorf("Remove() got %t, want %t", ok, tc.wantOk)
			}
			assertValues(t, l, tc.want)
		})
	}
}

func TestRemovePushSingleElement(t *testing.T) {
	l := newTestLinkedlist(t, []int{1})
	l.Remove(1)
	l.PushBack(2)
	assertValues(t, l, []int{2})
}

func TestLinkedlist_Reverse(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
	}{
		{name: "normal", values: []int{1, 2, 3}, want: []int{3, 2, 1}},
		{name: "single element", values: []int{1}, want: []int{1}},
		{name: "empty", values: []int{}, want: []int{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := newTestLinkedlist(t, tc.values)
			l.Reverse()
			assertValues(t, l, tc.want)
		})
	}
}

func TestReverseTwice(t *testing.T) {
	l := newTestLinkedlist(t, []int{1, 2, 3})
	l.Reverse()
	l.Reverse()
	assertValues(t, l, []int{1, 2, 3})
}

func TestTailPointer(t *testing.T) {
	l := newTestLinkedlist(t, []int{1, 2, 3})
	l.Remove(3)
	l.PushBack(4)
	assertValues(t, l, []int{1, 2, 4}) // Test that the new element is reachable from the head
}
