package gorustly

type Iterator[T any] interface {
	Next() Option[T]
}

type MapIterator[T any, U any] struct {
	iter Iterator[T]
	fn   func(T) U
}

func (m *MapIterator[T, U]) Next() Option[U] {
	if m.iter.Next().IsNone() {
		return None[U]()
	}
	return Some[U](m.fn(m.iter.Next().Unwrap()))
}

func (m *MapIterator[T, U]) IntoIter() Iterator[U] {
	return m
}

func NewMap[T any, U any](iter Iterator[T], fn func(T) U) *MapIterator[T, U] {
	return &MapIterator[T, U]{iter: iter, fn: fn}
}
