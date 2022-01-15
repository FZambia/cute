package cute

import "fmt"

func ExampleQueue_Add() {
	q := New[string]()
	err := q.Add("1", 0)
	fmt.Println(err)
	el, ok, err := q.Remove()
	fmt.Println(err)
	fmt.Println(el)
	fmt.Println(ok)

	// Output:
	// <nil>
	// <nil>
	// 1
	// true
}

func ExampleQueue_Wait() {
	q := New[string]()
	elemConsumed := make(chan struct{})
	go func() {
		err := q.Wait()
		if err != nil {
			return
		}
		if el, ok, err := q.Remove(); ok {
			if err != nil {
				return
			}
			fmt.Println(ok)
			fmt.Println(el)
		}
		close(elemConsumed)
	}()
	_ = q.Add("1", 0)
	<-elemConsumed

	// Output:
	// true
	// 1
}

func ExampleQueue_Cost() {
	q := New[string](Config{MaxCost: 1})
	err := q.Add("e1", 1)
	fmt.Println(err)
	err = q.Add("e2", 1)
	fmt.Println(err)

	// Output:
	// <nil>
	// max queue cost exceeded
}
