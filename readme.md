[![Build Status](https://github.com/FZambia/queue/workflows/build/badge.svg?branch=main)](https://github.com/FZambia/queue/actions)
[![GoDoc](https://pkg.go.dev/badge/FZambia/queue)](https://pkg.go.dev/github.com/FZambia/queue)

Package queue provides a generic unbounded queue for Go programming language.

This queue implementation also maintains the total cost of items in the queue, so you can take decisions based on the current queue total cost before adding element to the queue.

One possible usage scenario could be when you need to maintain a mailbox for some client, but don't want to use a buffered channel for it since the size of the buffer is hard to estimate. It should be small most of the time but still survive occasional spikes in load. This queue package may help then.

All methods are goroutine-safe.

Requires Go >= 1.18.

Credits:

* [Erik Dubbelboer](https://github.com/erikdubbelboer) and his blog post [Faster queues in Go](https://blog.dubbelboer.com/2015/04/25/go-faster-queue.html)
* [Klaus Post](https://github.com/klauspost) and [his contribution](https://github.com/centrifugal/centrifugo/pull/23) of queue into Centrifugo project
* The code evolved a bit and finally released here as a generic queue package.
