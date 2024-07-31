package tg

import (
	"context"
	"github.com/nikoksr/notify/service/telegram"
	"golang.org/x/time/rate"
	"log/slog"
	"messagio_tg_processing/internal/logger"
)

type Service struct {
	tg  *telegram.Telegram
	log *slog.Logger

	limitter *rate.Limiter
}

var MaxBursts = 1

func New(log *slog.Logger, apiToken string, chatIDs []int64, limitPerSecond float64) (*Service, error) {
	if log == nil {
		log = logger.NewEraseLogger()
	}
	telegramService, err := telegram.New(apiToken)
	if err != nil {
		return nil, err
	}

	telegramService.AddReceivers(chatIDs...)

	log.Info("telegram new", slog.Float64("limit per second", limitPerSecond))
	limitter := rate.NewLimiter(rate.Limit(limitPerSecond), MaxBursts)

	return &Service{tg: telegramService, log: log, limitter: limitter}, nil
}

func (s *Service) Send(ctx context.Context, subject string, message string) error {
	s.log.Info("limitter waiting...", slog.Float64("tokens", s.limitter.Tokens()))
	err := s.limitter.Wait(ctx)
	if err != nil {
		s.log.Error("limitter wait error", logger.Err(err))
		return err
	}

	s.log.Info("sending message...")
	err = s.tg.Send(ctx, subject, message)
	if err != nil {
		s.log.Error("message send error", logger.Err(err))
		return err
	}

	s.log.Info("message is successfully sent")
	return nil
}
