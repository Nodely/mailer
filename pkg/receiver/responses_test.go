package receiver

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrRender(t *testing.T) {
	render := ErrRender(errors.New("Some critical error"))
	assert.NotNil(t, render)
}
