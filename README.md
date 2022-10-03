# Telegram bot (ДЗ1) 

[[_TOC_]]

## ТЗ
### GOHW-1
- Команда добавления новой финансовой "траты". В трате должна присутствовать сумма, категория и дата. Но можете добавить еще поля, если считаете нужным. Придумайте, как оформить команду так, чтобы пользователю было удобно ее использовать.
- Хранение трат в памяти, базы данных пока не используем.
- Команда запроса отчета за последнюю неделю/месяц/год. В отчете должны быть суммы трат по категориям.

### GOHW-2
Новый функционал:
- Команда переключения бота на конкретную валюту - "выбрать валюту"
    1. После ввода команды бот предлагает выбрать интересующую валюту из четырех: USD, CNY, EUR, RUB
    2. При нажатии на нужную валюту переключаем бота на нее - результат получение трат конвертируется в выбранную валюту.
- Храним траты всегда в рублях, конвертацию используем только для отображения, ввода и отчетов
- Особенности
     * При запуске сервиса мы в отдельном потоке запрашиваем курсы валют.
     * Запрос курса валют происходит из любого из открытых источников.
     * Сервис должен завершаться gracefully.

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
