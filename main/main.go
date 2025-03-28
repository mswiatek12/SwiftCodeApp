package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"SwiftCodeApp/handlers"
	"SwiftCodeApp/repository"
	"SwiftCodeApp/service"
)

func main() {
	err := godotenv.Load("../config/.env")
	if err != nil {
		log.Fatalf("Unable to load .env file: %v", err)
	}

	fmt.Println(".env loaded")

	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer connPool.Close()

	fmt.Println("DB connection established.")

	sheetsService, err := service.NewSheetsService(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to Google Sheets API: %v", err)
	}

	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	readRange := "Sheet1!A1:G"

	banks, err := sheetsService.ReadSheet(spreadsheetId, readRange)
	if err != nil {
		log.Fatalf("Failed to read sheet: %v", err)
	}
	fmt.Println("Spreadsheet read.")
	repo := repository.NewBankRepository(connPool)
	handler := handlers.NewBankHandler(repo)

	err = repo.CreateRelation()
	if err != nil {
		fmt.Println("Error creating relation")
	}

	for _, bank := range banks {
		err := repo.InsertBankDb(bank)
		if err != nil {
			fmt.Printf("Error inserting data: %v", err)
		}
	}

	router := mux.NewRouter()

	router.HandleFunc("/v1/swift-codes/{swift-code}", handler.GetBankBySwiftCode).Methods("GET")
	router.HandleFunc("/v1/swift-codes/country/{countryISO2code}", handler.GetBanksByISO2Code).Methods("GET")
	router.HandleFunc("/v1/swift-codes", handler.InsertBank).Methods("POST")
	router.HandleFunc("/v1/swift-codes/{swift-code}", handler.DeleteBankBySwiftCode).Methods("DELETE")

	log.Println("App listening on port :8080")
	http.ListenAndServe(":8080", router)
}
