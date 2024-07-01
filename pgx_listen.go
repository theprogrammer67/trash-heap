package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
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
	listener.Handle("event", &Handler{})

	log.Println("listener started")
	listener.Listen(ctx)

	<-ctx.Done()
	log.Println("listener shutting down gracefully")

}

type Handler struct{}

func (h *Handler) HandleNotification(ctx context.Context, notification *pgconn.Notification, conn *pgx.Conn) error {
	var event Event

	err := json.Unmarshal([]byte(notification.Payload), &event)
	if err == nil {
		err = HandleEvent(event)
		if err == nil {
			if e := DeleteEvent(ctx, event.Id, conn); e != nil {
				err = fmt.Errorf("error deleting event: %v", e)
			}
		}
	}

	return err
}

func (h *Handler) HandleBacklog(ctx context.Context, channel string, conn *pgx.Conn) error {
	log.Println("HandleBacklog: " + channel)
	// read table events and delete it

	events, err := GetEvents(ctx, conn)
	if err == nil {
		for _, v := range events {
			e := HandleEvent(v)
			if e == nil {
				e = DeleteEvent(ctx, v.Id, conn)
			}

			if (e != nil) && (err == nil) {
				err = e
			}
		}
	}

	return err
}

func DeleteEvent(ctx context.Context, id uuid.UUID, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, "DELETE FROM events WHERE id = $1", id)
	return err
}

func HandleEvent(event Event) error {
	log.Println(event)
	return nil
}

func GetEvents(ctx context.Context, conn *pgx.Conn) (result []Event, err error) {
	err = pgxscan.Select(ctx, conn, &result, "SELECT id, payload FROM events ORDER BY created_at")
	return
}

type Event struct {
	Id      uuid.UUID `json:"id" db:"id"`
	Payload string    `json:"payload" db:"payload"`
}

type User struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
