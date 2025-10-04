package or

import "sync"

// or is a function that combines one or more done channels into one.
// The return channel closes as soon as any of the source channels are closed.
var Or func(channels ...<-chan interface{}) <-chan interface{}

func init() {
	Or = func(channels ...<-chan interface{}) <-chan interface{} {
		if len(channels) == 0 {
			return nil
		}

		if len(channels) == 1 {
			return channels[0]
		}

		orDone := make(chan interface{})
		var once sync.Once

		for _, ch := range channels {
			go func(c <-chan interface{}) {
				select {
				case <-c:
					once.Do(func() {
						close(orDone)
					})
				case <-orDone:
					// Another goroutine already closed the channel
				}
			}(ch)
		}

		return orDone
	}
}
