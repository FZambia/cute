package queue

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueueResize(t *testing.T) {
	q := New[string]()
	require.Equal(t, 0, q.Len())
	require.Equal(t, false, q.Closed())

	for i := 0; i < initialCapacity; i++ {
		q.Add("1", 0)
	}
	q.Add("resize here", 0)
	require.Equal(t, initialCapacity*2, cap(q.nodes))
	q.Remove()

	q.Add("new resize here", 0)
	require.Equal(t, initialCapacity*2, cap(q.nodes))
	q.Add("one more elem, no resize must happen", 0)
	require.Equal(t, initialCapacity*2, cap(q.nodes))

	require.Equal(t, initialCapacity+2, q.Len())
}

func TestQueueLen(t *testing.T) {
	q := New[string]()
	i := "12345"
	require.Equal(t, 0, q.Len())
	q.Add(i, len(i))
	q.Add(i, len(i))
	require.Equal(t, 2, q.Len())
	q.Remove()
	require.Equal(t, 1, q.Len())
}

func TestQueueSize(t *testing.T) {
	q := New[string]()
	i := "12345"
	require.Equal(t, 0, q.Cost())
	q.Add(i, len(i))
	q.Add(i, len(i))
	require.Equal(t, 10, q.Cost())
	q.Remove()
	require.Equal(t, 5, q.Cost())
}

func TestQueueWait(t *testing.T) {
	q := New[string]()
	q.Add("1", 0)
	q.Add("2", 0)

	ok := q.Wait()
	require.Equal(t, true, ok)
	s, ok := q.Remove()
	require.Equal(t, true, ok)
	require.Equal(t, "1", s)

	ok = q.Wait()
	require.Equal(t, true, ok)
	s, ok = q.Remove()
	require.Equal(t, true, ok)
	require.Equal(t, "2", s)

	go func() {
		q.Add("3", 0)
	}()

	ok = q.Wait()
	require.Equal(t, true, ok)
	s, ok = q.Remove()
	require.Equal(t, true, ok)
	require.Equal(t, "3", s)
}

func TestQueueClose(t *testing.T) {
	q := New[[]byte]()

	// test removing from empty queue
	_, ok := q.Remove()
	require.Equal(t, false, ok)

	i := []byte("1")

	q.Add(i, 0)
	q.Add(i, 0)
	q.Close()

	ok = q.Add(i, 0)
	require.Equal(t, false, ok)

	ok = q.Wait()
	require.Equal(t, false, ok)

	_, ok = q.Remove()
	require.Equal(t, false, ok)

	require.Equal(t, true, q.Closed())
}

func TestQueueCloseRemaining(t *testing.T) {
	q := New[string]()
	q.Add("1", 0)
	q.Add("2", 0)
	remaining := q.CloseRemaining()
	require.Equal(t, 2, len(remaining))
	ok := q.Add("3", 0)
	require.Equal(t, false, ok)
	require.Equal(t, true, q.Closed())
	remaining = q.CloseRemaining()
	require.Equal(t, 0, len(remaining))
}

func BenchmarkQueueAdd(b *testing.B) {
	q := New[[]byte]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		i := []byte("test")
		q.Add(i, len(i))
	}
	b.StopTimer()
	q.Close()
}

func addAndConsume[T any](q *Queue[T], item T, cost int, n int) {
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
		q.Add(item, cost)
	}
	<-done
}

func BenchmarkQueueAddConsume(b *testing.B) {
	q := New[[]byte]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addAndConsume(q, []byte("test"), 4, 10000)
		require.Equal(b, 0, q.Cost())
	}
	b.StopTimer()
	q.Close()
}
