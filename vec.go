package gorustly

type Vec[T any] struct {
	slice []T
}

func VecNew[T any]() Vec[T] {
	return Vec[T]{}
}

func VecNewWith[T any](slice []T) Vec[T] {
	return Vec[T]{slice: slice}
}

func (v *Vec[T]) Push(e T) {
	v.slice = append(v.slice, e)
}

func (v *Vec[T]) Pop() Option[T] {
	if len(v.slice) == 0 {
		return None[T]()
	}
	e := v.slice[len(v.slice)-1]
	v.slice = v.slice[:len(v.slice)-1]
	return Some[T](e)
}

func (v *Vec[T]) Len() int {
	return len(v.slice)
}

func (v *Vec[T]) IsEmpty() bool {
	return len(v.slice) == 0
}

func (v *Vec[T]) Get(i int) Option[T] {
	if i < 0 || i >= len(v.slice) {
		return None[T]()
	}
	return Some[T](v.slice[i])
}

func (v *Vec[T]) Set(i int, e T) {
	if i < 0 || i >= len(v.slice) {
		return
	}
	v.slice[i] = e
}

func (v *Vec[T]) Swap(i int, j int) {
	if i < 0 || i >= len(v.slice) || j < 0 || j >= len(v.slice) {
		return
	}
	v.slice[i], v.slice[j] = v.slice[j], v.slice[i]
}

func (v *Vec[T]) Insert(i int, e T) {
	if i < 0 || i > len(v.slice) {
		return
	}
	v.slice = append(v.slice, e)
	copy(v.slice[i+1:], v.slice[i:])
	v.slice[i] = e
}

func (v *Vec[T]) Remove(i int) Option[T] {
	if i < 0 || i >= len(v.slice) {
		return None[T]()
	}
	e := v.slice[i]
	v.slice = append(v.slice[:i], v.slice[i+1:]...)
	return Some[T](e)
}

func (v *Vec[T]) Clear() {
	v.slice = v.slice[:0]
}

func (v *Vec[T]) Extend(it IntoIterator[T]) {
	iter := it.IntoIter()
	for {
		if iter.Next().IsNone() {
			break
		}
		v.slice = append(v.slice, iter.Next().Unwrap())
	}
}

type VecIterator[T any] struct {
	v    *Vec[T]
	next int
}

func (vi *VecIterator[T]) Next() Option[T] {
	if vi.next >= len(vi.v.slice) {
		return None[T]()
	}
	e := vi.v.slice[vi.next]
	vi.next++
	return Some[T](e)
}

func (v *Vec[T]) IntoIter() Iterator[T] {
	return &VecIterator[T]{v: v}
}
