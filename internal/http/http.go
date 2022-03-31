package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog/log"
)

var (
	mutex   sync.Mutex
	running bool         = false
	srv     *http.Server = nil
)

func Start(addr string, errCh chan<- error) error {
	mutex.Lock()
	defer mutex.Unlock()
	if running {
		return fmt.Errorf("already running")
	}
	running = true

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = colorable.NewColorableStdout()
	router := gin.Default()
	router.Use(cors.Default())            // allow all origins
	router.Use(gzip.Gzip(gzip.BestSpeed)) // enable gzip compression
	router.SetHTMLTemplate(html)

	router.GET("*any", GetEndpoint)

	srv = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Info().Msgf("Serving at http://%s", addr)
		if err := srv.ListenAndServe(); err != nil {
			mutex.Lock()
			running = false
			srv = nil
			mutex.Unlock()
			errCh <- err
		}
	}()
	return nil
}

func Stop() error {
	mutex.Lock()
	defer mutex.Unlock()
	if srv != nil {
		if err := srv.Shutdown(context.Background()); err != nil {
			return err
		}
	}
	running = false
	return nil
}
