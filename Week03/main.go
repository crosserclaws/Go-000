package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

const timeout = 10 * time.Second

func main() {
	log.Println("[Program] Start.")

	deadlineCtx, deadlineCancel := context.WithTimeout(context.Background(), timeout)
	defer deadlineCancel()

	eg, ctx := errgroup.WithContext(deadlineCtx)
	svr, registeredShutdownCh := GetDefaultServer()

	// NOTE: We intentionally start a go routine by errorgroups in these functions for readability.
	goCaptureSignals(ctx, eg)
	goServer(ctx, eg, svr)
	goNotifyShutdownByContex(ctx, eg, svr)

	// Wait routines
	log.Println("[Program] Wait routines.")
	if err := eg.Wait(); err != nil {
		log.Println("Failed in an error group:", err)
	}

	// Wait shutdown
	log.Println("[Program] Ensure the shutdown is graceful.")
	<-registeredShutdownCh

	log.Println("[Program] End.")
}

func goCaptureSignals(ctx context.Context, eg *errgroup.Group) {
	eg.Go(func() error {
		log.Println("[Catcher] Start.")
		// The signals SIGKILL and SIGSTOP may not be caught by a program.
		ch := RegisterSignals(syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		var err error

		select {
		case s := <-ch:
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("Ask server to shutdown when capturing a registered signal:", s)
				err = fmt.Errorf("Capture a registered signal: %v", s)
			default:
				log.Println("Capture an unknown signal:", s)
			}
		case <-ctx.Done():
			log.Println("[Catcher] The context is done.")
		}
		log.Println("[Catcher] End.")
		return err
	})
}

func goServer(ctx context.Context, eg *errgroup.Group, svr *http.Server) {
	eg.Go(func() error {
		log.Println("[Server] Start.")
		err := svr.ListenAndServe()
		log.Println("[Server] End.")
		return err
	})
}

func goNotifyShutdownByContex(ctx context.Context, eg *errgroup.Group, svr *http.Server) {
	eg.Go(func() error {
		log.Println("[Notifier] Start.")
		select {
		case <-ctx.Done():
			log.Println("[Notifier] The context is done.")
			log.Println("[Notifier] End.")
			return svr.Shutdown(ctx)
		}
	})
}

// GetDefaultServer returns a default server with port 8080 and a logging msg when shutdown.
func GetDefaultServer() (*http.Server, <-chan struct{}) {
	ch := make(chan struct{})
	svr := &http.Server{Addr: ":8080"}
	svr.RegisterOnShutdown(func() {
		log.Println("Do registered shutdown.")
		close(ch)
	})
	return svr, ch
}

// RegisterSignals is a utility function registers given signals.
func RegisterSignals(sig ...os.Signal) <-chan os.Signal {
	log.Println("Register signals:", sig)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig...)
	return ch
}
