package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	httpconfig "github/WowVeryLogin/multer/api/multer/http/config"
	"github/WowVeryLogin/multer/internal/multer"
	"github/WowVeryLogin/multer/internal/multer/app"
	"github/WowVeryLogin/multer/pkg/httpclient"
	"github/WowVeryLogin/multer/pkg/pool"
	"github/WowVeryLogin/multer/pkg/pool/flexibale"
	"log"
	"net/http"
	"sync"
	"time"
)

func validateRequest(r *http.Request) (*multer.Request, error) {
	req := new(multer.Request)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, fmt.Errorf("decoding request: %w", err)
	}
	if len(req.Urls) > 20 {
		return nil, fmt.Errorf("request is too big: %d", len(req.Urls))
	}
	return req, nil
}

func sendError(w http.ResponseWriter, format string, err error) error {
	data := []byte(fmt.Sprintf(format, err))
	n, err := w.Write(data)
	if err != nil {
		return fmt.Errorf("writing response: %w", err)
	}
	if n != len(data) {
		return errors.New("could not write all data")
	}
	return nil
}

type HTTPApi struct {
	*httpconfig.Config
	p    pool.Pool
	cl   httpclient.Client
	done chan struct{}
}

func New(app *app.Application) *HTTPApi {
	api := &HTTPApi{
		Config: app.HTTPConfig,
		p: flexibale.New(&flexibale.Config{
			MaxWorkers: app.HTTPConfig.MaxRequests,
		}),
		cl:   app.HTTPClient,
		done: make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-api.done:
				return
			case err, ok := <-api.p.Errors():
				if !ok {
					return
				}
				log.Printf("http api error: %s", err)
			}
		}
	}()
	return api
}

func (h *HTTPApi) Close() {
	h.done <- struct{}{}
	close(h.done)
	h.p.Close()
}

func (h *HTTPApi) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var wg sync.WaitGroup
		wg.Add(1)
		h.p.Put(r.Context(), pool.Task(func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, h.RequestTimeout*time.Second)
			defer cancel()

			defer wg.Done()
			log.Println("processing request")
			req, err := validateRequest(r)
			if err != nil {
				log.Printf("error: %s", err)
				w.WriteHeader(http.StatusBadRequest)
				return sendError(w, "bad request params: %s", err)
			}

			m := multer.New(h.cl, flexibale.New(h.Config.RequestsPoolConfig))
			defer m.Close()

			resp, err := m.HandleRequest(ctx, *req)
			if err != nil {
				log.Printf("error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return sendError(w, "internal error: %s", err)
			}

			data, err := json.Marshal(resp)
			if err != nil {
				log.Printf("error: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return sendError(w, "internal error: %s", err)
			}

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(data)
			return err
		}))
		wg.Wait()
	}
}
