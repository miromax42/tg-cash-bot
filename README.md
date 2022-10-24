# Telegram bot (ДЗ1) 

[![pipeline status](https://gitlab.ozon.dev/miromaxxs/telegram-bot/badges/master/pipeline.svg)](https://gitlab.ozon.dev/miromaxxs/telegram-bot/-/commits/master) [![coverage report](https://gitlab.ozon.dev/miromaxxs/telegram-bot/badges/master/coverage.svg)](https://gitlab.ozon.dev/miromaxxs/telegram-bot/-/commits/master)

[[_TOC_]]

## ТЗ
### GOHW-1
- Команда добавления новой финансовой "траты". В трате должна присутствовать сумма, категория и дата. Но можете добавить еще поля, если считаете нужным. Придумайте, как оформить команду так, чтобы пользователю было удобно ее использовать.
- Хранение трат в памяти, базы данных пока не используем.
- Команда запроса отчета за последнюю неделю/месяц/год. В отчете должны быть суммы трат по категориям.

### GOHW-2
- Команда переключения бота на конкретную валюту - "/currency"
    1. После ввода команды бот предлагает выбрать интересующую валюту из четырех: USD, CNY, EUR, RUB
    2. При нажатии на нужную валюту переключаем бота на нее - результат получение трат конвертируется в выбранную валюту.
- Храним траты всегда в рублях, конвертацию используем только для отображения, ввода и отчетов
- Особенности
     * При запуске сервиса мы в отдельном потоке запрашиваем курсы валют.
     * Запрос курса валют происходит из любого из открытых источников.
     * Сервис должен завершаться gracefully.
  
### GOHW-3
- [x] Завести PostgreSQL в Docker, резметить таблички и схемы, любые действия с БД проводить только через миграции
  * в качестве базы теперь Postgres
  * миграции в `ent/migrate/migrations`
  * движок миграций `atlas`
- [x] Перенести хранение данных из памяти приложения в базу данных
- [x] Добавить бюджеты/лимиты: на траты в месяц, при проведении транзакций проверять согласованность данных (превышен ли бюджет и тд)
  > Добавлено поле пользовательских настроек Лимит на месяц, создан ендпоинт для его установки
- [x] Сгенерировать тестовые данные для таблицы расходов
  > Добавлена конфигурация `DB_TEST_USER_ID=<telegram_id>`, которая при установке не нулем/пустым, затирает информацию базы 
  > и для пользователя с данным telegram_id устанавливает тестовые параметры
- [x] Создать индексы на таблицу расходов, в комментариях к миграции пояснить выбор индексируемых колонок и типа индекса
  > Создан хэш-индекс на таблицу Расходов на поле создателя записи поскольку частая выборка по этому полю с условием равенства.
  > Поскольку для поля актуальна операция полного совпадения, хеш индекс очевидная замена б-три (быстрее, легче)
- [x] Доп: покрыть бизнес логику тестами, взаимодействие с базой замокать
  > Замоканы все интерфейсы проекта (`mockery`), табличные мок-тесты на хендлеры ()
- [x] Доп: интеграционные тесты на sql код
  > Интеграционные табл-тесты на постгрес в докере. Тест поднимает постгрес сам (`testcontainers`),
  > заполнение базы для теста (`testfixtures`)

### GOHW-5
- [x] Перевести бота на ведение структурированных логов в STDOUT.
- [x] Инструментировать код трейсингом. Создавать спан на каждое пришедшее сообщение. Корректно прокидывать контекст внутрь дерева функций и покрыть спанами важные части.
  > В качестве платформы трасировки Jaeger: [localhost](http://localhost:16686)
- [ ] Добавить метрики количества сообщений и времени обработки одного сообщения от пользователя. Разбить эти метрики по разным типам команд.
- [ ] Доп: Придумать и реализовать еще несколько полезных метрик для своего бота
- [ ] Доп: Добавить в свой композ-файл и настроить инфраструктуру сбора метрик и создать рабочий дашборд с несколькими панелями в Графане
- [ ] Доп: Добавить в композ-файл и настроить инфраструктуру сбора трейсов. Трейсы должны успешно искаться через веб-интерфейс Jaeger

## Граф зависимостей
![dependency graph](https://gitlab.ozon.dev/miromaxxs/telegram-bot/-/jobs/artifacts/master/raw/godepgraph.png?job=dependency-graph)

## Быстрый старт
0. Создать фаил `test.env` в корне проекта с содержимым:
  > `DB_TEST_USER_ID` можно посмотреть командой в боте `/ping`
  ```text
  TLG_TOKEN=5418307428:AAGaauK4oJwKCLKtwoxN7be8p2VKiVtBvus
  EXCHANGE_TOKEN=myjdYvKQjtY8oQ6LXtZhKYEuTP8t40o0
  EXCHANGE_BASE_CURRENCY=RUB
  DB_URL=postgres://postgres:pass@localhost:5432/test?sslmode=disable
  DB_TEST_USER_ID=
  ```
1. Установить taskfile
  ```bash 
  go install github.com/go-task/task/v3/cmd/task@latest
  ```
2. Запуск генераторов, контейнера с БД и приложения
  ```bash 
  task start-new
  ```
3. **Profit!**

## Конфигурация
* `TLG_TOKEN` - токен телеграм бота
* `EXCHANGE_TOKEN`- токен для [сервиса получения курсов валют](https://apilayer.com/marketplace/fixer-api)
* `EXCHANGE_BASE_CURRENCY` - код дефолтной валюты (USD, CNY, EUR, RUB)
* `DB_URL` - postgres URL (по-умолчанию `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`)
* `DB_TEST_USER_ID` - telegram ID юзера для которого создадутся тестовые данные (выставлять только в тестовой среде)

## Выбор библиотек
* `taskfile` - лучше чем Makefile во всем
* `telebot` - мне как бекенд разработчику, гораздо комфортнее с httpServer-подобной семантикой  
* `ent` - просто хороший ORM для golang (по сравнению с gorm, bun более производительный за счет кодогенерации)
* `atlas` - движок миграций поддерживающий декларативные миграции
* `testfixtures` - загрузка тестовых данных определенных в ямл
* `testcontainers` - запуск контейнеров из го (для интеграционного тестирования)
* `mockery` - кодогенерация моков, выбрал из-за более тесной интеграции с `testify`
* `cockroachdb/errors` - обработка ошибок аналогичная `pkg/errors`, но репа поддерживается и есть интеграция с `sentry`

## tAPI
### Получение id / alive signal
`/ping` - возвращает сигнал жизни с ID пользователя

---
### Добавление траты
`/exp <amount> <category>` - добавление траты в `<amount>` у.е. на категорию `<category>`. 
  #### Валидация
  * 0 < `amount` < 100000
  * `category` - любая строка длинны от 1 до 100
  #### Примеры
  * `/exp 12 FUN`
  * `/exp 11.99 rest`
  * `/exp 12.00 ZHkh`

---
### Отчет по тратам
`/all <time_token>` - траты за последнее время=`<time_token>` по кaтегориям
  #### Примеры
  * `/all day` - за последний день
  * `/all year` - за последний год
  * `/all 1m30s` - за последние полторы минуты

---
### Выбор валют
`/currency` - меню выбора валю из списка поддерживаемых (USD, CNY, EUR, RUB)

---
### Выставление лимита
`/limit <float>` - выставление лимита на месяц