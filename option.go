package gorustly

type Option[T any] struct {
	E *T
}

func Some[T any](e T) Option[T] {
	return Option[T]{&e}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func (o Option[T]) IsSome() bool {
	return o.E != nil
}

func (o Option[T]) IsNone() bool {
	return o.E == nil
}

func (o Option[T]) Unwrap() T {
	return *o.E
}

func (o Option[T]) UnwrapOr(def T) T {
	if o.E == nil {
		return def
	}
	return *o.E
}

func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.E == nil {
		return fn()
	}
	return *o.E
}
