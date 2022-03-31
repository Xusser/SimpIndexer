package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Xusser/SimpIndexer/build"
	"github.com/Xusser/SimpIndexer/internal/config"
	"github.com/Xusser/SimpIndexer/internal/http"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	o, _ := os.Stdout.Stat()
	if o.Mode()&os.ModeCharDevice == os.ModeCharDevice {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: colorable.NewColorableStdout(), TimeFormat: time.RFC3339})
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	log.Info().Str("BuildDate", build.BuildDate).Str("BuildCommit", build.BuildCommit).Msg("SimpIndexer")

	flag.Parse()
	if config.FlagTraceLog {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else if config.FlagDebugLog {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Trace().Msg("Initializing config")
	if err := config.Init(); err != nil {
		log.Fatal().Err(err).Msg("Fail to initialize config")
		return
	}

	configJson, _ := json.Marshal(config.Get())
	log.Info().RawJSON("config", configJson).Msg("Config initialized")

	httpCh := make(chan error, 1)
	signalCh := make(chan os.Signal, 1)

	if err := http.Start(fmt.Sprintf("%s:%d", config.Get().Host, config.Get().Port), httpCh); err != nil {
		log.Fatal().Err(err).Msg("Fail to start http module")
		return
	}

	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
	select {
	case s := <-signalCh:
		fmt.Println("")
		log.Warn().Msgf("Catch signal[%s]", s.String())
	case err := <-httpCh:
		fmt.Println("")
		log.Error().Err(err).Msg("Unexpected error from http module")
	}

	log.Info().Msg("Application exit")
}
