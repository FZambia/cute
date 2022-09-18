package cute

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueueResize(t *testing.T) {
	initialCapacity := 2
	q := New[string](Config{InitialCapacity: initialCapacity})
	require.Equal(t, 0, q.Len())
	require.Equal(t, false, q.Closed())

	for i := 0; i < initialCapacity; i++ {
		_ = q.Add("1", 0)
	}
	_ = q.Add("resize here", 0)
	require.Equal(t, initialCapacity*2, cap(q.nodes))
	_, _, _ = q.Remove()
	_ = q.Add("new resize here", 0)
	require.Equal(t, initialCapacity*2, cap(q.nodes))
	_ = q.Add("one more elem, no resize must happen", 0)
	require.Equal(t, initialCapacity*2, cap(q.nodes))

	_ = q.Add("one more elem, resize must happen", 0)
	require.Equal(t, initialCapacity*4, cap(q.nodes))

	_, _, _ = q.Remove()
	require.Equal(t, initialCapacity*2, cap(q.nodes))
	_, _, _ = q.Remove()
	require.Equal(t, initialCapacity*2, cap(q.nodes))
	_, ok, err := q.Remove()
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, initialCapacity, cap(q.nodes))
	_, ok, err = q.Remove()
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, initialCapacity, cap(q.nodes))
	_, ok, err = q.Remove()
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, initialCapacity, cap(q.nodes))
	_, ok, err = q.Remove()
	require.NoError(t, err)
	require.False(t, ok)
	require.Equal(t, initialCapacity, cap(q.nodes))
}

func TestQueueResizeToZero(t *testing.T) {
	q := New[string]()
	require.Equal(t, 0, q.Len())
	require.Equal(t, false, q.Closed())

	_ = q.Add("resize here", 0)
	require.Equal(t, 1, cap(q.nodes))
	_, _, _ = q.Remove()
	require.Equal(t, 0, cap(q.nodes))
	_ = q.Add("resize here", 0)
	require.Equal(t, 1, cap(q.nodes))
}

func TestQueueLen(t *testing.T) {
	q := New[string]()
	i := "12345"
	require.Equal(t, 0, q.Len())
	_ = q.Add(i, len(i))
	_ = q.Add(i, len(i))
	require.Equal(t, 2, q.Len())
	_, _, _ = q.Remove()
	require.Equal(t, 1, q.Len())
}

func TestQueueNew_Panics(t *testing.T) {
	require.Panics(t, func() {
		New[string](Config{}, Config{})
	})
}

func TestQueueCost(t *testing.T) {
	q := New[string]()
	i := "12345"
	require.Equal(t, 0, q.Cost())
	_ = q.Add(i, len(i))
	_ = q.Add(i, len(i))
	require.Equal(t, 10, q.Cost())
	_, _, _ = q.Remove()
	require.Equal(t, 5, q.Cost())
}

func TestQueueWait(t *testing.T) {
	q := New[string]()
	_ = q.Add("1", 0)
	_ = q.Add("2", 0)

	err := q.Wait()
	require.NoError(t, err)
	s, ok, err := q.Remove()
	require.NoError(t, err)
	require.Equal(t, true, ok)
	require.Equal(t, "1", s)

	err = q.Wait()
	require.NoError(t, err)
	s, ok, err = q.Remove()
	require.NoError(t, err)
	require.Equal(t, true, ok)
	require.Equal(t, "2", s)

	go func() {
		_ = q.Add("3", 0)
	}()

	err = q.Wait()
	require.NoError(t, err)
	s, ok, err = q.Remove()
	require.NoError(t, err)
	require.Equal(t, true, ok)
	require.Equal(t, "3", s)
}

func TestQueueClose(t *testing.T) {
	q := New[[]byte]()

	// test removing from empty queue
	_, ok, err := q.Remove()
	require.NoError(t, err)
	require.Equal(t, false, ok)

	i := []byte("1")

	_ = q.Add(i, 0)
	_ = q.Add(i, 0)
	q.Close()

	err = q.Add(i, 0)
	require.ErrorIs(t, err, ErrClosed)

	err = q.Wait()
	require.ErrorIs(t, err, ErrClosed)

	_, _, err = q.Remove()
	require.ErrorIs(t, err, ErrClosed)

	require.True(t, q.Closed())
}

func TestQueueCloseRemaining(t *testing.T) {
	q := New[string]()
	_ = q.Add("1", 0)
	_ = q.Add("2", 0)
	remaining := q.CloseRemaining()
	require.Equal(t, 2, len(remaining))
	err := q.Add("3", 0)
	require.ErrorIs(t, err, ErrClosed)
	require.True(t, q.Closed())
	remaining = q.CloseRemaining()
	require.Equal(t, 0, len(remaining))
}

var testQueue *Queue[[]byte]

func BenchmarkQueueNew(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testQueue = New[[]byte](Config{InitialCapacity: 2})
	}
	b.StopTimer()
}

func BenchmarkQueueAdd(b *testing.B) {
	q := New[[]byte]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		i := []byte("test")
		_ = q.Add(i, len(i))
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
			err := q.Wait()
			if err != nil {
				break
			}
			_, _, _ = q.Remove()
			count++
			if count == n {
				close(done)
				break
			}
		}
	}()
	for i := 0; i < n; i++ {
		_ = q.Add(item, cost)
	}
	<-done
}

func BenchmarkQueueAddConsume(b *testing.B) {
	q := New[[]byte](Config{InitialCapacity: 1})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addAndConsume(q, []byte("test"), 4, 100)
		require.Equal(b, 0, q.Cost())
	}
	b.StopTimer()
	q.Close()
}

func BenchmarkQueueCreateAddConsumeClose(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := New[[]byte](Config{InitialCapacity: 100})
		addAndConsume(q, []byte("test"), 4, 100)
		require.Equal(b, 0, q.Cost())
		q.Close()
	}
	b.StopTimer()
}
