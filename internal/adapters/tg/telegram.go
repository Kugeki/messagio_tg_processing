package tg

import (
	"context"
	"github.com/nikoksr/notify/service/telegram"
	"log/slog"
	"messagio_tg_processing/internal/logger"
)

type Service struct {
	tg  *telegram.Telegram
	log *slog.Logger
}

func New(log *slog.Logger, apiToken string, chatIDs []int64) (*Service, error) {
	if log == nil {
		log = logger.NewEraseLogger()
	}
	telegramService, err := telegram.New(apiToken)
	if err != nil {
		return nil, err
	}

	telegramService.AddReceivers(chatIDs...)

	return &Service{tg: telegramService, log: log}, nil
}

func (s *Service) Send(ctx context.Context, subject string, message string) error {
	s.log.Info("sending message...")

	err := s.tg.Send(ctx, subject, message)
	if err != nil {
		s.log.Error("message send error", logger.Err(err))
		return err
	}

	s.log.Info("message is successfully sent")
	return nil
}
