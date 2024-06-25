package main

import (
	"log"
	"time"
)

func TestZeroDate() {
	log.Println("zero date: ", ZeroDate(time.Now()))
	log.Println("zero time: ", ZeroTime(time.Now()))

	t, _ := time.Parse("15:04:05", time.Now().Format("15:04:05"))
	log.Println("zero date 1: ", t)

	t, _ = time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	log.Println("zero time 1: ", t)
}

func ZeroDate(t time.Time) time.Time {
	h, m, s := t.Clock()
	return time.Date(1, 1, 1, h, m, s, 0, t.Location())
}

func ZeroTime(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
