package linkedlist

type Item[V any] struct {
	Value V

	prev *Item[V]
	next *Item[V]
}
