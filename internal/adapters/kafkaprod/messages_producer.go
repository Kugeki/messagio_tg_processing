package kafkaprod

import (
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"log/slog"
	"messagio_tg_processing/internal/domain"
	"messagio_tg_processing/internal/logger"
)

type ProcMessagesProducer struct {
	p     sarama.AsyncProducer
	log   *slog.Logger
	topic string
}

func NewProcMessagesProducer(log *slog.Logger, brokerList []string,
	saramaCfg *sarama.Config, topic string) (*ProcMessagesProducer, error) {
	if log == nil {
		log = logger.NewEraseLogger()
	}
	log = log.With(slog.String("component", "adapters/kafkaprod/proc_messages_producer"))

	if saramaCfg == nil {
		saramaCfg = sarama.NewConfig()
	}

	conf := *saramaCfg

	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Compression = sarama.CompressionSnappy

	producer, err := sarama.NewAsyncProducer(brokerList, &conf)
	if err != nil {
		return nil, err
	}

	mp := &ProcMessagesProducer{p: producer, log: log, topic: topic}

	go func() {
		for err = range producer.Errors() {
			mp.log.Error("messages producer error", logger.Err(err))
		}
	}()

	return mp, nil
}

func (p *ProcMessagesProducer) Close() error {
	if p.p != nil {
		return p.p.Close()
	}
	return errors.New("MessagesProducer.Close: async producer is nil")
}

func (p *ProcMessagesProducer) Produce(msg *domain.Message) {
	p.p.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   nil,
		Value: newMessageValue(msg),
	}
}

type MessageValue struct {
	ID int `json:"ID"`

	bytes []byte
	err   error
}

func newMessageValue(msg *domain.Message) *MessageValue {
	v := &MessageValue{ID: msg.ID}

	var err error
	v.bytes, err = json.Marshal(v)
	if err != nil {
		v.err = err
	}
	return v
}

func (v *MessageValue) Encode() ([]byte, error) {
	if v.err != nil {
		return v.bytes, v.err
	}
	return v.bytes, nil
}

func (v *MessageValue) Length() int {
	return len(v.bytes)
}
