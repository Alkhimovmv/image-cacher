package linkedlist


type List[V any] struct {
	head *Item[V]
	tail *Item[V]
	len int
}

func New[V any]() *List[V] {
	return &List[V]{}
}

func (l *List[V]) Len() int {
	return l.len
}

func (l *List[V]) PushFront(value V) *Item[V] {
	item := &Item[V]{
		Value: value,
	}

	if l.len == 0 {
		l.head = item
		l.tail = item
		l.len++
		return item
	}

	item.next = l.head
	l.head.prev = item
	l.head = item
	l.len++

	return item
}

func (l *List[V]) PushBack(value V) *Item[V] {
	item := &Item[V]{
		Value: value,
	}

	if l.len == 0 {
		l.head = item
		l.tail = item
		l.len++
		return item
	}

	item.prev = l.tail
	l.tail.next = item
	l.tail = item
	l.len++

	return item
}