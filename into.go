package gorustly

type Into[T any] interface {
	Into() T
}

type IntoIterator[T any] interface {
	IntoIter() Iterator[T]
}
