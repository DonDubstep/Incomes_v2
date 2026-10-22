package calculate

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// Подключение к БД
func ConnectToDataBase(sql_path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", sql_path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Удаление всех таблиц в БД
func ResetDatabase(db *sql.DB) error {
	rows, err := db.Query(`
		SELECT name FROM sqlite_master
		WHERE type='table' AND name NOT LIKE 'sqlite_%'
	`)
	if err != nil {
		return fmt.Errorf("query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("scan table name: %w", err)
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate tables: %w", err)
	}

	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, table)); err != nil {
			return fmt.Errorf("drop table %s: %w", table, err)
		}
	}
	log.Print("wiped")
	return nil
}

// Функция создания таблицы
func CreateTable(db *sql.DB) {
	_, err := db.Exec(`

		CREATE TABLE IF NOT EXISTS period (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS type_money (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS currency (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			abbr TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS category (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS monthly_balances (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			period_id INTEGER NOT NULL,
			type_money_id INTEGER NOT NULL,
			currency_id INTEGER NOT NULL,
			actual_sum INTEGER,
			FOREIGN KEY (period_id) REFERENCES period(id),
			FOREIGN KEY (type_money_id) REFERENCES type_money(id),
			FOREIGN KEY (currency_id) REFERENCES currency(id)
		);

		CREATE TABLE IF NOT EXISTS incomes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			period_id INTEGER NOT NULL,
			amount INTEGER NOT NULL,
			description TEXT,
			type_money_id INTEGER NOT NULL,
			date TEXT NOT NULL,
			FOREIGN KEY (period_id) REFERENCES period(id),
			FOREIGN KEY (type_money_id) REFERENCES type_money(id)
		);

		CREATE TABLE IF NOT EXISTS outcomes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			period_id INTEGER NOT NULL,
			amount INTEGER NOT NULL,
			description TEXT,
			date TEXT NOT NULL,
			FOREIGN KEY (period_id) REFERENCES period(id)
		);
	
	`)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("created")
}

// Функция заполнения таблицы
func FillTable(db *sql.DB) {
	_, err := db.Exec(`

		INSERT INTO period 
			(name)
		VALUES 
			('2026-05'),
			('2026-06'),
			('2026-07'),
			('2026-08');

		INSERT INTO type_money 
			(name)
		VALUES 
			('Наличные'),
			('Безналичные');

		INSERT INTO currency 
			(name, abbr)
		VALUES 
			("Рубли","rub"),
			("Доллары","usd"),
			("Дирхамы","aed");

			
		INSERT INTO monthly_balances 
			(period_id, type_money_id, currency_id, actual_sum)
		VALUES 
			(1, 2, 1, 100000),
			(1, 1, 1, 50000),
			(1, 1, 2, 200),
			(1, 1, 3, 500),
			(2, 1, 1, 80000),
			(3, 1, 1, 100000),
			(4, 1, 1, 130000),
			(5, 1, 1, 150000),
			(5, 2, 1, 110000);

		INSERT INTO incomes 
			(period_id, amount, Description, type_money_id, date)
		VALUES 
			(1, 20000, NULL, 1, '2026-05-29'),
			(2, 30000, NULL, 1, '2026-06-29'),
			(3, 40000, NULL, 1, '2026-07-29'),
			(3, 10000, NULL, 2, '2026-08-10');

	`)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("filled")
}

type IncomesByPeriod struct {
	Period string  `json:"period"`
	Total  float64 `json:"total"`
}

type IncomesRecords struct {
	Date        string  `json:"date"`
	Total       float64 `json:"total"`
	Description string  `json:"description"`
	TypeOfMoney string  `json:"type_money"`
}

