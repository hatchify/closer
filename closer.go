package closer

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hatchify/atoms"
)

// New returns a new instance of Closer
func New() *Closer {
	var c Closer
	c.mc = make(chan error, 1)

	go c.listen()

	// Sleep for a millisecond so that listening go routine can run before we return
	time.Sleep(time.Millisecond)
	return &c
}

// Closer manages the closing of services
type Closer struct {
	// Message channel
	mc chan error
	// Closed state
	closed atoms.Bool
}

// listen will listen for closing signals (interrupt, terminate, abort, quit) and call close
func (c *Closer) listen() {
	// sc represents the signal channel
	sc := make(chan os.Signal, 1)
	// Listen for signal notifications
	// Discussion topic: Should we include SIGQUIT? If we catch the signal, we won't get to see the unwind
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	// Signal received
	<-sc
	c.Close(nil)
}

// Wait will wait until it receives a close notification
// If an error is the reason for closure - said error will be returned
func (c *Closer) Wait() (err error) {
	err = <-c.mc
	return
}

// Close will close selected instance of Closer
func (c *Closer) Close(err error) (ok bool) {
	if ok = c.closed.Set(true); !ok {
		return
	}

	c.mc <- err
	return
}
