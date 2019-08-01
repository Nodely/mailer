package receiver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nodely/notify/pkg/queue"
	"github.com/stretchr/testify/assert"
)

func TestNotfy(t *testing.T) {
	q, err := queue.NewQueue()
	assert.Nil(t, err)

	r := New(context.TODO(), 0, q)
	r.Launch()

	proceedRequest(t, queue.Message{}, r.notify, func(w *httptest.ResponseRecorder) {
		assert.Equal(t, 400, w.Code)
	})

	proceedRequest(t, queue.Message{
		TemplateID: "some-id",
	}, r.notify, func(w *httptest.ResponseRecorder) {
		fmt.Printf("%s\n", w.Body)
		assert.Equal(t, 200, w.Code)
	})
}

func proceedRequest(t *testing.T, input queue.Message, handleFn func(w http.ResponseWriter, r *http.Request), callback func(w *httptest.ResponseRecorder)) {

	bytInput, _ := json.Marshal(input)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(bytInput))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("content-type", "application/json")

	rr := httptest.NewRecorder()

	nextHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//ctx := context.WithValue(r.Context(), data.ContextLoggerKey, log)
			next.ServeHTTP(w, r)
		})
	}

	handler := nextHandler(http.HandlerFunc(handleFn))
	handler.ServeHTTP(rr, req)

	callback(rr)
}
