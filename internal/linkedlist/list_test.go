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

func TestList_RemoveMiddleItem(t *testing.T) {
	list := New[int]()

	first := list.PushBack(1)
	second := list.PushBack(2)
	third := list.PushBack(3)

	list.Remove(second)

	if list.Len() != 2 {
		t.Fatalf("ожидали длину списка 2, получили %d", list.Len())
	}

	if list.head != first {
		t.Fatal("ожидали, что первый элемент останется головой списка")
	}

	if list.tail != third {
		t.Fatal("ожидали, что третий элемент останется хвостом списка")
	}

	if first.next != third {
		t.Fatal("ожидали, что первый элемент будет ссылаться на третий")
	}

	if third.prev != first {
		t.Fatal("ожидали, что третий элемент будет ссылаться на первый")
	}
}

func TestList_RemoveFirstItem(t *testing.T) {
	list := New[int]()

	first := list.PushBack(1)
	second := list.PushBack(2)
	third := list.PushBack(3)

	list.Remove(first)

	if list.Len() != 2 {
		t.Fatalf("ожидали длину списка 2, получили %d", list.Len())
	}

	if list.head != second {
		t.Fatal("ожидали, что второй элемент станет головой списка")
	}

	if list.tail != third {
		t.Fatal("ожидали, что третий элемент останется хвостом списка")
	}

	if second.prev != nil {
		t.Fatal("ожидали, что у новой головы не будет предыдущего элемента")
	}

	if second.next != third {
		t.Fatal("ожидали, что второй элемент будет ссылаться на третий")
	}
}

func TestList_RemoveLastItem(t *testing.T) {
	list := New[int]()

	first := list.PushBack(1)
	second := list.PushBack(2)
	third := list.PushBack(3)

	list.Remove(third)

	if list.Len() != 2 {
		t.Fatalf("ожидали длину списка 2, получили %d", list.Len())
	}

	if list.head != first {
		t.Fatal("ожидали, что первый элемент останется головой списка")
	}

	if list.tail != second {
		t.Fatal("ожидали, что второй элемент станет хвостом списка")
	}

	if second.next != nil {
		t.Fatal("ожидали, что у нового хвоста не будет следующего элемента")
	}

	if second.prev != first {
		t.Fatal("ожидали, что второй элемент будет ссылаться назад на первый")
	}
}

func TestList_RemoveOnceItem(t *testing.T) {
	list := New[int]()

	item := list.PushBack(1)

	list.Remove(item)

	if list.Len() != 0 {
		t.Fatalf("ожидали длину списка 0, получили %d", list.Len())
	}

	if list.head != nil {
		t.Fatal("ожидали, что голова списка станет nil")
	}

	if list.tail != nil {
		t.Fatal("ожидали, что хвост списка станет nil")
	}
}

func TestList_RemoveClearsItemLinks(t *testing.T) {
	list := New[int]()

	list.PushBack(1)
	item := list.PushBack(2)
	list.PushBack(3)

	list.Remove(item)

	if item.prev != nil {
		t.Fatal("ожидали, что ссылка удалённого элемента на предыдущий элемент станет nil")
	}

	if item.next != nil {
		t.Fatal("ожидали, что ссылка удалённого элемента на следующий элемент станет nil")
	}
}
