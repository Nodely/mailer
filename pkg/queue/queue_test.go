package queue

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalQueue(t *testing.T) {
	q, err := NewQueue()
	assert.Nil(t, err)

	err = q.Send(&Message{
		TemplateID: "some-id",
		Locale:     "en_US",
	})
	assert.Nil(t, err)

	// read message
	msg, err := q.Read()
	assert.Nil(t, err)
	assert.Equal(t, "some-id", msg.TemplateID)

	// read empty queue
	msg, err = q.Read()
	assert.Nil(t, err)
	assert.Nil(t, msg)
}

func TestMessageBind(t *testing.T) {
	m := &Message{
		TemplateID: "some-id",
	}

	r, _ := http.NewRequest("GET", "/", nil)

	err := m.Bind(r)
	assert.Nil(t, err)

	m = &Message{}
	err = m.Bind(r)
	assert.NotNil(t, err)
}
