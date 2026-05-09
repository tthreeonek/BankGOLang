# Bank API – Банковский REST API на Go

Выполнен как учебный проект. Реализована регистрация/аутентификация пользователей, управление банковскими счетами и виртуальными картами, переводы, кредитование (аннуитетные платежи), аналитика, интеграция с ЦБ РФ (SOAP) и SMTP-уведомления.

## 🔧 Технологии

- **Язык:** Go 1.23+
- **База данных:** PostgreSQL 18
- **Маршрутизация:** gorilla/mux
- **Аутентификация:** JWT (golang-jwt/jwt/v5)
- **Логирование:** logrus
- **Пароли:** bcrypt
- **Шифрование карт:** PGP (openpgp) + HMAC-SHA256
- **CVV:** bcrypt
- **SMTP:** gomail.v2
- **SOAP (ЦБ РФ):** beevik/etree
- **Планировщик платежей:** время через time.Ticker

## 📁 Структура проекта
```
bank-api/
├── cmd/server/ # Входная точка, DI
├── internal/
│ ├── config/ # Загрузка .env
│ ├── model/ # Структуры данных с валидацией
│ ├── repository/ # Интерфейсы и реализация PostgreSQL
│ ├── service/ # Бизнес-логика, интеграции
│ ├── handler/ # HTTP-обработчики
│ ├── middleware/ # JWT-проверка
│ └── router/ # Регистрация маршрутов
├── pkg/ # Вспомогательные пакеты (luhn, crypto)
├── migrations/ # SQL-миграции
├── .env.example # Пример переменных окружения
├── go.mod / go.sum
└── README.md
```
## ⚙️ Установка и запуск

1. **Клонируйте репозиторий**
   ```bash
   git clone https://github.com/tthreeonek/BankGOLang.git
   cd BankGOLang

2. **Настройте переменные окружения**

Создайте в корне проекта файл ```.env``` и замените данные в файле ```.env``` реальными данными:

```text
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=yourpassword
    DB_NAME=bankdb
    SSL_MODE=disable
    JWT_SECRET=your-secret-key
    AES_KEY=32-символьная-строка-для-шифрования
    HMAC_SECRET=секрет-для-hmac
    SMTP_HOST=smtp.example.com
    SMTP_PORT=587
    SMTP_USER=noreply@example.com
    SMTP_PASS=your-smtp-password
```
3. **Создайте базу данных и примените миграции**

```bash
    psql -U postgres -c "CREATE DATABASE bankdb;"
    psql -U postgres -d bankdb -f migrations/001_init.sql
```
4. **Установите зависимости и запустите сервер**

```bash
go mod tidy
go run cmd/server/main.go
```
Сервер стартует на порту ```8080.```

## 🔐 Аутентификация
Все защищённые эндпоинты требуют заголовок:

```text
Authorization: Bearer <jwt_token>
```
Токен получается при логине и действителен 24 часа.

## 🌐 API Endpoints
**Публичные**


```POST /register```	Регистрация	```{"username":"...","email":"...","password":"..."}```

```POST /login```	Вход (возвращает JWT)	```{"email":"...","password":"..."}```

**Защищённые (префикс /api)**

```POST /api/accounts```	Создать счёт

```GET	/api/accounts```	Список счетов пользователя

```POST /api/accounts/{id}/deposit```	Пополнить счёт (сумма в body)

```GET	/api/accounts/{id}/predict?days=N```	Прогноз баланса на N дней (макс. 365)

```POST /api/cards```	Выпустить виртуальную карту ({"account_id":1})

```POST /api/transfer```	Перевод ```({"from_account_id":1,"to_account_id":2,"amount":100})```

```POST /api/credits```	Оформление кредита ```({"account_id":1,"amount":100000,"term_months":12})```

```GET	/api/credits/{creditId}/schedule```	График платежей по кредиту

```GET	/api/analytics/monthly?month=2026-05```	Статистика доходов/расходов

```GET	/api/analytics/credit-load```	Текущая кредитная нагрузка

## 🧪 Примеры тестовых запросов (curl)

**1.** **Регистрация**
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"secret123"}'
```
**2.** **Логин**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secret123"}'
```
Сохраните полученный токен в переменную $TOKEN

**3.** **Создание счёта**
```bash
curl -X POST http://localhost:8080/api/accounts \
  -H "Authorization: Bearer $TOKEN"
```
**4.** **Выпуск карты**
```bash
curl -X POST http://localhost:8080/api/cards \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"account_id": 1}'
```
**5.** **Пополнение**
```bash
curl -X POST http://localhost:8080/api/accounts/1/deposit \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"amount": 5000}'
```
**6.** **Перевод**
```bash
curl -X POST http://localhost:8080/api/transfer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"from_account_id":1,"to_account_id":2,"amount":1000}'
```
**7.** **Оформление кредита**
```bash
curl -X POST http://localhost:8080/api/credits \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"account_id":1,"amount":100000,"term_months":12}'
```
**8.** **График платежей**
```bash
curl http://localhost:8080/api/credits/1/schedule \
  -H "Authorization: Bearer $TOKEN"
```
**9.** **Аналитика**
```bash
curl "http://localhost:8080/api/analytics/monthly?month=2025-05" \
  -H "Authorization: Bearer $TOKEN"
```
**10.** **Прогноз баланса**
```bash
curl "http://localhost:8080/api/accounts/1/predict?days=60" \
  -H "Authorization: Bearer $TOKEN"

```
## 🛡️ Безопасность
**Пароли** – хешируются bcrypt (12 раундов).

**Данные карт** – номер и срок шифруются PGP (алгоритм openpgp), целостность номера проверяется **HMAC-SHA256.**

**CVV** – хешируется bcrypt (не требует расшифровки).

**Авторизация** – проверка JWT на каждом защищённом маршруте; дополнительно проверяется принадлежность счёта/карты текущему пользователю.

**Все SQL-запросы** параметризованы, транзакции для переводов и создания графиков платежей.

## 📅 Планировщик кредитных платежей
Фоновая горутина запускается каждые 12 часов. Логика:

Находит просроченные платежи (```status='pending' и due_date <= NOW()```).

Если на счёте достаточно средств – автоматическое списание, статус ```paid```, отправляется email-уведомление.

Если средств недостаточно – начисляется штраф 10% к сумме платежа, новая дата через месяц.

## 📈 Интеграция с ЦБ РФ
Сервис ```CBRService``` отправляет SOAP-запрос к ```DailyInfoWebServ/DailyInfo.asmx```. Полученная ключевая ставка увеличивается на 5% (маржа банка) и используется при расчёте аннуитетного платежа.

## 📩 SMTP-уведомления
При успешном автоматическом списании кредитного платежа отправляется письмо на email пользователя. Настройки SMTP задаются в ```.env```.

## 📌 Дальнейшее развитие
Миграция на актуальную PGP-библиотеку (openpgp устарел)

Добавление 2FA

Административная панель

Полное покрытие unit-тестами
