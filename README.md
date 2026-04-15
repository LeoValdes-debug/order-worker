# order-worker

Воркер для обработки событий заказов из RabbitMQ.

## Стек

- Go 1.22
- RabbitMQ 3.12
- Docker / Docker Compose

## Как работает

Воркер подписывается на очередь `orders` и обрабатывает два типа событий:

- `created` - отправляет пользователю подтверждение заказа
- `cancelled` - отправляет уведомление об отмене

При ошибке обработки сообщение возвращается в очередь (`nack + requeue`).

## Запуск

```bash
docker-compose up --build
```

RabbitMQ Management UI доступен на `http://localhost:15672` (guest / guest).

## Пример сообщения

```json
{
  "order_id": 123,
  "user_id": 456,
  "total_price": 1500.00,
  "status": "created",
  "created_at": "2026-04-15T10:00:00Z"
}
```

Отправить тестовое сообщение можно через Management UI: Queues -> orders -> Publish message.

## Структура

```
cmd/worker/         - точка входа
internal/
  consumer/         - подключение к RabbitMQ, чтение сообщений
  processor/        - логика обработки событий
  notifier/         - отправка уведомлений
```
