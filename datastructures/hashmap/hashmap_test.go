package hashmap

import "testing"

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
			if hash := fnv_1a(tc.key); hash != tc.want {
				t.Errorf("FNV-1a on \"%s\" string should return %d, got %d", tc.key, tc.want, hash)
			}
		})
	}
}
