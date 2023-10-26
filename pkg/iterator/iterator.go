package iterator

type ListIter[T []E, E any] struct {
	val    T
	length int
	index  int
}

var _ Iterable[any] = &ListIter[[]any, any]{}

func (it *ListIter[T, E]) Next() (E, bool) {
	var e E
	if it.index >= it.length {
		return e, true
	}
	res := it.val[it.index]
	it.index++
	return res, false
}

func NewListIter[T []E, E any](array T) *ListIter[T, E] {
	if array == nil {
		return &ListIter[T, E]{}
	}
	return &ListIter[T, E]{
		val:    array,
		length: len(array),
		index:  0,
	}
}
