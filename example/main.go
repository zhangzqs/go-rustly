package main

import (
	"fmt"
	"time"

	. "github.com/zhangzqs/go-rustly"
)

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
	c.infos.Extend(NewMap[Pair[Into[Stringer], Into[Stringer]], Pair[Stringer, Stringer]](
		it.IntoIter(),
		func(p Pair[Into[Stringer], Into[Stringer]]) Pair[Stringer, Stringer] {
			return Pair[Stringer, Stringer]{first: p.first.Into(), second: p.second.Into()}
		},
	))
}

type SharedState struct {
	completed bool
	waker     Option[Waker]
}

type TimerFuture struct {
	sharedState *Mutex[SharedState]
}

func (f *TimerFuture) Poll(cx FutureContext) PollResult[Option[struct{}]] {
	state := f.sharedState.Lock()
	defer f.sharedState.Drop()
	if state.completed {
		return PollReady[Option[struct{}]](Option[struct{}]{})
	}
	state.waker = Some[Waker](cx.Waker)
	return PollPending[Option[struct{}]]()
}

func TimerFutureNew(duration time.Duration) *TimerFuture {
	sharedState := MutexNew[SharedState](&SharedState{
		completed: false,
		waker:     None[Waker](),
	})
	go func() {
		time.Sleep(duration)
		state := sharedState.Lock()
		defer sharedState.Drop()
		// 通知executor当前future已经完成，executor会调用waker的wake方法
		state.completed = true
		if state.waker.IsSome() {
			state.waker.Unwrap().Wake()
		}
	}()
	return &TimerFuture{sharedState: &sharedState}
}

type JoinFuture struct {
	futureA Option[*TimerFuture]
	futureB Option[*TimerFuture]
}

func (f *JoinFuture) Poll(cx FutureContext) PollResult[struct{}] {
	if f.futureA.IsSome() {
		if f.futureA.Unwrap().Poll(cx).IsReady() {
			f.futureA.Take()
		}
	}
	if f.futureB.IsSome() {
		if f.futureB.Unwrap().Poll(cx).IsReady() {
			f.futureB.Take()
		}
	}
	if f.futureA.IsNone() && f.futureB.IsNone() {
		// 两个future都完成了
		return PollReady[struct{}](struct{}{})
	} else {
		return PollPending[struct{}]()
	}
}

func main() {
	// executor, spawner := NewExecutorAndSpawner()
	// join := &JoinFuture{
	// 	futureA: Some(TimerFutureNew(time.Second * 2)),
	// 	futureB: Some(TimerFutureNew(time.Second * 3)),
	// }
	// spawner.Spawn(join)
	// executor.Run()

	ch1 := make(chan int)

	for v := range ch1 {
		fmt.Println(v)
	}

	for {
		select {
		case v := <-ch1:
			fmt.Println(v)
		}
	}
}
