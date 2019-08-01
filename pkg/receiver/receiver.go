package receiver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/nodely/notify/pkg/queue"
	uuid "github.com/satori/go.uuid"
)

// Receiver app
type Receiver struct {
	port  string
	queue queue.Queue
}

// New Receiver instance
func New(ctx context.Context, port int, queue queue.Queue) *Receiver {
	return &Receiver{
		port:  string(port),
		queue: queue,
	}
}

// Launch application server
func (a *Receiver) Launch() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(router chi.Router) {
		router.Post("/notify", a.notify)
	})

	return http.ListenAndServe(":"+a.port, r)
}

func (a *Receiver) notify(w http.ResponseWriter, r *http.Request) {
	var data queue.Message
	if err := render.Bind(r, &data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// generate message id
	data.ID = uuid.Must(uuid.NewV4(), nil).String()

	log.Printf("Message Queued with ID: [%s]", data.ID)

	// send message to queue
	a.queue.Send(&data)

	w.Write([]byte("{message_id:\"" + data.ID + "\"}"))
}
