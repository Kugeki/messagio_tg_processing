package main

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	"messagio_tg_processing/internal/adapters/kafkaprod"
	"messagio_tg_processing/internal/adapters/tg"
	"messagio_tg_processing/internal/config"
	"messagio_tg_processing/internal/logger"
	"messagio_tg_processing/internal/ports/kafkacons"
	"messagio_tg_processing/internal/usecases"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err)
	}

	cfg, err := config.ReadEnvToConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	wg := sync.WaitGroup{}
	defer func() {
		slogger.Info("starting shuitdown")
		wg.Wait()
		slogger.Info("shutdown is completed")
	}()

	sarama.Logger = logger.NewSaramaLogger(slogger, slog.LevelDebug)
	sarama.DebugLogger = logger.NewSaramaLogger(slogger, slog.LevelDebug)

	saramaCfg := sarama.NewConfig()
	saramaCfg.ClientID = cfg.Kafka.ClientID
	saramaCfg.Version = sarama.V3_6_0_0

	slogger.Info("creating tg sender")
	tgSender, err := tg.New(slogger, cfg.Telegram.APIToken, []int64{cfg.Telegram.ChatID})
	if err != nil {
		slogger.Error("tg sender create", logger.Err(err))
		return
	}

	slogger.Info("creating msg producer")
	msgProducer, err := kafkaprod.NewProcMessagesProducer(slogger, cfg.Kafka.BrokerList, saramaCfg, cfg.Kafka.Producer.Topic)
	if err != nil {
		slogger.Error("msg producer create", logger.Err(err))
		return
	}
	defer func() {
		wg.Add(1)
		go func() {
			err = msgProducer.Close()
			if err != nil {
				slogger.Error("msg consumer close", logger.Err(err))
			}
			wg.Done()
		}()
	}()

	msgUC := usecases.NewMessageUC(tgSender, msgProducer)

	slogger.Info("creating msg consumer")
	msgConsumer, err := kafkacons.NewMessagesConsumer(slogger, msgUC,
		cfg.Kafka.BrokerList, saramaCfg, cfg.Kafka.Consumer.Group, cfg.Kafka.Consumer.Topics)
	if err != nil {
		slogger.Error("msg producer create", logger.Err(err))
		return
	}
	defer func() {
		wg.Add(1)
		go func() {
			err = msgConsumer.Close()
			if err != nil {
				slogger.Error("msg consumer close", logger.Err(err))
			}
			wg.Done()
		}()
	}()

	slogger.Info("start consuming")
	go func() {
		msgConsumer.StartConsume(ctx)
	}()

	<-ctx.Done()
}
