# Docker Deployment

Этот каталог содержит файлы для развертывания парсера Хабра с помощью Docker.

## Файлы

- `Dockerfile` - Docker файл для сборки Go приложения
- `docker-compose.yml` - Конфигурация для запуска сервиса
- `README.md` - Инструкции по развертыванию

## Запуск

1. Убедитесь, что у вас установлен Docker и Docker Compose
2. Перейдите в корневую директорию проекта
3. Выполните команду:

```bash
docker-compose -f scripts/docker-compose.yml up --build
```

## Остановка

```bash
docker-compose -f scripts/docker-compose.yml down
```

## Переменные окружения

Сервис использует следующие переменные окружения (шаблон в `.env.example`):

- `BASE_URL` - Базовый URL Хабра
- `POST_URL` - URL для постов
- `NEWS_URL` - URL для новостей
- `ARTICLE_URL` - URL для статей
- `SEARCH_URL` - URL для поиска

Для локальной разработки:
```bash
cp .env.example .env
# Отредактируйте .env файл при необходимости
```

## Порты

Сервис доступен на порту 80 внутри контейнера и пробрасывается на порт 80 хоста.