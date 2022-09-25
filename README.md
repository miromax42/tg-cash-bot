# Telegram bot (ДЗ1) 
## ТЗ
- Команда добавления новой финансовой "траты". В трате должна присутствовать сумма, категория и дата. Но можете добавить еще поля, если считаете нужным. Придумайте, как оформить команду так, чтобы пользователю было удобно ее использовать.
- Хранение трат в памяти, базы данных пока не используем.
- Команда запроса отчета за последнюю неделю/месяц/год. В отчете должны быть суммы трат по категориям.
## Конфигурация
Приложение конфигурируется из переменных окружения:
* `TLG_TOKEN` - токен телеграм бота
## Выбор библиотек
* `telebot` - мне как бекенд разработчику, гораздо комфортнее с httpServer-подобной семантикой  
* `ent` - просто хороший ORM для golang (по сравнению с gorm, bun более производительный за счет кодогенерации)
## tAPI
1. `/exp <amount> <category>` - добавление траты в `<amount>` у.е. на категорию `<category>`.

    **Валидация:**
    * 0 < `amount` < 100000
    * `category` - любая строка длинны от 1 до 100
   
    **Примеры:**
    * `/exp 12 FUN`
    * `/exp 11000 rest`
    * `/exp 9 ZHkh`
2. `/all <time_token>` - траты за последнее время=`<time_token>` по кaтегориям

    **Валидация:**
    * `time_token` - одно из значений (`day`,`week`,`month`.`year`)  или в [формате duration](https://pkg.go.dev/time#ParseDuration)
    * `category` - любая строка длинны от 1 до 100

   **Примеры:**
    * `/all day` - за последний день
    * `/all year` - за последний год
    * `/all 1m30s` - за последние полторы минуты