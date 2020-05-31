/*
Package catena implements a simple social graph API.
*/
package catena

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/bbengfort/catena/config"
	"github.com/bbengfort/catena/logs"
	"github.com/julienschmidt/httprouter"

	// register database drivers
	_ "github.com/lib/pq"
)

// Version of the Catena server and package.
const Version = "v0.1"

// Content types
const ctjson = "application/json; charset=utf-8"

// New creates a Catena API server with the specified options and returns it.
func New(conf config.Config) (api *Catena, err error) {
	// TODO: validate the config

	// Implement basic requests logger
	// TODO: add logger config to config
	logger := logs.New("catena")
	logger.EnableColors()
	logger.SetBackend(os.Stdout)
	logger.SetLogLevel(logs.LevelInfo)

	// TODO: add config to routes
	// TODO: better middleware handling
	mux := Routes()
	handler := logs.NewHTTPLogger("http", mux)
	handler.EnableColors()

	server := &http.Server{
		Addr:         conf.BindAddr(),
		Handler:      handler,
		ErrorLog:     log.New(os.Stderr, "[http] ", log.LstdFlags),
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		IdleTimeout:  conf.IdleTimeout,
	}

	return &Catena{
		conf:    conf,
		mux:     mux,
		server:  server,
		logger:  logger,
		healthy: false,
		done:    make(chan bool),
	}, nil
}

// Catena is an API server.
type Catena struct {
	sync.RWMutex
	conf    config.Config
	db      *sql.DB
	mux     *httprouter.Router
	server  *http.Server
	logger  *logs.Logger
	healthy bool
	done    chan bool
}

// Serve the API
func (c *Catena) Serve() (err error) {
	// set healthy before starting the server
	c.setHealth(true)

	// capture os signals to gracefully shutdown
	c.osSignals()

	// listen and serve
	// TODO: handle serveTLS
	c.logger.Status("server is ready to handle requests at %s", c.conf.Endpoint())
	if err = c.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	// wait until shutdown is complete
	<-c.done
	c.logger.Status("server(s) stopped")
	return nil
}

// Shutdown the API server gracefully
func (c *Catena) Shutdown() (err error) {
	c.logger.Caution("shutting down server(s)...")
	c.setHealth(false)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c.server.SetKeepAlivesEnabled(false)
	if err := c.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("could not gracefully shutdown server: %s", err)
	}

	close(c.done)
	return nil
}

func (c *Catena) setHealth(health bool) {
	c.Lock()
	c.healthy = health
	c.Unlock()
}

func (c *Catena) osSignals() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		c.Shutdown()
	}()
}
