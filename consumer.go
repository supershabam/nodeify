package nodeify

import "time"

type Consumer struct {
	Fetcher Fetcher
	Since   time.Time
	Period  time.Duration
	err     error
	done    chan struct{}
}

func (c *Consumer) Consume() <-chan Module {
	c.done = make(chan struct{})
	out := make(chan Module)
	go func() {
		defer close(out)
		for {
			select {
			case <-c.done:
				return
			default:
				last := time.Now()
				modules, err := c.Fetcher.Fetch(c.Since)
				if err != nil {
					c.err = err
					return
				}
				for _, module := range modules {
					out <- module
				}
				c.Since = last

				// sleep inbetween polls, but also be cancelable
				select {
				case <-c.done:
					return
				case <-time.After(c.Period):
				}
			}
		}
	}()
	return out
}

func (c *Consumer) Stop() {
	close(c.done)
}

func (c Consumer) Err() error {
	return c.err
}
