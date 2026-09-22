package vector

type Vector struct {
	data     []int
	capacity int
	length   int
}

func (v *Vector) Push(value int) {
	if v.capacity == v.length {
		if v.capacity == 0 {
			v.capacity = 1
		} else {
			v.capacity *= 2
		}
		empty := make([]int, v.capacity)
		copy(empty, v.data)
		v.data = empty
	}
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

func (v *Vector) Get(index int) int {
	return v.data[index]
}
