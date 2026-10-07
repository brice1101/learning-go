package hashmap

import (
	"strconv"
	"testing"
)

func newTestHashMap(putKeys []string, putValues []int) *HashMap {
	h := New()
	for i := range len(putKeys) {
		h.Put(tc.putKeys[i], putValues[i])
	}
	return h
}

func TestKnownHashValues(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want uint64
	}{
		{name: "Empty string", key: "", want: 14695981039346656037},
		{name: "a", key: "a", want: 0xaf63dc4c8601ec8c},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if hash := fnv1a(tc.key); hash != tc.want {
				t.Errorf("FNV-1a on \"%s\" string should return %d, got %d", tc.key, tc.want, hash)
			}
		})
	}
}

func TestLen(t *testing.T) {
	tests := []struct {
		name      string
		putKeys   []string
		putValues []int
		want      int
	}{
		{name: "Empty", want: 0},
		{name: "One put", putKeys: []string{"apple"}, putValues: []int{1}, want: 1},
		{name: "Multiple puts", putKeys: []string{"a", "b", "c"}, putValues: []int{1, 2, 3}, want: 3},
		{name: "Put same value twice", putKeys: []string{"a", "a"}, putValues: []int{1, 2}, want: 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHashMap(tc.putKeys, tc.putValues)
			if length := h.Len(); length != tc.want {
				t.Errorf("Expected length %d, got %d", tc.want, length)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name      string
		putKeys   []string
		putValues []int
		getKey    string
		want      int
		wantOk    bool
	}{
		{name: "Only Value", putKeys: []string{"a"}, putValues: []int{1}, getKey: "a", want: 1, wantOk: true},
		{name: "Value not present", putKeys: []string{"a", "b"}, putValues: []int{1, 2}, getKey: "c", wantOk: false},
		{name: "On empty", getKey: "a", wantOk: false},
		{name: "Put same value twice", putKeys: []string{"a", "a"}, putValues: []int{1, 2}, getKey: "a", want: 2, wantOk: true},
		{name: "Empty string key", putKeys: []string{""}, putValues: []int{1}, getKey: "", want: 1, wantOk: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHashMap(tc.putKeys, tc.putValues)
			value, ok := h.Get(tc.getKey)
			if value != tc.want|ok != tc.wantOk {
				t.Errorf("Get() expected (%d, %t), got (%d, %t)", tc.want, tc.wantOk, value, ok)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		putKeys   []string
		putValues []int
		deleteKey string
		wantOk    bool
		wantLen   int
	}{
		{name: "Present key", putKeys: []string{"a", "b"}, putValues: []int{1, 2}, deleteKey: "a", wantOk: true, wantLen: 1},
		{name: "Absent key", putKeys: []string{"a", "b"}, putValues: []int{1, 2}, deleteKey: "c", wantOk: false, wantLen: 2},
		{name: "Empty hash", deleteKey: "a", wantOk: false, wantLen: 0},
		{name: "Duplicate values", putKeys: []string{"a", "a"}, putValues: []int{1, 2}, deleteKey: "a", wantOk: true, wantLen: 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHashMap(tc.putKeys, tc.putValues)
			ok := h.Delete(tc.deleteKey)
			if ok != tc.wantOk|h.Len() != tc.wantLen {
				t.Errorf("Delete() expects %t and length %d, got %t and length %d", tc.wantOk, tc.wantLen, ok, h.Len())
			}
			_, getOk := h.Get(tc.deleteKey)
			if getOk != tc.wantOk {
				t.Errorf("Get() on deleted value expects false, got %t", getOk)
			}
		})
	}
}

func TestResizes(t *testing.T) {
	var numElements float32 = 1000
	var maxLoadFactor float32 = 0.75
	h := New()
	for i := range int(numElements) {
		h.Put(strconv.Itoa(i), i)
	}
	for i := range int(numElements) {
		value, ok := h.Get(strconv.Itoa(i))
		if value != i|ok != true {
			t.Errorf("Get(%d) expected %d, got %d", i, i, value)
		}
	}
	if load := float32(len(h.buckets)) / numElements; load > maxLoadFactor {
		t.Errorf("Load factor should be below %d, got %d", maxLoadFactor, load)
	}
}

/*
func TestCollisions(t *testing.T) {
	tests := []struct {
		name      string
		deleteKey string
	}{
		{},
	}

	h := HashMap{buckets: make([]*entry, 8), hashFunction: func(string) uint64 { return 0 }, length: 0}
	for i := range 3 {
		h.Push(strconv.Itoa(i), i)
	}
	for i := range 3 {
		if val := h.Get(strconv.Itoa(i)); val != i {
			t.Errorf("Get() expected %d, got %d", i, val)
		}
	}
}
*/
