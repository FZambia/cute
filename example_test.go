package queue

import (
	"fmt"
)

func ExampleQueue_Add() {
	q := New[string]()
	q.Add("1", 0)
	el, ok := q.Remove()
	fmt.Println(el)
	fmt.Println(ok)

	// Output:
	// 1
	// true
}

func ExampleQueue_Wait() {
	q := New[string]()
	elemConsumed := make(chan struct{})
	go func() {
		ok := q.Wait()
		if !ok {
			return
		}
		if el, ok := q.Remove(); ok {
			fmt.Println(el)
		}
		close(elemConsumed)
	}()
	q.Add("1", 0)
	<-elemConsumed

	// Output:
	// 1
}

func ExampleQueue_Cost() {
	q := New[string]()
	elemConsumed := make(chan struct{})
	go func() {
		ok := q.Wait()
		if !ok {
			fmt.Println("queue closed")
			close(elemConsumed)
			return
		}
	}()
	q.Add("s", 1)
	if q.Cost() >= 1 {
		q.Close()
	}
	<-elemConsumed

	// Output:
	// queue closed
}
