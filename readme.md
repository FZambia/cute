[![Build Status](https://github.com/FZambia/cute/workflows/build/badge.svg?branch=main)](https://github.com/FZambia/cute/actions)
[![GoDoc](https://pkg.go.dev/badge/FZambia/cute)](https://pkg.go.dev/github.com/FZambia/cute)
[![codecov.io](https://codecov.io/gh/FZambia/cute/branch/main/graphs/badge.svg)](https://codecov.io/github/FZambia/cute?branch=main)

Package `cu[T]e` provides a generic unbounded FIFO queue for Go programming language.

This queue implementation additionally maintains the total cost of currently queued elements, so the queue can be bound to the cost of the elements.

One possible usage scenario could be when you need to maintain a mailbox for some client, but don't want to use a buffered channel for it since the size of the buffer is hard to estimate. It should be small most of the time but still survive occasional spikes in load. This queue package may help then.

All methods are goroutine-safe.

Requires Go >= 1.18 since uses type parameters (generics).

Credits:

* [Erik Dubbelboer](https://github.com/erikdubbelboer) and his blog post [Faster queues in Go](https://blog.dubbelboer.com/2015/04/25/go-faster-queue.html)
* [Klaus Post](https://github.com/klauspost) and [his contribution](https://github.com/centrifugal/centrifugo/pull/23) of queue into [Centrifugo](https://github.com/centrifugal/centrifugo) project
* The code evolved a bit and was finally released here as a generic queue package.
