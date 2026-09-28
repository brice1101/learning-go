package linkedlist

type Node struct {
	value int
	next  *Node
}
type Linkedlist struct {
	head *Node
	tail *Node
}

func (l *Linkedlist) PushFront(value int) {
	n := Node{value, l.head}
	if l.head == nil {
		l.tail = &n
	}
	l.head = &n
}

func (l *Linkedlist) PushBack(value int) {
	n := Node{value, nil}
	if l.head == nil {
		l.head = &n
	} else {
		l.tail.next = &n
	}
	l.tail = &n
}

func (l *Linkedlist) PopFront() (int, bool) {
	if l.head == nil {
		return 0, false
	}
	value := l.head.value
	l.head = l.head.next
	if l.head == nil {
		l.tail = nil
	}
	return value, true
}

func (l *Linkedlist) Find(value int) (int, bool) {
	if l.head == nil {
		return 0, false
	}
	n := l.head
	i := 0
	for n != nil {
		if n.value == value {
			return i, true
		}
		n = n.next
		i++
	}
	return 0, false
}

func (l *Linkedlist) Remove(value int) bool {
	if l.head == nil {
		return false
	}
	// case first item
	n := l.head
	if n.value == value {
		l.PopFront()
		return true
	}

	// case middle
	var prev *Node
	for {
		prev = n
		n = n.next
		if n == nil {
			return false
		} else if n.value == value {
			if n.next == nil {
				l.tail = prev
			}
			prev.next = n.next
			return true
		}
	}
}

func (l *Linkedlist) Len() int {
	n := l.head
	if n == nil {
		return 0
	}
	counter := 0
	for n != nil {
		n = n.next
		counter++
	}
	return counter
}

func (l *Linkedlist) Reverse() {
	curr := l.head
	var prev *Node = nil
	var next *Node = nil
	for curr != nil {
		next = curr.next
		curr.next = prev
		prev = curr
		curr = next
	}
	l.tail = l.head
	l.head = prev

}
