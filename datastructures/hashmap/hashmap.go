package hashmap

type HashMap struct {
	buckets      []*entry
	hashFunction func(string) uint64
	length       int
}

type entry struct {
	key   string
	value int
	next  *entry
}

func fnv1a(key string) uint64 {
	var hash uint64 = 14695981039346656037
	for _, c := range []byte(key) {
		hash = hash ^ uint64(c)
		hash = hash * 1099511628211
	}
	return hash
}

func New() *HashMap {
	return &HashMap{buckets: make([]*entry, 8), hashFunction: fnv1a, length: 0}
}
