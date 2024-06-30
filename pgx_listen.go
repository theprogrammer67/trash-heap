package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgxlisten"
)

func PgxListen() {
	connStr := "postgresql://proxy_server:123@localhost/test"
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	listener := pgxlisten.Listener{}
	listener.Connect = func(ctx context.Context) (*pgx.Conn, error) {
		return pgx.Connect(ctx, connStr)
	}
	listener.LogError = func(ctx context.Context, err error) {
		log.Println(err)
	}
	log.Println("listener started")

	listener.Handle("event", &Handler{})
	listener.Listen(ctx)

	<-ctx.Done()
	log.Println("listener shutting down gracefully")

}

type Handler struct{}

func (h *Handler) HandleNotification(ctx context.Context, notification *pgconn.Notification, conn *pgx.Conn) error {
	log.Println(notification.Payload)
	return nil
}

func (h *Handler) HandleBacklog(ctx context.Context, channel string, conn *pgx.Conn) error {
	log.Println(channel)
	// read table events and delete it

	return nil
}
