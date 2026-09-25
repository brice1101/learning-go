package vector

import (
	"slices"
	"testing"
)

func newTestVector(t *testing.T, capacity int, values []int) *Vector { // Sets up vector fields directly
	t.Helper()
	if len(values) > capacity {
		t.Fatalf("newTestVector: %d values don't fit in capacity %d", len(values), capacity)
	}
	v := &Vector{data: make([]int, capacity), length: len(values)}
	copy(v.data, values)
	return v
}

func assertContents(t *testing.T, v *Vector, want []int) {
	t.Helper()
	got := make([]int, v.Len())
	for i := range got {
		value, ok := v.Get(i)
		if !ok {
			t.Fatalf("Get(%d) returned false, but Len() is %d", i, v.Len())
		}
		got[i] = value
	}
	if !slices.Equal(got, want) {
		t.Errorf("contents = %v, want %v", got, want)
	}
}

func TestNewVectorHasZeroLength(t *testing.T) {
	v := New()
	if v.length != 0 { // If new vector is not zero length
		t.Errorf("Expected length of newVector to be 0, but got %d", v.Len())
	}
}

func TestVector_NewWithCapacity(t *testing.T) {
	tests := []struct {
		name    string
		cap     int
		wantCap int
		wantLen int
	}{
		{name: "normal capacity", cap: 4, wantCap: 4, wantLen: 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := NewWithCapacity(tc.cap)
			if v.Cap() != tc.wantCap || v.Len() != tc.wantLen {
				t.Errorf("NewWithCapacity(%d) got (cap %d, len %d), want (cap %d, len %d)", tc.cap, v.Cap(), v.Len(), tc.wantCap, tc.wantLen)
			}
		})
	}
}

func TestNewWithNegativeCapacity(t *testing.T) {
	v := NewWithCapacity(-1)
	if v != nil {
		t.Errorf("NewWithCapacity should return nil, got %v", v)
	}
}

func TestVector_Push(t *testing.T) {
	tests := []struct {
		name     string
		capacity int
		values   []int
		pushes   []int
		wantCap  int
	}{
		{name: "into zero capacity", capacity: 0, pushes: []int{1}, wantCap: 1},
		{name: "one past capacity doubles", capacity: 2, pushes: []int{1, 2, 3}, wantCap: 4},
		{name: "with spare room", capacity: 4, values: []int{1, 2, 3}, pushes: []int{4}, wantCap: 4},
		{name: "through multiple doubles", capacity: 0, pushes: []int{1, 2, 3, 4, 5}, wantCap: 8},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := newTestVector(t, tc.capacity, tc.values)
			for _, value := range tc.pushes {
				v.Push(value)
			}
			assertContents(t, v, append(tc.values, tc.pushes...))
			if got := v.Cap(); got != tc.wantCap {
				t.Errorf("Cap() = %d, want %d", got, tc.wantCap)
			}
		})
	}
}

func TestVector_Pop(t *testing.T) {
	tests := []struct {
		name      string
		values    []int
		wantValue int
		wantOk    bool
		want      []int
	}{
		{name: "only element", values: []int{1}, wantValue: 1, wantOk: true, want: []int{}},
		{name: "empty vector", values: []int{}, want: []int{}},
		{name: "multiple values", values: []int{1, 2, 3, 4}, wantOk: true, wantValue: 4, want: []int{1, 2, 3}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := newTestVector(t, len(tc.values), tc.values)
			got, ok := v.Pop()
			if got != tc.wantValue || ok != tc.wantOk {
				t.Errorf("Pop() = (%d, %t), want (%d, %t)", got, ok, tc.wantValue, tc.wantOk)
			}
			assertContents(t, v, tc.want)
			if ok && v.data[v.Len()] != 0 {
				t.Errorf("vacated slot data[%d] = %d, want 0", v.Len(), v.data[v.Len()])
			}
		})

	}
}

func TestVector_Get(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		getIndex int
		cap      int
		want     int
		wantOk   bool
	}{
		{name: "only value", values: []int{1}, getIndex: 0, cap: 1, want: 1, wantOk: true},
		{name: "negative index", values: []int{1}, getIndex: -1, cap: 1, wantOk: false},
		{name: "index equals length", values: []int{1, 2}, getIndex: 2, cap: 2, wantOk: false},
		{name: "index within cap, past length", values: []int{1, 2, 3}, getIndex: 3, cap: 4, wantOk: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := newTestVector(t, len(tc.values), tc.values)
			got, ok := v.Get(tc.getIndex)
			if got != tc.want || ok != tc.wantOk {
				t.Errorf("Get(%d) = (%d, %t), want (%d, %t)", tc.getIndex, got, ok, tc.want, tc.wantOk)
			}
		})
	}
}

