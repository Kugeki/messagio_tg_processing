# messagio_tg_processing
Для [тестового задания в Messagio](https://github.com/Kugeki/messagio_assignment).

Микросервис, который читает сообщения из Kafka
и отправляет их в телеграм-канал.

Сделан для демонстрации работоспособности тестового задания.
Качеству кода времени почти не уделялось.

> [!CAUTION]
> Рекомендуется выставлять TELEGRAM_LIMIT_PER_SECOND не
> менее 1 из-за ограничений телеграмма на отправку сообщений
> в минуту во избежание превышения rate лимита.

## Запуск
```
docker compose -d --build
```

## Пример конфигурации
.env файл:
```
TELEGRAM_API_TOKEN=api_token
TELEGRAM_CHAT_ID=chat_id
TELEGRAM_LIMIT_PER_SECOND=0.5

KAFKA_CLIENT_ID=messagio-tg-processing
KAFKA_BROKER_LIST=localhost:9092

KAFKA_CONSUMER_GROUP=messagio-tg-processing
KAFKA_CONSUMER_TOPICS=messages-to-process

KAFKA_PRODUCER_TOPIC=processed-messages
```