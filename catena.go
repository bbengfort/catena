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

	"github.com/julienschmidt/httprouter"

	// register database drivers
	_ "github.com/lib/pq"
)

// Version of the Catena server and package.
const Version = "v0.1"

// Content types
const ctjson = "application/json; charset=utf-8"

// New creates a Catena API server with the specified options and returns it.
func New() (api *Catena, err error) {
	// Implement basic requests logger
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	mux := Routes()
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return &Catena{
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
	db      *sql.DB
	mux     *httprouter.Router
	server  *http.Server
	logger  *log.Logger
	healthy bool
	done    chan bool
}

// Serve the API
func (c *Catena) Serve() (err error) {
	// set healthy before starting the server
	c.setHealth(true)

	// capture os signals to gracefully shutdown
	c.signals()

	// listen and serve
	// TODO: handle serveTLS
	c.logger.Println("server is ready to handle requests at", "http://localhost:8080")
	if err = c.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	// wait until shutdown is complete
	<-c.done
	c.logger.Println("server(s) stopped")
	return nil
}

// Shutdown the API server gracefully
func (c *Catena) Shutdown() (err error) {
	c.logger.Println("shutting down server(s)...")
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

func (c *Catena) signals() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		c.Shutdown()
	}()
}
