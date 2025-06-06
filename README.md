# quote-book
# Quotes API

Простой REST API-сервис на Go для хранения и управления цитатами.

## Функционал

- Добавление новой цитаты (POST /quotes)
- Получение всех цитат (GET /quotes)
- Получение случайной цитаты (GET /quotes?author=random)
- Фильтрация цитат по автору (GET /quotes?author=Имя)
- Удаление цитаты по ID (DELETE /quotes/{id})

## Запуск

