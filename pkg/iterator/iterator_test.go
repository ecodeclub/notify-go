package iterator

import (
	"reflect"
	"testing"
)

type MyStr string

func equal[T []E, E any](t *testing.T, it *ListIter[T, E], want T) {
	res := make([]E, 0, 10)
	for {
		got, done := it.Next()
		if done {
			break
		}
		res = append(res, got)
	}
	if !reflect.DeepEqual(res, want) {
		t.Errorf("Next() got = %v, want %v", res, want)
	}
}

func TestListIter_Nil(t *testing.T) {
	var a []int
	lr := NewListIter(a)
	want := make([]int, 0, 10)
	equal[[]int, int](t, lr, want)
}

func TestListIter_Empty(t *testing.T) {
	lr := NewListIter([]int{})
	want := make([]int, 0, 10)
	equal[[]int, int](t, lr, want)
}

func TestListIter_Base(t *testing.T) {
	lr := NewListIter[[]int, int]([]int{1, 2, 3, 4})
	want := []int{1, 2, 3, 4}
	equal[[]int, int](t, lr, want)
}

func TestListIter_TypeAlias(t *testing.T) {
	lr := NewListIter[[]MyStr, MyStr]([]MyStr{"a", "b", "c", "d"})
	want := []MyStr{"a", "b", "c", "d"}
	equal[[]MyStr, MyStr](t, lr, want)
}
