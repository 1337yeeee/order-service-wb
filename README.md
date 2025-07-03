# Демонстрационный сервис с Kafka, PostgreSQL, кешем

Перед началом работы:
```bash
git clone https://github.com/1337yeeee/order-service-wb && cd order-service-wb
```

Создание файла с переменными окружения:
```bash
cp example.env .env
```

Настройте переменные окружения по Вашему усмотрению.

Контейнеризация и запуск:
```bash
docker-compose up -d --build
```

---
Приложение доступно по вдресу `http://localhost:8080`
