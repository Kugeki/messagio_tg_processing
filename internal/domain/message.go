package domain

import "context"

type Message struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	Processed bool   `json:"processed"`
}

type Sender interface {
	Send(ctx context.Context, subject string, message string) error
}

type Producer interface {
	Produce(msg *Message)
}
