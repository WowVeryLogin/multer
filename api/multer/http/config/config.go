package config

import (
	"github/WowVeryLogin/multer/pkg/pool/flexibale"
	"time"
)

type Config struct {
	GracefulTimeout    time.Duration     `json:"graceful_timeout"`
	MaxRequests        int               `json:"max_request"`
	MaxURLs            int               `json:"max_urls"`
	RequestTimeout     time.Duration     `json:"timeout"`
	Addr               string            `json:"addr"`
	RequestsPoolConfig *flexibale.Config `json:"pool"`
}
