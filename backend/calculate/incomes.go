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

// Функция создания таблицы
func CreateTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS category (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS type_of_money (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS period (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS incomes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			period_id INTEGER NOT NULL,
			amount INTEGER NOT NULL,
			description TEXT,
			type_money_id INTEGER NOT NULL,
			date TEXT NOT NULL,
			FOREIGN KEY (period_id) REFERENCES period(id),
			FOREIGN KEY (type_money_id) REFERENCES type_of_money(id)
		);

		CREATE TABLE IF NOT EXISTS type_of_money (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
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
		INSERT INTO type_of_money 
			(name)
		VALUES 
			('Наличные'),
			('Безналичные');

		INSERT INTO period 
			(name)
		VALUES 
			('2026-05'),
			('2026-06'),
			('2026-07'),
			('2026-08');
			
		INSERT INTO incomes 
			(period_id, amount, Description, type_money_id, date)
		VALUES 
			(1, 20000, NULL, 1, '2026-05-29'),
			(2, 30000, NULL, 1, '2026-06-29'),
			(3, 40000, NULL, 1, '2026-07-29'),
			(3, 10000, NULL, 2, '2026-08-10'),
			(4, 50000, NULL, 2, '2026-09-14');
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
	TypeOfMoney string  `json:"type_of_money"`
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
				inner join type_of_money t on i.type_money_id = t.id
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
				inner join type_of_money t on i.type_money_id = t.id
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
