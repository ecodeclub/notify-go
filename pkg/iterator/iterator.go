// Copyright 2021 ecodeclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
