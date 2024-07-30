package usecases

import (
	"context"
	"messagio_tg_processing/internal/domain"
)

type MessageUC struct {
	sender   domain.Sender
	producer domain.Producer
}

func NewMessageUC(sender domain.Sender, producer domain.Producer) *MessageUC {
	return &MessageUC{sender: sender, producer: producer}
}

func (uc *MessageUC) SendMessage(ctx context.Context, msg *domain.Message) error {
	return uc.sender.Send(ctx, "", msg.Content)
}

func (uc *MessageUC) ProduceMessage(msg *domain.Message) {
	uc.producer.Produce(msg)
}