func TestVector_Set(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		setIndex int
		setValue int
		want     []int
		wantOk   bool
	}{
		{name: "overwrites in place", values: []int{1, 2}, setIndex: 0, setValue: 3, want: []int{3, 2}, wantOk: true},
		{name: "past end", values: []int{1, 2}, setIndex: 3, setValue: 4, want: []int{1, 2}, wantOk: false},
		{name: "negative index", values: []int{1}, setIndex: -1, setValue: 2, want: []int{1}, wantOk: false},
		{name: "index within cap, past length", values: []int{1, 2, 3}, setIndex: 3, setValue: 4, want: []int{1, 2, 3}, wantOk: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := newTestVector(t, len(tc.values), tc.values)
			ok := v.Set(tc.setIndex, tc.setValue)
			assertContents(t, v, tc.want)
			if ok != tc.wantOk {
				t.Errorf("Set(%d, %d) = %t, want %t", tc.setIndex, tc.setValue, ok, tc.wantOk)
			}
		})
	}
}

func TestVector_Insert(t *testing.T) {
	tests := []struct {
		name        string
		values      []int
		insertValue int
		insertIndex int
		cap         int
		want        []int
		wantOk      bool
		wantCap     int
	}{
		{name: "middle, full", values: []int{1, 2}, insertValue: 3, insertIndex: 1, cap: 3, want: []int{1, 3, 2}, wantOk: true, wantCap: 3},
		{name: "start, full", values: []int{1, 2}, insertValue: 3, insertIndex: 0, cap: 3, want: []int{3, 1, 2}, wantOk: true, wantCap: 3},
		{name: "end, full", values: []int{1, 2}, insertValue: 3, insertIndex: 2, cap: 3, want: []int{1, 2, 3}, wantOk: true, wantCap: 3},
		{name: "middle, with room", values: []int{1, 2, 3}, insertValue: 4, insertIndex: 1, cap: 4, want: []int{1, 4, 2, 3}, wantOk: true, wantCap: 4},
		{name: "start, with room", values: []int{1, 2, 3}, insertValue: 4, insertIndex: 0, cap: 4, want: []int{4, 1, 2, 3}, wantOk: true, wantCap: 4},
		{name: "end, with room", values: []int{1, 2, 3}, insertValue: 4, insertIndex: 3, cap: 4, want: []int{1, 2, 3, 4}, wantOk: true, wantCap: 4},
		{name: "negative index", values: []int{1}, insertValue: 2, insertIndex: -1, cap: 1, want: []int{1}, wantOk: false, wantCap: 1},
		{name: "past end", values: []int{1}, insertValue: 2, insertIndex: 2, cap: 1, want: []int{1}, wantOk: false, wantCap: 1},
		{name: "on nil", values: nil, insertValue: 1, insertIndex: 0, cap: 0, want: []int{1}, wantOk: true, wantCap: 1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := newTestVector(t, tc.cap, tc.values)
			ok := v.Insert(tc.insertIndex, tc.insertValue)
			if ok != tc.wantOk {
				t.Errorf("Insert(%d, %d) got %t, want %t", tc.insertIndex, tc.insertValue, ok, tc.wantOk)
			}
			if v.Cap() != tc.wantCap {
				t.Errorf("Insert(%d, %d) got cap %d, want cap %d", tc.insertIndex, tc.insertValue, v.Cap(), tc.wantCap)
			}
			assertContents(t, v, tc.want)
		})
	}
}

func TestVector_Remove(t *testing.T) {
	tests := []struct {
		name        string
		values      []int
		removeIndex int
		want        []int
		wantValue   int
		wantOk      bool
	}{
		{name: "middle", values: []int{1, 2, 3}, removeIndex: 1, want: []int{1, 3}, wantValue: 2, wantOk: true},
		{name: "start", values: []int{1, 2}, removeIndex: 0, want: []int{2}, wantValue: 1, wantOk: true},
		{name: "end", values: []int{1, 2}, removeIndex: 1, want: []int{1}, wantValue: 2, wantOk: true},
		{name: "negative index", values: []int{1}, removeIndex: -1, want: []int{1}, wantOk: false},
		{name: "past end", values: []int{1}, removeIndex: 1, want: []int{1}, wantOk: false},
		{name: "on nil", values: nil, removeIndex: 0, want: nil, wantOk: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := newTestVector(t, len(tc.values), tc.values)
			got, ok := v.Remove(tc.removeIndex)
			if got != tc.wantValue || ok != tc.wantOk {
				t.Errorf("Remove(%d) got (%d, %t), want (%d, %t)", tc.removeIndex, got, ok, tc.wantValue, tc.wantOk)
			}
			assertContents(t, v, tc.want)
			if ok && v.data[v.Len()] != 0 {
				t.Errorf("vacated slot data[%d] = %d, want 0", v.Len(), v.data[v.Len()])
			}
		})
	}
}

func TestCapacityGrowth(t *testing.T) {
	vector := New()
	for i := range 1000 {
		vector.Push(i)
		vector.Pop()
	}
	if vector.Cap() != 1 {
		t.Errorf("Expected capacity of vector to be 1 after 1000 pushes and pops, but got %d", vector.Cap())
	}
}
