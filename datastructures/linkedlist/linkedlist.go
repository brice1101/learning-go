package linkedlist

type Node struct {
	value    int
	nextNode *Node
}

type LinkedList struct {
	headPointer *Node
	tailPointer *Node
}

func (l *LinkedList) PushFront(v int) {

}

func (l *LinkedList) PushBack(v int) {

}

func (l *LinkedList) PopFront() int {
	return 0
}

func (l *LinkedList) Find(v int) int {
	return 0
}

func (l *LinkedList) Remove(v int) bool {
	return false
}

func (l *LinkedList) Len() int {
	return 0
}

func (l *LinkedList) Reverse() {

}
