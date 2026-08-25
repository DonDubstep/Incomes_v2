# Защита от нежелательного входа
Буду пробовать пробрасывать трафик через Tailscale. В случае неудачи - стандартная аутентификация

# Структура программы
Backend - Go. Забирает данные из бд и выдаёт фронтенду
Frontend - React. Визуализирует данные на графиках и циферках
SQLite - база

# БД
### period
id | period

### Incomes 
id | period_id | amount | Description | type_money_id | date

### Outcomes 
id | period_id | amount | category_id | Description | date

### Category
id | name

### Type of money
id | name|
1. Наличные
2. Безналичные

# Go функционал
### GET
- <b>GetIncomesByMonth</b> Возвращает сумму доходов за месяц (сумма)
- <b>GetAllIncomesRecordsByMonth</b> Возвращает все записи о доходах за месяц (Дата, сумма, описание, тип)
- <b>GetOutcomesByMonth</b> Возвращает сумму расходов за месяц (сумма)
- <b>GetAllOutcomesRecordsByMonth</b> Возвращает все записи расходов за месяц (дата, сумма, описание, категория)


### POST
- <b>SetIncome</b> Добавление дохода в базу (дата, месяц, сумма, описание, тип денег(нал/безнал))
- <b>SetOutcome</b> Добавление расхода(дата, месяц, сумма, описание, категория)

### DELETE
