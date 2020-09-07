package app

import (
	"github/WowVeryLogin/multer/internal/multer/config"
	"github/WowVeryLogin/multer/pkg/httpclient"
)

type Application struct {
	*config.Config
	HTTPClient httpclient.Client
}

func New(cfg *config.Config) *Application {
	return &Application{
		Config:     cfg,
		HTTPClient: httpclient.New(cfg.HTTPClientConfig),
	}
}
