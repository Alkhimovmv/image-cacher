package linkedlist

type List[V any] struct {
	head *Item[V]
	tail *Item[V]
	len  int
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

func (l *List[V]) Remove(item *Item[V]) {
	if item.prev != nil {
		item.prev.next = item.next
	} else {
		l.head = item.next
	}

	if item.next != nil {
		item.next.prev = item.prev
	} else {
		l.tail = item.prev
	}

	item.prev = nil
	item.next = nil

	l.len--
}

func (l *List[V]) MoveToFront(item *Item[V]) {
	if item == nil || l.head == item {
		return
	}

	if item.prev != nil {
		item.prev.next = item.next
	}

	if item.next != nil {
		item.next.prev = item.prev
	} else {
		l.tail = item.prev
	}

	item.prev = nil
	item.next = l.head
	l.head.prev = item
	l.head = item
}
