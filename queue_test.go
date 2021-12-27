package queue

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type raw []byte

func (i raw) Size() int {
	return len(i)
}

func TestByteQueueResize(t *testing.T) {
	q := New[raw]()
	require.Equal(t, 0, q.Len())
	require.Equal(t, false, q.Closed())

	for i := 0; i < initialCapacity; i++ {
		q.Add(raw(strconv.Itoa(i)))
	}
	q.Add(raw("resize here"))
	require.Equal(t, initialCapacity*2, cap(q.nodes))
	q.Remove()

	q.Add(raw("new resize here"))
	require.Equal(t, initialCapacity*2, cap(q.nodes))
	q.Add(raw("one more item, no resize must happen"))
	require.Equal(t, initialCapacity*2, cap(q.nodes))

	require.Equal(t, initialCapacity+2, q.Len())
}

func TestByteQueueSize(t *testing.T) {
	q := New[raw]()
	require.Equal(t, 0, q.Size())
	q.Add(raw("1"))
	q.Add(raw("2"))
	require.Equal(t, 2, q.Size())
	q.Remove()
	require.Equal(t, 1, q.Size())
}

func TestByteQueueWait(t *testing.T) {
	q := New[raw]()
	q.Add(raw("1"))
	q.Add(raw("2"))

	ok := q.Wait()
	require.Equal(t, true, ok)
	s, ok := q.Remove()
	require.Equal(t, true, ok)
	require.Equal(t, "1", string(s))

	ok = q.Wait()
	require.Equal(t, true, ok)
	s, ok = q.Remove()
	require.Equal(t, true, ok)
	require.Equal(t, "2", string(s))

	go func() {
		q.Add(raw("3"))
	}()

	ok = q.Wait()
	require.Equal(t, true, ok)
	s, ok = q.Remove()
	require.Equal(t, true, ok)
	require.Equal(t, "3", string(s))
}

func TestByteQueueClose(t *testing.T) {
	q := New[raw]()

	// test removing from empty queue
	_, ok := q.Remove()
	require.Equal(t, false, ok)

	q.Add(raw("1"))
	q.Add(raw("2"))
	q.Close()

	ok = q.Add(raw("3"))
	require.Equal(t, false, ok)

	ok = q.Wait()
	require.Equal(t, false, ok)

	_, ok = q.Remove()
	require.Equal(t, false, ok)

	require.Equal(t, true, q.Closed())
}

func TestByteQueueCloseRemaining(t *testing.T) {
	q := New[raw]()
	q.Add(raw("1"))
	q.Add(raw("2"))
	remaining := q.CloseRemaining()
	require.Equal(t, 2, len(remaining))
	ok := q.Add(raw("3"))
	require.Equal(t, false, ok)
	require.Equal(t, true, q.Closed())
	remaining = q.CloseRemaining()
	require.Equal(t, 0, len(remaining))
}

func BenchmarkQueueAdd(b *testing.B) {
	q := New[raw]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Add(raw("test"))
	}
	b.StopTimer()
	q.Close()
}

func addAndConsume[T Sizeable](q *Queue[T], item T, n int) {
	// Add to queue and consume in another goroutine.
	done := make(chan struct{})
	go func() {
		count := 0
		for {
			ok := q.Wait()
			if !ok {
				continue
			}
			q.Remove()
			count++
			if count == n {
				close(done)
				break
			}
		}
	}()
	for i := 0; i < n; i++ {
		q.Add(item)
	}
	<-done
}

func BenchmarkQueueAddConsume(b *testing.B) {
	q := New[raw]()
	item := raw("test")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addAndConsume(q, item, 100)
	}
	b.StopTimer()
	q.Close()
}
