package utils

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

func listenAndServe(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("cannot start server at \"%v\" caused %v", server.Addr, err)
	}
	log.Infof("server started at \"%v\"", server.Addr)
}

func waitInterruptSignal() {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)
	<-ch
}

func shutdownServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server was shutdown with error %v", err)
	}
	log.Info("server already shutdown")
}

func StartHttpServer(address string, handler http.Handler) {
	server := &http.Server{
		Addr:    address,
		Handler: handler,
	}
	log.Infof("starting server at \"%v\"...", server.Addr)

	go listenAndServe(server)
	waitInterruptSignal()

	log.Infof("shutdown server...")
	shutdownServer(server)
}
