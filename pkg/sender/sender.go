package sender

import "context"

// Sender app
type Sender struct {
}

// New Sender app instance
func New(ctx context.Context) *Sender {
	return &Sender{}
}

// Launch queue listening
func (a *Sender) Launch() error {
	return nil
}
