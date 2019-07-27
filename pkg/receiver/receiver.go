package receiver

import "context"

// Receiver app
type Receiver struct {
}

// New Receiver instance
func New(ctx context.Context, port int) *Receiver {
	return &Receiver{}
}

// Launch application server
func (a *Receiver) Launch() error {
	return nil
}
