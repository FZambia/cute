package queue

import (
	"testing"
)

//
//func TestQueueResize(t *testing.T) {
//	initialCapacity := 2
//	q := New[string](WithInitialCapacity(initialCapacity))
//	require.Equal(t, 0, q.Len())
//	require.Equal(t, false, q.Closed())
//
//	for i := 0; i < initialCapacity; i++ {
//		q.Add("1", 0)
//	}
//	q.Add("resize here", 0)
//	require.Equal(t, initialCapacity*2, cap(q.nodes))
//	q.Remove()
//	q.Add("new resize here", 0)
//	require.Equal(t, initialCapacity*2, cap(q.nodes))
//	q.Add("one more elem, no resize must happen", 0)
//	require.Equal(t, initialCapacity*2, cap(q.nodes))
//
//	q.Add("one more elem, resize must happen", 0)
//	require.Equal(t, initialCapacity*4, cap(q.nodes))
//
//	q.Remove()
//	require.Equal(t, initialCapacity*2, cap(q.nodes))
//	q.Remove()
//	require.Equal(t, initialCapacity*2, cap(q.nodes))
//	_, ok := q.Remove()
//	require.True(t, ok)
//	require.Equal(t, initialCapacity, cap(q.nodes))
//	_, ok = q.Remove()
//	require.True(t, ok)
//	require.Equal(t, initialCapacity, cap(q.nodes))
//	_, ok = q.Remove()
//	require.True(t, ok)
//	require.Equal(t, initialCapacity, cap(q.nodes))
//	_, ok = q.Remove()
//	require.False(t, ok)
//	require.Equal(t, initialCapacity, cap(q.nodes))
//}
//
//func TestQueueResizeToZero(t *testing.T) {
//	q := New[string](0)
//	require.Equal(t, 0, q.Len())
//	require.Equal(t, false, q.Closed())
//
//	q.Add("resize here", 0)
//	require.Equal(t, 1, cap(q.nodes))
//	q.Remove()
//	require.Equal(t, 0, cap(q.nodes))
//	q.Add("resize here", 0)
//	require.Equal(t, 1, cap(q.nodes))
//}
//
//func TestQueueLen(t *testing.T) {
//	q := New[string]()
//	i := "12345"
//	require.Equal(t, 0, q.Len())
//	q.Add(i, len(i))
//	q.Add(i, len(i))
//	require.Equal(t, 2, q.Len())
//	q.Remove()
//	require.Equal(t, 1, q.Len())
//}
//
//func TestQueueNew_Panics(t *testing.T) {
//	require.Panics(t, func() {
//		New[string](2, 4)
//	})
//}
//
//func TestQueueCost(t *testing.T) {
//	q := New[string]()
//	i := "12345"
//	require.Equal(t, 0, q.Cost())
//	q.Add(i, len(i))
//	q.Add(i, len(i))
//	require.Equal(t, 10, q.Cost())
//	q.Remove()
//	require.Equal(t, 5, q.Cost())
//}
//
//func TestQueueWait(t *testing.T) {
//	q := New[string]()
//	q.Add("1", 0)
//	q.Add("2", 0)
//
//	ok := q.Wait()
//	require.Equal(t, true, ok)
//	s, ok := q.Remove()
//	require.Equal(t, true, ok)
//	require.Equal(t, "1", s)
//
//	ok = q.Wait()
//	require.Equal(t, true, ok)
//	s, ok = q.Remove()
//	require.Equal(t, true, ok)
//	require.Equal(t, "2", s)
//
//	go func() {
//		q.Add("3", 0)
//	}()
//
//	ok = q.Wait()
//	require.Equal(t, true, ok)
//	s, ok = q.Remove()
//	require.Equal(t, true, ok)
//	require.Equal(t, "3", s)
//}
//
//func TestQueueClose(t *testing.T) {
//	q := New[[]byte]()
//
//	// test removing from empty queue
//	_, ok := q.Remove()
//	require.Equal(t, false, ok)
//
//	i := []byte("1")
//
//	q.Add(i, 0)
//	q.Add(i, 0)
//	q.Close()
//
//	ok = q.Add(i, 0)
//	require.Equal(t, false, ok)
//
//	ok = q.Wait()
//	require.Equal(t, false, ok)
//
//	_, ok = q.Remove()
//	require.Equal(t, false, ok)
//
//	require.Equal(t, true, q.Closed())
//}
//
//func TestQueueCloseRemaining(t *testing.T) {
//	q := New[string]()
//	q.Add("1", 0)
//	q.Add("2", 0)
//	remaining := q.CloseRemaining()
//	require.Equal(t, 2, len(remaining))
//	ok := q.Add("3", 0)
//	require.Equal(t, false, ok)
//	require.Equal(t, true, q.Closed())
//	remaining = q.CloseRemaining()
//	require.Equal(t, 0, len(remaining))
//}

var testQueue *Queue[[]byte]

func BenchmarkQueueInit(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testQueue = New[[]byte](WithInitialCapacity(2), WithMaxCost(10))
	}
	b.StopTimer()
}

//
//func BenchmarkQueueAdd(b *testing.B) {
//	q := New[[]byte]()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		i := []byte("test")
//		q.Add(i, len(i))
//	}
//	b.StopTimer()
//	q.Close()
//}
//
//func addAndConsume[T any](q *Queue[T], item T, cost int, n int) {
//	// Add to queue and consume in another goroutine.
//	done := make(chan struct{})
//	go func() {
//		count := 0
//		for {
//			ok := q.Wait()
//			if !ok {
//				continue
//			}
//			q.Remove()
//			count++
//			if count == n {
//				close(done)
//				break
//			}
//		}
//	}()
//	for i := 0; i < n; i++ {
//		q.Add(item, cost)
//	}
//	<-done
//}
//
//func BenchmarkQueueAddConsume(b *testing.B) {
//	q := New[[]byte]()
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		addAndConsume(q, []byte("test"), 4, 100)
//		require.Equal(b, 0, q.Cost())
//	}
//	b.StopTimer()
//	q.Close()
//}
//
//func BenchmarkQueueCreateAddConsumeClose(b *testing.B) {
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		q := New[[]byte](100)
//		addAndConsume(q, []byte("test"), 4, 100)
//		require.Equal(b, 0, q.Cost())
//		q.Close()
//	}
//	b.StopTimer()
//}
