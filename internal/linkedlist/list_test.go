package linkedlist

import "testing"

func TestList_PushFrontAddsItemsToFront(t *testing.T) {
	list := New[int]()

	first := list.PushFront(1)
	second := list.PushFront(2)

	if list.Len() != 2 {
		t.Fatalf("ожидали длину списка 2, получили %d", list.Len())
	}

	if list.head != second {
		t.Fatal("ожидали, что второй добавленный элемент станет головой списка")
	}

	if list.tail != first {
		t.Fatal("ожидали, что первый добавленный элемент станет хвостом списка")
	}

	if list.head.Value != 2 {
		t.Fatalf("ожидали значение головы 2, получили %d", list.head.Value)
	}

	if list.tail.Value != 1 {
		t.Fatalf("ожидали значение хвоста 1, получили %d", list.tail.Value)
	}
}

func TestList_PushBackAddsItemsToBack(t *testing.T) {
	list := New[int]()

	first := list.PushBack(1)
	second := list.PushBack(2)

	if list.Len() != 2 {
		t.Fatalf("ожидали длину списка 2, получили %d", list.Len())
	}

	if list.head != first {
		t.Fatal("ожидали, что первый добавленный элемент станет головой списка")
	}

	if list.tail != second {
		t.Fatal("ожидали, что второй добавленный элемент станет хвостом списка")
	}

	if list.head.Value != 1 {
		t.Fatalf("ожидали значение головы 1, получили %d", list.head.Value)
	}

	if list.tail.Value != 2 {
		t.Fatalf("ожидали значение хвоста 2, получили %d", list.tail.Value)
	}
}
