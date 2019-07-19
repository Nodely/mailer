package processor

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newProcessor(locale string) *Processor {
	return New()
}

var testEmailPayload = `
{
	"template": "welcome",
	"locale": "en_US"
}
`

func TestNewProcessor(t *testing.T) {
	proc := newProcessor("en_US")
	err := proc.Do([]byte(testEmailPayload))
	assert.Nil(t, err, "error should be nil")
}

func TestProcessorInvalidPayload(t *testing.T) {
	proc := newProcessor("en_US")

	log.Print("Trying to pass empty payload")
	err := proc.Do(nil)
	assert.NotNil(t, err, "empty payload should be handled")

	log.Print("Trying to pass invalid payload")
	err = proc.Do([]byte("some non json string"))
	assert.NotNil(t, err, "invalid payload should be handled")

	log.Print("Trying to pass wrong template id")
	err = proc.Do([]byte(`{"template": "welcome-wrong"}`))
	assert.NotNil(t, err, "invalid template id is not handled")
}
