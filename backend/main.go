package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// Функция создания таблицы
func createTable(db *sql.DB) {
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
func fillTable(db *sql.DB) {
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
			(0, 20000, NULL, 1, '2026-05-29'),
			(1, 30000, NULL, 1, '2026-06-29'),
			(2, 40000, NULL, 1, '2026-07-29'),
			(2, 10000, NULL, 2, '2026-08-10'),
			(3, 50000, NULL, 2, '2026-09-14');
	`)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("filled")
}

// Подключение к БД
func connectToDataBase(sql_path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", sql_path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type IncomesByPeriod struct {
	Period string
	Total  float64
}

// Выводит все доходы из БД
func GetAllIncomes(db *sql.DB) ([]IncomesByPeriod, error) {
	query := `SELECT period_id, SUM(amount) AS total FROM Incomes GROUP BY period_id`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	var result []IncomesByPeriod

	for rows.Next() {
		var r IncomesByPeriod
		if err := rows.Scan(&r.Period, &r.Total); err != nil {
			return nil, fmt.Errorf("ощибка чтения строки: %w", err)
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

type IncomesRecords struct {
	Date        string
	Total       float64
	Description string
	Type        string
}

// Выводит все записи доходов за выбраный месяц
func GetAllIncomesRecordsByMonth(db *sql.DB, req_period string) ([]IncomesRecords, error) {
	// TODO: implement
}

// Выводит все записи за выбранный диапазон
func GetIncomesRecordsByRange(db *sql.DB, req_period_start string, req_period_end string) ([]IncomesRecords, error) {
	// TODO: implement
}

func main() {
	sql_path := "storage.db"
	db, err := connectToDataBase(sql_path)
	if err != nil {
		fmt.Println("ошибка открытия БД %w", err)
		return
	}
	defer db.Close()
	createTable(db)
	fillTable(db)

	// var result []IncomesByPeriod
	// result, err = GetAllIncomes(db)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// for _, el := range result {
	// 	fmt.Printf("%s\t%f\n", el.Period, el.Total)
	// }

	result, err := GetIncomesByMonth(db, "2026-06")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("result = ", result)

}
