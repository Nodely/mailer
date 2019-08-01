package queue

import (
	"container/list"
	"errors"
	"fmt"
	"net/http"
)

// Queue interface
type Queue interface {
	Send(msg *Message) error
	Read() (*Message, error)
}

// LocalQueue internal immplementation
type LocalQueue struct {
	q *list.List
}

// Message struct
type Message struct {
	ID         string                 `json:"message_id"`
	TemplateID string                 `json:"template"`
	Locale     string                 `json:"locale"`
	Values     map[string]interface{} `json:"props"`
}

// Bind func
func (a *Message) Bind(r *http.Request) error {
	if a.TemplateID == "" {
		return errors.New("template id is not defined")
	}
	if a.Locale == "" {
		a.Locale = "en_US"
	}
	fmt.Printf("Bind: %+v", a)
	return nil
}

// NewQueue func
func NewQueue() (Queue, error) {
	return &LocalQueue{
		q: list.New(),
	}, nil
}

// Send func
func (app *LocalQueue) Send(msg *Message) error {
	app.q.PushBack(msg)
	return nil
}

// Read func
func (app *LocalQueue) Read() (*Message, error) {
	if app.q.Len() > 0 {
		e := app.q.Front()
		app.q.Remove(e)
		return e.Value.(*Message), nil
	}
	return nil, nil
}
