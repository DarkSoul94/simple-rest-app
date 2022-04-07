package models

import "time"

type Book struct {
	ID           uint64
	Name         string
	Author       string
	CreationDate time.Time
}
