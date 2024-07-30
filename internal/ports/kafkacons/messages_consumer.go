package kafkacons

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"log/slog"
	"messagio_tg_processing/internal/domain"
	"messagio_tg_processing/internal/logger"
)

type MessagesUsecase interface {
	SendMessage(ctx context.Context, msg *domain.Message) error
	ProduceMessage(msg *domain.Message)
}

type MessagesConsumer struct {
	log *slog.Logger

	msgUC  MessagesUsecase
	cg     sarama.ConsumerGroup
	topics []string
}

func NewMessagesConsumer(log *slog.Logger, msgUC MessagesUsecase, brokerList []string,
	saramaCfg *sarama.Config, group string, topics []string) (*MessagesConsumer, error) {
	if log == nil {
		log = logger.NewEraseLogger()
	}
	log = log.With(slog.String("component", "ports/kafkacons/messages_consumer"))

	if saramaCfg == nil {
		saramaCfg = sarama.NewConfig()
	}

	conf := *saramaCfg

	consumerGroup, err := sarama.NewConsumerGroup(brokerList, group, &conf)
	if err != nil {
		return nil, err
	}

	return &MessagesConsumer{log: log, msgUC: msgUC, cg: consumerGroup, topics: topics}, nil
}

func (c *MessagesConsumer) Close() error {
	if c.cg != nil {
		return c.cg.Close()
	}
	return errors.New("MessagesConsumer.Close: cg is nil")
}

// StartConsume is blocking
func (c *MessagesConsumer) StartConsume(ctx context.Context) {
	for {
		if err := c.cg.Consume(ctx, c.topics, c); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return
			}
			c.log.Error("error from consume",
				slog.String("function", "StartConsume"),
				logger.Err(err),
			)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func (c *MessagesConsumer) ConsumeClaim(ses sarama.ConsumerGroupSession, cm sarama.ConsumerGroupClaim) error {
	log := c.log.With(
		slog.String("handler", "processed message"),
		slog.String("member_id", ses.MemberID()),
		slog.Int64("generation_id", int64(ses.GenerationID())),
	)

	log.Info("consume claim")
	for {
		select {
		case claimMsg, ok := <-cm.Messages():
			if !ok {
				log.Info("claim message channel was closed")
				return nil
			}

			log = c.log.With(
				slog.String("handler", "processed message"),
				slog.String("topic", claimMsg.Topic),
				slog.Int64("partition", int64(claimMsg.Partition)),
				slog.Int64("offset", claimMsg.Offset),
				slog.Time("timestamp", claimMsg.Timestamp),
				slog.Time("block_timestamp", claimMsg.BlockTimestamp),
			)
			err := c.HandleMessage(ses.Context(), log, claimMsg)
			if err != nil {
				log.Error("handle claim message", logger.Err(err))
			}

			ses.MarkMessage(claimMsg, "")
		case <-ses.Context().Done():
			ses.Commit()
			return nil
		}
	}
}

func (c *MessagesConsumer) HandleMessage(ctx context.Context, log *slog.Logger, claimMsg *sarama.ConsumerMessage) error {
	log.Info("starting handle message")

	msg := domain.Message{}
	err := json.Unmarshal(claimMsg.Value, &msg)
	if err != nil {
		log.Warn("json unmarshal", logger.Err(err))
		return err
	}

	err = c.msgUC.SendMessage(ctx, &msg)
	if err != nil {
		log.Error("send message", logger.Err(err))
		return err
	}

	log.Info("message is handled")

	c.msgUC.ProduceMessage(&msg)

	log.Info("message is produced")

	return nil
}

func (c *MessagesConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *MessagesConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}
