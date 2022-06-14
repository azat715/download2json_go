# downloader2json

## Общий функционал

Загрузка фотографий с сайта https://jsonplaceholder.typicode.com

## Запуск

```bash
docker-compose up
```

## Основные технологии

Использован паттерн Worker pool

## Структура приложения

- internal/core - основная логика скрипта, враппер для errors
- internal/logger - логгер
- internal/models - модели для парсинга json
- internal/wpool - воркер пул

## TODO 

- добавить тесты
- урлы добавить в env
