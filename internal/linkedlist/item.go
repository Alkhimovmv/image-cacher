package linkedlist

type Item[V any] struct {
	Value V

	prev *Item[V]
	next *Item[V]
}

func (i *Item[V]) Next() *Item[V] {
	return i.next
}

func (i *Item[V]) Prev() *Item[V] {
	return i.prev
}

func (l *List[V]) Front() *Item[V] {
	return l.head
}

func (l *List[V]) Back() *Item[V] {
	return l.tail
}
