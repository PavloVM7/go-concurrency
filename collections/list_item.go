package collections

type listItem[T any] struct {
	prev  *listItem[T]
	next  *listItem[T]
	value T
}

func (li *listItem[T]) insert(item *listItem[T]) {
	item.prev = li.prev
	item.next = li
	li.prev = item
}
func (li *listItem[T]) append(item *listItem[T]) {
	item.prev = li
	item.next = li.next
	li.next = item
}
func (li *listItem[T]) removeYourself() {
	if li.prev != nil {
		li.prev.next = li.next
	}
	if li.next != nil {
		li.next.prev = li.prev
	}
}

func swapListItems[T any](item1, item2 *listItem[T]) {
	item1.value, item2.value = item2.value, item1.value
}
