package vector

type Vector struct {
	data   []int
	length int
}

func New() *Vector {
	return &Vector{}
}

func NewWithCapacity(capacity int) *Vector {
	if capacity < 0 {
		return nil
	}
	return &Vector{
		data: make([]int, capacity),
	}
}

func (v *Vector) Len() int {
	return v.length
}

func (v *Vector) Cap() int {
	return len(v.data)
}

func (v *Vector) Get(i int) (int, bool) {
	if i >= v.length || i < 0 {
		return 0, false
	}
	return v.data[i], true
}

func (v *Vector) Set(i int, value int) bool {
	if i >= v.length || i < 0 {
		return false
	}
	v.data[i] = value
	return true
}

func (v *Vector) Push(value int) {
	v.growIfFull(v.length)
	v.data[v.length] = value
	v.length++
}

func (v *Vector) Pop() (int, bool) {
	if v.length == 0 {
		return 0, false
	}
	value := v.data[v.length-1]
	v.data[v.length-1] = 0
	v.length--
	return value, true
}

func (v *Vector) Insert(i int, value int) bool {
	if i > v.length || i < 0 {
		return false
	}
	v.growIfFull(i)
	v.data[i] = value
	v.length++
	return true
}

func (v *Vector) Remove(i int) (int, bool) {
	if i >= v.length || i < 0 {
		return 0, false
	}
	value := v.data[i]
	for j := i + 1; j < v.length; j++ {
		v.data[j-1] = v.data[j]
	}
	v.data[v.length-1] = 0
	v.length--
	return value, true
}

func (v *Vector) growIfFull(gap int) {
	if capacity := len(v.data); v.length == capacity {
		if capacity == 0 {
			capacity = 1
		} else {
			capacity *= 2
		}
		empty := make([]int, capacity)
		copy(empty[:gap], v.data[:gap])
		copy(empty[gap+1:], v.data[gap:])
		v.data = empty
	} else {
		copy(v.data[gap+1:], v.data[gap:])
	}
}
