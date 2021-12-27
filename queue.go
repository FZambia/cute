// Package queue provides a generic unbounded queue for Go programming language.
package queue

import (
	"sync"
)

// Queue is a generic unbounded queue of any T.
// It can optionally maintain the total cost of currently queued elements.
// All queue methods are goroutine-safe. When using cost be aware that integer overflow
// is not handled here, supposing that the queue will be closed well below max int size.
type Queue[T any] struct {
	mu      sync.RWMutex
	cond    *sync.Cond
	nodes   []queuedItem[T]
	head    int
	tail    int
	cnt     int
	cost    int
	closed  bool
	initCap int
}

// Wrap T and keep its cost.
type queuedItem[T any] struct {
	elem T
	cost int
}

// DefaultInitialCapacity of the underlying slice to keep queued elements.
const DefaultInitialCapacity = 1

// New returns a new Queue. The caller can optionally override initial capacity
// of the queue (which is DefaultInitialCapacity by default). Queue capacity never
// goes down below the initial capacity.
func New[T any](initialCapacity ...uint) *Queue[T] {
	if len(initialCapacity) > 1 {
		panic("too many arguments passed")
	}
	initCap := DefaultInitialCapacity
	if len(initialCapacity) > 0 {
		initCap = int(initialCapacity[0])
	}
	sq := &Queue[T]{
		initCap: initCap,
		nodes:   make([]queuedItem[T], initCap),
	}
	sq.cond = sync.NewCond(&sq.mu)
	return sq
}

// Add a T to the back of the queue.
// It will return false if the queue is closed, in that case the T is dropped.
func (q *Queue[T]) Add(elem T, cost int) bool {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return false
	}
	if q.cnt == len(q.nodes) {
		// Also tested a growth rate of 1.5, see: http://stackoverflow.com/questions/2269063/buffer-growth-strategy
		// In Go this resulted in a higher memory usage.
		n := q.cnt * 2
		if n == 0 {
			n = 1
		}
		q.resize(n)
	}
	i := queuedItem[T]{
		elem: elem,
		cost: cost,
	}
	q.nodes[q.tail] = i
	q.tail = (q.tail + 1) % len(q.nodes)
	q.cost += cost
	q.cnt++
	q.cond.Signal()
	q.mu.Unlock()
	return true
}

// Close the queue and discard all entries in the queue. All goroutines in wait() will return.
func (q *Queue[T]) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.closed = true
	q.cnt = 0
	q.nodes = nil
	q.cost = 0
	q.cond.Broadcast()
}

// CloseRemaining will close the queue and return all entries in the queue.
// All goroutines in wait() will return.
func (q *Queue[T]) CloseRemaining() []T {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return []T{}
	}
	rem := make([]T, 0, q.cnt)
	for q.cnt > 0 {
		i := q.nodes[q.head]
		q.head = (q.head + 1) % len(q.nodes)
		q.cnt--
		rem = append(rem, i.elem)
	}
	q.closed = true
	q.cnt = 0
	q.nodes = nil
	q.cost = 0
	q.cond.Broadcast()
	return rem
}

// Closed returns true if the queue has been closed.
// The call cannot guarantee that the queue hasn't been closed
// while the function returns, so only "true" has a definite meaning.
func (q *Queue[T]) Closed() bool {
	q.mu.RLock()
	c := q.closed
	q.mu.RUnlock()
	return c
}

// Wait for a message to be added.
// If there are items on the queue will return immediately. Will return false
// if the queue is closed. Otherwise, returns true.
func (q *Queue[T]) Wait() bool {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return false
	}
	if q.cnt != 0 {
		q.mu.Unlock()
		return true
	}
	q.cond.Wait()
	q.mu.Unlock()
	return true
}

// Remove will remove a T from the queue.
// If false is returned, it either means:
// 1) there were no items on the queue, or
// 2) the queue is closed.
func (q *Queue[T]) Remove() (T, bool) {
	q.mu.Lock()
	if q.cnt == 0 {
		q.mu.Unlock()
		var t T
		return t, false
	}
	i := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.cnt--
	q.cost -= i.cost

	if n := len(q.nodes) / 2; n >= q.initCap && q.cnt <= n {
		q.resize(n)
	}

	q.mu.Unlock()
	return i.elem, true
}

// Len returns the current length of the queue.
func (q *Queue[T]) Len() int {
	q.mu.RLock()
	l := q.cnt
	q.mu.RUnlock()
	return l
}

// Cost returns the current total cost of the queue.
func (q *Queue[T]) Cost() int {
	q.mu.RLock()
	s := q.cost
	q.mu.RUnlock()
	return s
}

// Write mutex must be held when calling.
func (q *Queue[T]) resize(n int) {
	nodes := make([]queuedItem[T], n)

	if n == 0 {
		q.tail = 0
		q.head = 0
		q.nodes = nodes
		return
	}

	if q.head < q.tail {
		copy(nodes, q.nodes[q.head:q.tail])
	} else {
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.tail])
	}

	q.tail = q.cnt % n
	q.head = 0
	q.nodes = nodes
}