// Выводит все доходы из БД
func GetAllIncomes(db *sql.DB) ([]IncomesByPeriod, error) {
	query := `	SELECT p.name, COALESCE(SUM(i.amount), 0) 
				FROM period p 
				LEFT JOIN incomes i 
				ON i.period_id = p.id
				GROUP BY p.id, p.name
				ORDER BY p.id;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	var result []IncomesByPeriod

	for rows.Next() {
		var r IncomesByPeriod
		if err := rows.Scan(&r.Period, &r.Total); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %w", err)
		}
		result = append(result, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обходе результата: %w", err)
	}

	return result, nil
}

// Выводит доход за выбранный месяц
func GetIncomesByMonth(db *sql.DB, req_period string) (int, error) {
	query := `	select sum(i.amount) 
				from incomes i 
				inner join period p on i.period_id = p.id
				where p.name = ?
				`

	var result int
	err := db.QueryRow(query, req_period).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("запрос вернул ничего GetIncomesByMonth %w", err) // Возвращаем 0 и nil (ошибки нет, просто ничего не нашли)
		}
		return 0, fmt.Errorf("ошибка запроса GetIncomesByMonth %w", err)
	}
	return result, nil

}

// Выводит все записи доходов за выбраный месяц
func GetAllIncomesRecordsByMonth(db *sql.DB, req_period string) ([]IncomesRecords, error) {
	query := `	select i.date, i.amount, COALESCE(i.description, ''), t.name 
				from incomes i
				inner join period p on i.period_id = p.id
				inner join type_money t on i.type_money_id = t.id
				where p.name = ?`

	rows, err := db.Query(query, req_period)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	var result []IncomesRecords
	for rows.Next() {
		var el IncomesRecords
		if err := rows.Scan(&el.Date, &el.Total, &el.Description, &el.TypeOfMoney); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки %w", err)
		}
		result = append(result, el)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обходе результата: %w", err)
	}
	return result, nil
}

// Выводит все записи за выбранный диапазон
func GetIncomesRecordsByRange(db *sql.DB, req_period_start string, req_period_end string) ([]IncomesRecords, error) {
	query := `	select i.date, i.amount, COALESCE(i.description, ''), t.name 
				from incomes i
				inner join period p on i.period_id = p.id
				inner join type_money t on i.type_money_id = t.id
				where p.name >= ? and p.name <= ?`

	rows, err := db.Query(query, req_period_start, req_period_end)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	var result []IncomesRecords
	for rows.Next() {
		var el IncomesRecords
		if err := rows.Scan(&el.Date, &el.Total, &el.Description, &el.TypeOfMoney); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки %w", err)
		}
		result = append(result, el)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обходе результата: %w", err)
	}
	return result, nil
}

func GetAmountOfMoneyRub(db *sql.DB) (int, error) {
	query := `	SELECT SUM(m.actual_sum)
				FROM monthly_balances m
				WHERE m.period_id = (
					select m.period_id
					from monthly_balances m
					inner join period p on p.id = m.period_id
					order by p.name desc limit 1
				);`

	var result int
	err := db.QueryRow(query).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("запрос вернул ничего GetAmountOfMoneyRub %w", err) // Возвращаем 0 и nil (ошибки нет, просто ничего не нашли)
		}
		return 0, fmt.Errorf("ошибка запроса GetAmountOfMoneyRub %w", err)
	}
	return result, nil
}

// Возвращает доход за последний месяц
func GetAmountOfLastIncome(db *sql.DB) (int, error) {
	query := `SELECT SUM(amount)
			  FROM incomes 
			  where period_id = (
				select i.period_id
				from incomes i
				inner join period p on p.id = i.period_id
				order by p.name desc limit 1
				);`

	var result int
	err := db.QueryRow(query).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("запрос вернул ничего GetAmountOfLastIncome %w", err) // Возвращаем 0 и nil (ошибки нет, просто ничего не нашли)
		}
		return 0, fmt.Errorf("ошибка запроса GetAmountOfLastIncome %w", err)
	}
	return result, nil
}

// Возвращает прирост/уменьшение общей суммы относительно предыдущего периода
func GetTotalLastChange(db *sql.DB) (int, error) {
	query := `	WITH totals AS (
				SELECT
					mb.period_id,
					p.name AS period_name,
					SUM(mb.actual_sum) AS total_sum
				FROM monthly_balances mb
				JOIN period p ON p.id = mb.period_id
				GROUP BY mb.period_id, p.name
			),
			ranked AS (
				SELECT
					period_id,
					period_name,
					total_sum,
					LAG(total_sum) OVER (ORDER BY period_name) AS prev_total
				FROM totals
			)
			SELECT
				total_sum - prev_total   AS diff
			FROM ranked
			ORDER BY period_name DESC
			LIMIT 1;
			;`

	var result int
	err := db.QueryRow(query).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("запрос вернул ничего GetAmountOfLastIncome %w", err) // Возвращаем 0 и nil (ошибки нет, просто ничего не нашли)
		}
		return 0, fmt.Errorf("ошибка запроса GetAmountOfLastIncome %w", err)
	}
	return result, nil
}

// Возвращает сумму денег в каждой валюте в виде таблицы
