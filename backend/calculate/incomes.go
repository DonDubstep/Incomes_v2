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
			name TEXT NOT NULL UNIQUE,
			cource DOUBLE NOT NULL
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
			amount INTEGER NOT NULL,
			description TEXT,
			type_money_id INTEGER NOT NULL,
			date TEXT NOT NULL,
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
			(name, abbr, cource)
		VALUES 
			("Рубли","rub", 1.0),
			("Доллары","usd", 80.65),
			("Дирхамы","aed", 21.96);

			
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
			(amount, Description, type_money_id, date)
		VALUES 
			(20000, NULL, 1, '2026-05-29'),
			(30000, NULL, 1, '2026-06-29'),
			(40000, NULL, 1, '2026-07-29'),
			(10000, NULL, 2, '2026-08-10');

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

// Общая сумма за последний месяц
func GetAmountOfMoneyRub(db *sql.DB) (int, error) {
	query := `	WITH latest AS (
					SELECT mb.currency_id,
						mb.actual_sum,
						ROW_NUMBER() OVER (
							PARTITION BY mb.type_money_id, mb.currency_id
							ORDER BY p.name DESC
						) AS rn
					FROM monthly_balances mb
					JOIN period p ON p.id = mb.period_id
				)
				SELECT COALESCE(SUM(l.actual_sum * c.cource), 0) AS total_rub
				FROM latest l
				JOIN currency c ON c.id = l.currency_id
				WHERE l.rn = 1;`

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
	query := `	SELECT COALESCE(SUM(amount), 0) AS income_last_month
				FROM incomes
				WHERE strftime('%Y-%m', date) = (SELECT MAX(name) FROM period);`

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
	query := `	WITH ranked AS (
					SELECT name, ROW_NUMBER() OVER (ORDER BY name DESC) AS rn
					FROM period
				),
				targets AS (
					SELECT rn, name FROM ranked WHERE rn IN (1, 2)   -- 1 = последний, 2 = предыдущий
				),
				balances AS (
					SELECT t.rn AS target_rn,
						mb.currency_id,
						mb.actual_sum,
						ROW_NUMBER() OVER (
							PARTITION BY t.rn, mb.type_money_id, mb.currency_id
							ORDER BY p.name DESC
						) AS rk
					FROM targets t
					CROSS JOIN monthly_balances mb
					JOIN period p ON p.id = mb.period_id
					WHERE p.name <= t.name
				)
				SELECT COALESCE(SUM(CASE WHEN b.target_rn = 1 THEN b.actual_sum * c.cource END), 0)
					- COALESCE(SUM(CASE WHEN b.target_rn = 2 THEN b.actual_sum * c.cource END), 0)
					AS change_rub
				FROM balances b
				JOIN currency c ON c.id = b.currency_id
				WHERE b.rk = 1;`

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
