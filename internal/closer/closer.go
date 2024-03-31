package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	funcs []func() error
	done  chan struct{}
}

var globalCloser = New()

// New returns new Closer, if []os.Signal is specified Closer will automatically call CloseAll when one of signals is received from OS
func New(signals ...os.Signal) *Closer {
	c := Closer{done: make(chan struct{})}
	if len(signals) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, signals...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}

	return &c
}

func (c *Closer) Add(f ...func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

func Add(f ...func() error) {
	globalCloser.Add(f...)
}

func (c *Closer) Wait() {
	<-c.done
}

func Wait() {
	globalCloser.Wait()
}

func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errCh := make(chan error, len(funcs))
		for _, f := range funcs {
			f := f
			go func() {
				errCh <- f()
			}()
		}

		for i := 0; i < cap(errCh); i++ {
			if err := <-errCh; err != nil {
				log.Println("error returned from Closer")
			}
		}
	})
}

func CloseAll() {
	globalCloser.CloseAll()
}
