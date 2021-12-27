[![Build Status](https://github.com/FZambia/queue/workflows/build/badge.svg?branch=main)](https://github.com/FZambia/queue/actions)
[![GoDoc](https://pkg.go.dev/badge/FZambia/queue)](https://pkg.go.dev/github.com/FZambia/queue)

Package queue provides a generic unbounded queue for Go programming language.

This queue implementation can optionally maintain the total cost of items in the queue, so you can take decisions based on the current max queue total cost.

All methods are goroutine-safe.

Requires Go >= 1.18.

Credits:

* [Erik Dubbelboer](https://github.com/erikdubbelboer) and his blog post [Faster queues in Go](https://blog.dubbelboer.com/2015/04/25/go-faster-queue.html)
* [Klaus Post](https://github.com/klauspost) and [his contribution](https://github.com/centrifugal/centrifugo/pull/23) of queue into Centrifugo project
* The code evolved a bit and finally released here as a generic queue package.