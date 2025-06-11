package poller

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

var ErrNoData = errors.New("no contains any data")

type TaskFunc func() error

type Poller struct {
	BaseDelay time.Duration
	MaxDelay  time.Duration
	Jitter    bool

	stop   chan struct{}
	closed atomic.Bool
}

func NewPoller() *Poller {
	return &Poller{
		BaseDelay: 5 * time.Second,
		MaxDelay:  5 * time.Second,
		Jitter:    false,
		stop:      make(chan struct{}, 1),
	}
}

func (p *Poller) SetBaseDelay(delay time.Duration) *Poller {
	p.BaseDelay = delay
	return p
}

func (p *Poller) SetMaxDelay(delay time.Duration) *Poller {
	p.MaxDelay = delay
	return p
}

func (p *Poller) SetJitter(jitter bool) *Poller {
	p.Jitter = jitter
	return p
}

func (p *Poller) Spawn(task TaskFunc) {
	log.Println("Poller is ready")

	attempt := 0

	for {
		select {
		case <-p.stop:
			return
		default:
			delay := p.backoff(attempt)

			t := time.NewTimer(delay)

			select {
			case <-t.C:
				err := task()
				if err != nil && errors.Is(err, ErrNoData) {
					attempt++
				} else {
					attempt = 0
				}
			case <-p.stop:
				t.Stop()
				return
			}
		}
	}
}

func (p *Poller) Close() error {
	if p.closed.CompareAndSwap(false, true) {
		close(p.stop)
	}

	return nil
}

func (p *Poller) backoff(attempt int) time.Duration {
	exp := math.Pow(2, float64(attempt))
	delay := time.Duration(float64(exp) * float64(p.BaseDelay))

	if p.Jitter {
		delay = time.Duration(rand.Float64() * float64(delay))
	}

	if p.isMaxDelay(delay) {
		return p.MaxDelay
	}

	return delay
}

func (p *Poller) isMaxDelay(delay time.Duration) bool {
	return p.MaxDelay > 0 && delay > p.MaxDelay
}
