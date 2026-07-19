package lru

import "testing"

func TestLRU_SetAddsNewItem(t *testing.T) {
	cache := New[string](100)

	cache.Set("hash", "cat-file", 10)

	value, ok := cache.Get("hash")

	if !ok {
		t.Fatal("ожидали, что ключ hash будет найден")
	}

	if value != "cat-file" {
		t.Fatalf("ожидали значение cat-file, получили %s", value)
	}
}

func TestLRU_GetReturnsFalseForMissingKey(t *testing.T) {
	cache := New[string](100)

	value, ok := cache.Get("missing")

	if ok {
		t.Fatal("ожидали, что для отсутствующего ключа ok будет false")
	}

	if value != "" {
		t.Fatalf("ожидали пустую строку для отсутствующего ключа, получили %q", value)
	}
}

func TestLRU_GetMovesItemToFront(t *testing.T) {
	cache := New[string](100)

	cache.Set("first", "file-1", 10)
	cache.Set("second", "file-2", 10)
	cache.Set("third", "file-3", 10)

	_, ok := cache.Get("first")
	if !ok {
		t.Fatal("ожидали, что ключ first будет найден")
	}

	items := cache.Flush()

	if len(items) != 3 {
		t.Fatalf("ожидали 3 элемента, получили %d", len(items))
	}

	if items[0].Key != "first" {
		t.Fatalf("ожидали, что ключ first станет первым после Get, получили %s", items[0].Key)
	}
}

func TestLRU_SetUpdatesExistingItem(t *testing.T) {
	cache := New[string](100)

	cache.Set("cat", "old-file", 10)
	cache.Set("cat", "new-file", 15)

	value, ok := cache.Get("cat")
	if !ok {
		t.Fatal("ожидали, что ключ cat будет найден")
	}

	if value != "new-file" {
		t.Fatalf("ожидали значение new-file, получили %s", value)
	}

	items := cache.Flush()
	if len(items) != 1 {
		t.Fatalf("ожидали 1 элемент после обновления ключа, получили %d", len(items))
	}

	if items[0].Size != 15 {
		t.Fatalf("ожидали размер 15, получили %d", items[0].Size)
	}
}

func TestLRU_SetEvictsLeastRecentlyUsedItems(t *testing.T) {
	cache := New[string](25)

	cache.Set("first", "file-1", 10)
	cache.Set("second", "file-2", 10)
	evicted := cache.Set("third", "file-3", 10)

	if len(evicted) != 1 {
		t.Fatalf("ожидали 1 вытесненный элемент, получили %d", len(evicted))
	}

	if evicted[0] != "file-1" {
		t.Fatalf("ожидали, что вытеснится file-1, получили %s", evicted[0])
	}

	if _, ok := cache.Get("first"); ok {
		t.Fatal("ожидали, что ключ first будет удалён из кеша")
	}

	if _, ok := cache.Get("second"); !ok {
		t.Fatal("ожидали, что ключ second останется в кеше")
	}

	if _, ok := cache.Get("third"); !ok {
		t.Fatal("ожидали, что ключ third останется в кеше")
	}
}

func TestLRU_GetAffectsEvictionOrder(t *testing.T) {
	cache := New[string](25)

	cache.Set("first", "file-1", 10)
	cache.Set("second", "file-2", 10)

	_, ok := cache.Get("first")
	if !ok {
		t.Fatal("ожидали, что ключ first будет найден")
	}

	evicted := cache.Set("third", "file-3", 10)

	if len(evicted) != 1 {
		t.Fatalf("ожидали 1 вытесненный элемент, получили %d", len(evicted))
	}

	if evicted[0] != "file-2" {
		t.Fatalf("ожидали, что вытеснится file-2, получили %s", evicted[0])
	}

	if _, ok := cache.Get("first"); !ok {
		t.Fatal("ожидали, что ключ first останется в кеше")
	}

	if _, ok := cache.Get("second"); ok {
		t.Fatal("ожидали, что ключ second будет удалён из кеша")
	}

	if _, ok := cache.Get("third"); !ok {
		t.Fatal("ожидали, что ключ third останется в кеше")
	}
}

func TestLRU_RestoreDataLoadsItemsInOrder(t *testing.T) {
	cache := New[string](100)

	restoreItems := []CacheItem[string]{
		{Key: "first", Value: "file-1", Size: 10},
		{Key: "second", Value: "file-2", Size: 20},
		{Key: "third", Value: "file-3", Size: 30},
	}

	cache.RestoreData(restoreItems)

	flushed := cache.Flush()

	if len(flushed) != 3 {
		t.Fatalf("ожидали 3 элемента после восстановления, получили %d", len(flushed))
	}

	if flushed[0].Key != "first" {
		t.Fatalf("ожидали first на позиции 0, получили %s", flushed[0].Key)
	}

	if flushed[1].Key != "second" {
		t.Fatalf("ожидали second на позиции 1, получили %s", flushed[1].Key)
	}

	if flushed[2].Key != "third" {
		t.Fatalf("ожидали third на позиции 2, получили %s", flushed[2].Key)
	}

	value, ok := cache.Get("second")
	if !ok {
		t.Fatal("ожидали, что ключ second будет найден после восстановления")
	}

	if value != "file-2" {
		t.Fatalf("ожидали значение file-2, получили %s", value)
	}
}
