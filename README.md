# Уровень L0

## Задание - [**ТЫК**](https://github.com/MetaException/wb-internship_l0/blob/master/Task.md)

## Настройка

Файлы конфигурации приложения config.go (internal/*/config.go) для настройки параметров по умолчанию.

Можно использовать переменные среды:
```
APISERVER_BINDADDR - адрес, на котором API сервер будет принимать входящие соединения.

NATS_URL - строка подключения к NATS
NATS_STREAMNAME - название потока (stream) в NATS, к которому нужно подписаться
NATS_STREAM_SUBJECT - subject (тема) в NATS, на который будет подписываться клиент.

POSTGRESQL_URL - строка подключения к бд postgresql
```
## Использование

Обращение к API: http://localhost:8080/api/orders?order_uid=b563feb7b2b84b6test

Скрипт fill.py - скрпит для публикации сгенерированных данных в поток.
```
python fill.py
```
