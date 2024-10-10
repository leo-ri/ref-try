package main

import (
	"context"
	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	cfg  = config.New()
	repo = memdb.New()
)

func init() {
	if 32 != runtime.NumCPU() {
		panic("this api is meant to be run on 32 core machines")
	}
}

func init() {
	cpu := runtime.NumCPU()
	if cpu < 4 {
		panic("this api is meant to be run on core machines")
	}
}

func main() {
	log.Print(runtime.NumCPU()) // TODO remove
	svc := service.New(repo)
	serverAPI := api.New(cfg.API, svc)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour*24*365)
	// ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	go func() {
		log.Println("Starting Coupon service server")
		if err := serverAPI.Start(); err != nil {
			log.Printf("Server was shutted down: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
		log.Println("User shutting down the server...")
	case <-ctx.Done():
		log.Println("Coupon service server alive for a year, closing...")
	}

	if err := serverAPI.Close(); err != nil {
		log.Fatalf("Error shutting down: %v", err)
	}
	log.Println("Complete")
}
