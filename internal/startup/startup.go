package startup

import (
	"flag"
	"fmt"

	"github.com/lorenzogood/x/internal/flagenv"
	"go.uber.org/zap"
)

var logLevel = zap.LevelFlag("log-level", zap.InfoLevel, "zap log level.")

func Run(envPrefix string) {
	if envPrefix != "" {
		flagenv.Prefix = envPrefix + "_"
	}

	flag.Parse()
	flagenv.Parse()
	flag.Parse()

	cfg := zap.NewProductionConfig()

	cfg.Level.SetLevel(*logLevel)

	l, err := cfg.Build(
		zap.WithCaller(false),
		zap.AddStacktrace(zap.PanicLevel),
	)
	if err != nil {
		panic(fmt.Errorf("error initializing logger: %w", err))
	}

	zap.ReplaceGlobals(l)
}