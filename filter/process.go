package filter

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")

// Service defines external service that can process batches of items.
type Service interface {
	GetLimits() (n uint64, p time.Duration)
	Process(ctx context.Context, batch Batch) error
}

// Batch is a batch of items.
type Batch []Item

// Item is some abstract item.
type Item struct{}

type Client struct {
	Service
	n  uint64
	t  time.Time
	mu sync.Mutex
}

func (c *Client) Run(ctx context.Context, items ...Item) error {
	for n := uint64(len(items)); n > 0; n = uint64(len(items)) {
		c.mu.Lock()
		t := time.Now()
		if t.After(c.t) {
			var duration time.Duration
			c.n, duration = c.GetLimits()
			c.t = t.Add(duration)
		}
		if n > c.n {
			n = c.n
		}
		c.n -= n
		c.mu.Unlock()
		if n == 0 {
			z := time.NewTimer(time.Until(c.t))
			var err error
			select {
			case <-z.C:
			case <-ctx.Done():
				err = ctx.Err()
			}
			z.Stop()
			if err != nil {
				return err
			}
			continue
		}
		err := c.Process(ctx, items[:n])
		if err != nil {
			return err
		}
		items = items[n:]
	}
	return nil
}
