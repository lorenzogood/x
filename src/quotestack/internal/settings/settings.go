package settings

import (
	"log/slog"
	"os"
	"sync"

	"github.com/hjson/hjson-go/v4"
)

type Settings struct {
	Port int `json:"web_port"`
}

var (
	settings_init sync.Once
	settings      Settings
)

func Get() *Settings {
	settings_init.Do(func() {
		location := os.Getenv("QUOTESTACK_CONFIG_LOCATION")
		if location == "" {
			slog.Error("failed while locating configuration", "error", "QUOTESTACK_CONFIG_LOCATION is unset")
			os.Exit(1)
		}

		contents, err := os.ReadFile(location)
		if err != nil {
			slog.Error("failed while loading configuration", "error", err)
			os.Exit(1)
		}

		err = hjson.Unmarshal(contents, &settings)
		if err != nil {
			slog.Error("failed while unmarshalling configuration", "error", err)
			os.Exit(1)
		}
	})

	return &settings
}
