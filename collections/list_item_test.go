package collections

import (
	"reflect"
	"testing"
)

func Test_swapListItems(t *testing.T) {

	type args[T any] struct {
		list   *ConcurrentLinkedList[T]
		index1 int
		index2 int
		want   []int
	}
	type testCase[T any] struct {
		name string
		args args[T]
	}
	tests := []testCase[int]{
		{"first<->last", args[int]{NewConcurrentLinkedListItems[int](1, 2), 0, 1, []int{2, 1}}},
		{"first<-1->last", args[int]{NewConcurrentLinkedListItems[int](1, 2, 3), 0, 2, []int{3, 2, 1}}},
		{"2<->3", args[int]{NewConcurrentLinkedListItems[int](1, 2, 3, 4), 1, 2, []int{1, 3, 2, 4}}},
		{"first<->2", args[int]{NewConcurrentLinkedListItems[int](1, 2, 3), 0, 1, []int{2, 1, 3}}},
		{"2<->last", args[int]{NewConcurrentLinkedListItems[int](1, 2, 3), 1, 2, []int{1, 3, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.list.mu.RLock()
			item1, _ := tt.args.list.getByIndex(tt.args.index1)
			item2, _ := tt.args.list.getByIndex(tt.args.index2)
			tt.args.list.mu.RUnlock()
			swapListItems(item1, item2)
			actual := tt.args.list.ToArray()
			if !reflect.DeepEqual(actual, tt.args.want) {
				t.Errorf("swapListItems() got: %v, want: %v", actual, tt.args.want)
			}
		})
	}
}
