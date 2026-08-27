package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

// Функция создания таблицы
func createTable() {
	db, err := sql.Open("sqlite", "storage.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
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

func fillTable() {
	db, err := sql.Open("sqlite", "storage.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`
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
			(2, 10000, NULL, 1, '2026-08-10'),
			(3, 50000, NULL, 1, '2026-09-14');
	`)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("filled")
}

func main() {
	createTable()
	fillTable()
	fmt.Println("Hello!")
}
