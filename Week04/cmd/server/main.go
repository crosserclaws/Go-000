package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pb "week04/api/video/v1"
)

const timeout = 10 * time.Second

const (
	port = ":50051"
)

func main() {
	log.Println("[Program] Start.")

	svr := grpc.NewServer()
	svc := InitializeVideoInfoService()
	pb.RegisterVideoInformerServer(svr, svc)

	lis := getListenerOrFail()
	eg, ctx := errgroup.WithContext(context.Background())
	// NOTE: We intentionally start a go routine by errorgroups in these functions for readability.
	goCaptureSignals(ctx, eg)
	goServer(ctx, eg, svr, lis)

	// Wait routines
	log.Println("[Program] Wait routines.")
	if err := eg.Wait(); err != nil {
		log.Println("Failed in an error group:", err)
	}
	log.Println("[Program] End.")
}

func getListenerOrFail() net.Listener {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	return lis
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

// RegisterSignals is a utility function registers given signals.
func RegisterSignals(sig ...os.Signal) <-chan os.Signal {
	log.Println("Register signals:", sig)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig...)
	return ch
}

func goServer(ctx context.Context, eg *errgroup.Group, svr *grpc.Server, lis net.Listener) {
	eg.Go(func() error {
		log.Println("[Notifier] Start.")
		select {
		case <-ctx.Done():
			log.Println("[Notifier] The context is done.")
			log.Println("[Notifier] End.")
			svr.GracefulStop()
			return nil
		}
	})

	eg.Go(func() error {
		log.Println("[Server] Start.")
		err := svr.Serve(lis)
		log.Println("[Server] End.")
		return err
	})
}
