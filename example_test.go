package cute

import "fmt"

func ExampleQueue_Add() {
	q := New[string]()
	err := q.Add("1", 0)
	if err != nil {
		panic(err)
	}
	el, ok, err := q.Remove()
	if err != nil {
		panic(err)
	}
	fmt.Println(ok)
	fmt.Println(el)

	// Output:
	// true
	// 1
}

func ExampleQueue_Wait() {
	q := New[string]()
	elemConsumed := make(chan struct{})
	go func() {
		err := q.Wait()
		if err != nil {
			panic(err)
		}
		if el, ok, err := q.Remove(); ok {
			if err != nil {
				panic(err)
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
	q := New[string](Config{MaxCost: 10})
	err := q.Add("e1", 4)
	fmt.Println(err)
	err = q.Add("e2", 6)
	fmt.Println(err)
	err = q.Add("e3", 1)
	fmt.Println(err)

	// Output:
	// <nil>
	// <nil>
	// max queue cost exceeded
}
