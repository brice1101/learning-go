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
	l.head = &n
}

func (l *Linkedlist) PushBack(value int) {
	n := Node{value, nil}
	l.tail.next = &n
}

func (l *Linkedlist) PopFront() (int, bool) {
	if l.head == nil {
		return 0, false
	}
	value := l.head.value
	l.head = l.head.next
	return value, true
}

func (l *Linkedlist) Find(value int) (int, bool) {
	n := l.head
	i := 0
	for n.next != nil {
		if n.value == value {
			return i, true
		}
		n = n.next
		i++
	}
	if n.value == value {
		return i, true
	}
	return 0, false
}

func (l *Linkedlist) Remove(value int) bool {
	// case first item
	n := l.head
	if n.value == value {
		l.head = l.head.next
		return true
	}

	// case middle
	prev := *n
	for n.next != nil {
		if n.value == value {
			prev.next = n.next
			return true
		}
		prev = *n
		n = n.next

	}

	// case end
	if n.value == value {
		prev.next = nil
		return true
	}
	return false
}

func (l *Linkedlist) Len() int {
	n := l.head
	if n == nil {
		return 0
	}
	counter := 1
	for n.next != nil {
		n = n.next
		counter++
	}
	return counter
}

func (l *Linkedlist) Reverse() {
	n := l.head
	flipPointer(n, n.next)
}

func flipPointer(nodeA *Node, nodeB *Node) {
	if nodeB.next != nil {
		flipPointer(nodeB, nodeB.next)
	}
	nodeB.next = nodeA
}
