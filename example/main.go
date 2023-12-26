package main

import . "github.com/zhangzqs/go-rustly"

type Stringer interface {
	String() string
}

type StringWrapper string

func (s StringWrapper) String() string {
	return string(s)
}

type Pair[T any, U any] struct {
	first  T
	second U
}

type Collector struct {
	infos Vec[Pair[Stringer, Stringer]]
}

func Extend[T IntoIterator[Pair[Into[Stringer], Into[Stringer]]]](c *Collector, it T) {
	c.infos.Extend(MapNew[Pair[Into[Stringer], Into[Stringer]], Pair[Stringer, Stringer]](
		it.IntoIter(),
		func(p Pair[Into[Stringer], Into[Stringer]]) Pair[Stringer, Stringer] {
			return Pair[Stringer, Stringer]{first: p.first.Into(), second: p.second.Into()}
		},
	))
}

func main() {
}
