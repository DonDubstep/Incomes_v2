package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"incomes/calculate"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

type Server struct {
	db *sql.DB
}

func (s *Server) GetAmountOfMoneyRubHandler(w http.ResponseWriter, r *http.Request) {
	result, err := calculate.GetAmountOfMoneyRub(s.db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сервера:\n%s", err), http.StatusInternalServerError)
		return
	}

	// Указываем, что будем отправлять в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// Кодируем срез в json и отправляем
	json.NewEncoder(w).Encode(result)
}

func (s *Server) GetAmountOfLastIncomeHandler(w http.ResponseWriter, r *http.Request) {
	result, err := calculate.GetAmountOfLastIncome(s.db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сервера:\n%s", err), http.StatusInternalServerError)
		return
	}

	// Указываем, что будем отправлять в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// Кодируем срез в json и отправляем
	json.NewEncoder(w).Encode(result)
}

func (s *Server) GetTotalLastChangeHandler(w http.ResponseWriter, r *http.Request) {
	result, err := calculate.GetTotalLastChange(s.db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сервера:\n%s", err), http.StatusInternalServerError)
		return
	}

	// Указываем, что будем отправлять в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// Кодируем срез в json и отправляем
	json.NewEncoder(w).Encode(result)
}

func (s *Server) GetTotalAmountInCurrenciesHandler(w http.ResponseWriter, r *http.Request) {
	result, err := calculate.GetTotalAmountInCurrencies(s.db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сервера:\n%s", err), http.StatusInternalServerError)
		return
	}

	// Указываем, что будем отправлять в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// Кодируем срез в json и отправляем
	json.NewEncoder(w).Encode(result)
}

func (s *Server) GetLast3IncomesHandlerHandler(w http.ResponseWriter, r *http.Request) {
	result, err := calculate.GetLast3Incomes(s.db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сервера:\n%s", err), http.StatusInternalServerError)
		return
	}

	// Указываем, что будем отправлять в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// Кодируем срез в json и отправляем
	json.NewEncoder(w).Encode(result)
}

func main() {
	sql_path := os.Getenv("DB_PATH")
	if sql_path == "" {
		sql_path = "../data/storage.db"
	}
	db, err := calculate.ConnectToDataBase(sql_path)
	if err != nil {
		fmt.Println("ошибка открытия БД %w", err)
		return
	}
	defer db.Close()

	// if err = calculate.ResetDatabase(db); err != nil {
	// 	fmt.Print(err)
	// }
	// calculate.CreateTable(db)
	// calculate.FillTable(db)

	server := Server{db}

	http.HandleFunc("/api/5", server.GetAmountOfMoneyRubHandler)
	http.HandleFunc("/api/6", server.GetAmountOfLastIncomeHandler)
	http.HandleFunc("/api/7", server.GetTotalLastChangeHandler)
	http.HandleFunc("/api/8", server.GetTotalAmountInCurrenciesHandler)
	http.HandleFunc("/api/9", server.GetLast3IncomesHandlerHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	fmt.Println("Сервер запущен на порту", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
