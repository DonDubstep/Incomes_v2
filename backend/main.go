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
			period TEXT NOT NULL UNIQUE
		);

		CREATE TABLE IF NOT EXISTS incomes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			period_id INTEGER NOT NULL,
			amount INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			description TEXT,
			date TEXT NOT NULL,
			FOREIGN KEY (period_id) REFERENCES period(id),
			FOREIGN KEY (category_id) REFERENCES category(id)
		);
	
	`)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("created")
}

func main() {
	createTable()
	fmt.Println("Hello!")
}
