package gorustly

import "reflect"

type Result[T any, E any] struct {
	value T
	err   E
}

func Ok[T any, E any](value T) Result[T, E] {
	return Result[T, E]{
		value: value,
	}
}

func Err[T any, E any](err E) Result[T, E] {
	return Result[T, E]{
		err: err,
	}
}

func (r *Result[T, E]) IsOk() bool {
	return reflect.ValueOf(r.value).IsNil()
}

func (r *Result[T, E]) IsErr() bool {
	return reflect.ValueOf(r.err).IsNil()
}

func (r *Result[T, E]) Error() E {
	return r.err
}

func (o *Result[T, E]) Unwrap() T {
	if o.IsErr() {
		panic("can't unwrap err value")
	}
	return o.value
}

func (o *Result[T, E]) UnwrapOr(def T) T {
	if o.IsErr() {
		return def
	}
	return o.value
}

func (o *Result[T, E]) UnwrapOrElse(fn func() T) T {
	if o.IsErr() {
		return fn()
	}
	return o.value
}

func (o *Result[T, E]) UnwrapErr() E {
	if o.IsOk() {
		panic("can't unwrap err value")
	}
	return o.err
}

func (o *Result[T, E]) UnwrapErrOr(def E) E {
	if o.IsOk() {
		return def
	}
	return o.err
}

func (o *Result[T, E]) UnwrapErrOrElse(fn func() E) E {
	if o.IsOk() {
		return fn()
	}
	return o.err
}

func (o *Result[T, E]) And(res Result[any, E]) Result[any, E] {
	if o.IsOk() {
		return res
	}

	return Result[any, E]{
		err: o.err,
	}
}
