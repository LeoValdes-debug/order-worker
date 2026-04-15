package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/leovaldes-debug/order-worker/internal/consumer"
	"github.com/leovaldes-debug/order-worker/internal/notifier"
	"github.com/leovaldes-debug/order-worker/internal/processor"
)

func main() {
	_ = godotenv.Load()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL is required")
	}

	queue := os.Getenv("QUEUE_NAME")
	if queue == "" {
		queue = "orders"
	}

	c, err := consumer.New(rabbitURL, queue)
	if err != nil {
		log.Fatalf("init consumer: %v", err)
	}
	defer c.Close()

	n := notifier.New()
	p := processor.New(n)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := c.Run(ctx, p.Handle); err != nil {
		log.Printf("consumer stopped: %v", err)
	}

	log.Println("worker stopped")
}
