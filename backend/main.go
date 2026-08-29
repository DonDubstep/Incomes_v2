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

func (s *Server) GetAllIncomesHandler(w http.ResponseWriter, r *http.Request) {
	result, err := calculate.GetAllIncomes(s.db)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сервера:\n%s", err), http.StatusInternalServerError)
		return
	}

	// Указываем, что будем отправлять в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// Кодируем срез в json и отправляем
	json.NewEncoder(w).Encode(result)
}

func (s *Server) GetIncomesByMonthHandler(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("period")
	result, err := calculate.GetIncomesByMonth(s.db, month)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сервера:\n%s", err), http.StatusInternalServerError)
		return
	}

	// Указываем, что будем отправлять в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// Кодируем срез в json и отправляем
	json.NewEncoder(w).Encode(result)
}

func (s *Server) GetAllIncomesRecordsByMonthHandler(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("period")
	result, err := calculate.GetAllIncomesRecordsByMonth(s.db, month)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка сервера:\n%s", err), http.StatusInternalServerError)
		return
	}

	// Указываем, что будем отправлять в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// Кодируем срез в json и отправляем
	json.NewEncoder(w).Encode(result)
}

func (s *Server) GetIncomesRecordsByRangeHandler(w http.ResponseWriter, r *http.Request) {
	start_month := r.URL.Query().Get("from")
	end_month := r.URL.Query().Get("to")
	result, err := calculate.GetIncomesRecordsByRange(s.db, start_month, end_month)
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
		sql_path = "storage.db"
	}
	db, err := calculate.ConnectToDataBase(sql_path)
	if err != nil {
		fmt.Println("ошибка открытия БД %w", err)
		return
	}
	defer db.Close()
	calculate.CreateTable(db)
	calculate.FillTable(db)

	server := Server{db}

	http.HandleFunc("/1", server.GetAllIncomesHandler)
	http.HandleFunc("/2", server.GetIncomesByMonthHandler)
	http.HandleFunc("/3", server.GetAllIncomesRecordsByMonthHandler)
	http.HandleFunc("/4", server.GetIncomesRecordsByRangeHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	fmt.Println("Сервер запущен на порту", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
