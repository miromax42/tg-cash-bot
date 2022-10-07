# Telegram bot (ДЗ1) 

[![pipeline status](https://gitlab.ozon.dev/miromaxxs/telegram-bot/badges/master/pipeline.svg)](https://gitlab.ozon.dev/miromaxxs/telegram-bot/-/commits/master) [![coverage report](https://gitlab.ozon.dev/miromaxxs/telegram-bot/badges/master/coverage.svg)](https://gitlab.ozon.dev/miromaxxs/telegram-bot/-/commits/master)

[[_TOC_]]

## ТЗ
- Команда добавления новой финансовой "траты". В трате должна присутствовать сумма, категория и дата. Но можете добавить еще поля, если считаете нужным. Придумайте, как оформить команду так, чтобы пользователю было удобно ее использовать.
- Хранение трат в памяти, базы данных пока не используем.
- Команда запроса отчета за последнюю неделю/месяц/год. В отчете должны быть суммы трат по категориям.

## Конфигурация
* `TLG_TOKEN` - токен телеграм бота

## Выбор библиотек
* `telebot` - мне как бекенд разработчику, гораздо комфортнее с httpServer-подобной семантикой  
* `ent` - просто хороший ORM для golang (по сравнению с gorm, bun более производительный за счет кодогенерации)

## tAPI

### Добавление траты
`/exp <amount> <category>` - добавление траты в `<amount>` у.е. на категорию `<category>`. 
  #### Валидация
  * 0 < `amount` < 100000
  * `category` - любая строка длинны от 1 до 100
  #### Примеры
  * `/exp 12 FUN`
  * `/exp 11000 rest`
  * `/exp 9 ZHkh`
---
### Отчет по тратам
`/all <time_token>` - траты за последнее время=`<time_token>` по кaтегориям
  #### Примеры
  * `/all day` - за последний день
  * `/all year` - за последний год
  * `/all 1m30s` - за последние полторы минуты
