package main

import (
	"errors"
	"mcp/bff/router"
	"mcp/core/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	log.InitLogx()
	gin.SetMode(gin.DebugMode)

	srv := &http.Server{
		Addr:    ":8090",
		Handler: router.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logrus.Info("Shutting down server...")

	// rpc.ShutdownClients()

	logrus.Info("exited")
}
