package config

import (
	"encoding/json"
	"flag"
	"fmt"
	httpapi "github/WowVeryLogin/multer/api/multer/http/config"
	"github/WowVeryLogin/multer/pkg/httpclient"
	"os"
)

type Config struct {
	HTTPConfig       *httpapi.Config    `json:"http"`
	HTTPClientConfig *httpclient.Config `json:"client"`
}

func (c *Config) Parse() error {
	path := flag.String("c", "multer.json", "path to config file")
	flag.Parse()
	f, err := os.Open(*path)
	if err != nil {
		return fmt.Errorf("opening config file: %w", err)
	}
	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		return fmt.Errorf("parsing config file: %w", err)
	}
	return nil
}
