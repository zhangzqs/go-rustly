package gorustly

type PollResult[T any] struct {
	// Ready is true if the future is ready
	Ready bool
	// Value is the value of the future if it is ready
	Value Option[T]
}

func PollReady[T any](value T) PollResult[T] {
	return PollResult[T]{Ready: true, Value: Some[T](value)}
}

func PollPending[T any]() PollResult[T] {
	return PollResult[T]{Ready: false, Value: None[T]()}
}

func (r PollResult[T]) IsReady() bool {
	return r.Ready
}

func (r PollResult[T]) GetValue() Option[T] {
	return r.Value
}

func (r PollResult[T]) IsPending() bool {
	return !r.Ready
}

type Waker interface {
	Wake()
}

type FutureContext struct {
	Waker Waker
}

type Future[T any] interface {
	Poll(cx FutureContext) PollResult[T]
}

type Task struct {
	future     Mutex[Option[Future[struct{}]]]
	taskSender chan<- *Task
}

func (t *Task) Wake() {
	t.taskSender <- t
}

type Executor struct {
	readyQueue <-chan *Task
}

func (e *Executor) Run() {
	for {
		task := <-e.readyQueue
		future_slot := task.future.Lock()

		future := future_slot.Take()
		if future.IsSome() {
			result := future.Unwrap().Poll(FutureContext{Waker: task})
			if result.IsPending() {
				// 还没执行完，放过队列里，等待下次poll
				*future_slot = Some[Future[struct{}]](future.Unwrap())
			}
		}
		task.future.Drop()
	}
}

type Spawner struct {
	taskSender chan<- *Task
}

func (s *Spawner) Spawn(future Future[struct{}]) {
	a := Some[Future[struct{}]](future)
	task := &Task{
		future:     MutexNew[Option[Future[struct{}]]](&a),
		taskSender: s.taskSender,
	}
	s.taskSender <- task
}

func NewExecutorAndSpawner() (Executor, Spawner) {
	taskChannel := make(chan *Task, 1000)
	return Executor{readyQueue: taskChannel}, Spawner{taskSender: taskChannel}
}
