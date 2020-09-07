package main

import (
	"context"
	"errors"
	httpapi "github/WowVeryLogin/multer/api/multer/http"
	"github/WowVeryLogin/multer/internal/multer/app"
	"github/WowVeryLogin/multer/internal/multer/config"
	"github/WowVeryLogin/multer/pkg/graceful"
	"os"
	"sync"
	"time"

	"log"
	"net/http"
)

func main() {
	log.SetOutput(os.Stdout)
	cfg := new(config.Config)
	if err := cfg.Parse(); err != nil {
		log.Fatalf("can not parse config: %s", err)
	}
	app := app.New(cfg)

	api := httpapi.New(app)
	defer api.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", api.HandlerFunc())
	srv := http.Server{
		Addr:    cfg.HTTPConfig.Addr,
		Handler: mux,
	}

	var wg sync.WaitGroup
	wg.Add(1)
	graceful.InitGraceful([]func(){
		func() {
			defer wg.Done()
			log.Println("closing server")
			ctx, cancel := context.WithTimeout(
				context.Background(),
				cfg.HTTPConfig.GracefulTimeout*time.Second,
			)
			defer cancel()
			err := srv.Shutdown(ctx)
			if err != nil {
				log.Fatalf("can not shutdown server: %s", err)
			}
		},
	})

	log.Println("starting server")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("error in http server: %s", err)
	}

	wg.Wait()
}
