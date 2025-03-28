package handlers

import (
	"SwiftCodeApp/models"
	"SwiftCodeApp/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type BankHandler struct {
	Repo *repository.BankRepository
}

func NewBankHandler(repo *repository.BankRepository) *BankHandler {
	return &BankHandler{Repo: repo}
}

func (h *BankHandler) GetBankBySwiftCode(w http.ResponseWriter, r *http.Request) {
	swiftCode := mux.Vars(r)["swift-code"]

	log.Printf("Extracted SwiftCode: %s", swiftCode)

	sqlStatement := "SELECT * FROM branch WHERE branch.SwiftCode = $1"

	res, err := h.Repo.DB.Query(context.Background(), sqlStatement, swiftCode)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	var bank models.Branch
	if !res.Next() {
		http.Error(w, "Branch not found", http.StatusNotFound)
		return
	}

	err = res.Scan(&bank.Address, &bank.BankName, &bank.CountryISO2, &bank.CountryName, &bank.IsHeadquarter, &bank.SwiftCode)
	if err != nil {
		http.Error(w, "Error scanning row", http.StatusInternalServerError)
		return
	}

	if bank.IsHeadquarter {
		sqlBranches := "SELECT * FROM branch WHERE substring(swiftCode FROM 1 FOR 8) = substring($1 FROM 1 FOR 8) AND swiftCode != $1"
		branchesRes, err := h.Repo.DB.Query(context.Background(), sqlBranches, bank.SwiftCode)

		responseData := map[string]interface{}{
			"address":       bank.Address,
			"bankName":      bank.BankName,
			"countryISO2":   bank.CountryISO2,
			"countryName":   bank.CountryName,
			"isHeadquarter": bank.IsHeadquarter,
			"swiftCode":     bank.SwiftCode,
			"branches":      []models.BranchResponse{}, // an empty slice for naow
		}

		branches := responseData["branches"].([]models.BranchResponse)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning for branches: %v", err), http.StatusInternalServerError)
			return
		}

		for branchesRes.Next() {
			var entry models.Branch
			err := branchesRes.Scan(&entry.Address, &entry.BankName, &entry.CountryISO2, &entry.CountryName, &entry.IsHeadquarter, &entry.SwiftCode)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error scanning branch data: %v", err), http.StatusInternalServerError)
				return
			}
			branch := models.BranchResponse{
				Address:       entry.Address,
				BankName:      entry.BankName,
				CountryISO2:   entry.CountryISO2,
				IsHeadquarter: entry.IsHeadquarter,
				SwiftCode:     entry.SwiftCode,
			}
			branches = append(branches, branch)
		}

		responseData["branches"] = branches
		jsonResponse, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error serializing data to JSON: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)

	} else {
		responseData := map[string]interface{}{
			"address":       bank.Address,
			"bankName":      bank.BankName,
			"countryISO2":   bank.CountryISO2,
			"countryName":   bank.CountryName,
			"isHeadquarter": bank.IsHeadquarter,
			"swiftCode":     bank.SwiftCode,
		}

		jsonResponse, err := json.Marshal(responseData)
		if err != nil {
			http.Error(w, fmt.Sprintf("Couldn't serialize to JSON: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}
}


func (h *BankHandler) GetBanksByISO2Code(w http.ResponseWriter, r *http.Request) {
	countryISO2 := mux.Vars(r)["countryISO2code"]

	sqlStatement := "SELECT * FROM branch WHERE branch.CountryISO2 = $1"
	res, err := h.Repo.DB.Query(context.Background(), sqlStatement, countryISO2)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying DB: %v", err), http.StatusInternalServerError)
		return
	}
	defer res.Close()

	var branches []models.Branch
	var branchResponses []models.BranchResponse

	for res.Next() {
		var entry models.Branch
		err := res.Scan(&entry.Address, &entry.BankName, &entry.CountryISO2, &entry.CountryName, &entry.IsHeadquarter, &entry.SwiftCode)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning result: %v", err), http.StatusInternalServerError)
			return
		}
		branches = append(branches, entry)
		branchResponse := models.BranchResponse{
			Address:       entry.Address,
			BankName:      entry.BankName,
			CountryISO2:   entry.CountryISO2,
			IsHeadquarter: entry.IsHeadquarter,
			SwiftCode:     entry.SwiftCode,
		}
		branchResponses = append(branchResponses, branchResponse)
	}

	if len(branchResponses) == 0 {
		http.Error(w, "No branches found for the given country ISO2", http.StatusNotFound)
		return
	}

	responseData := map[string]interface{}{
		"countryISO2": branchResponses[0].CountryISO2,
		"countryName": branches[0].CountryName,
		"swiftCodes":  branchResponses,
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not parse data into JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (h *BankHandler) InsertBank(w http.ResponseWriter, r *http.Request) {
	var bank models.Branch
	err := json.NewDecoder(r.Body).Decode(&bank)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	_, err = h.Repo.DB.Exec(context.Background(), "INSERT INTO branch (Address, BankName, CountryISO2, CountryName, IsHeadquarter, SwiftCode) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (SwiftCode) DO NOTHING;",
		bank.Address,
		bank.BankName,
		bank.CountryISO2,
		bank.CountryName,
		bank.IsHeadquarter,
		bank.SwiftCode)

	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't insert data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Data saved."}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"message": "Err"}`, http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func (h *BankHandler) DeleteBankBySwiftCode(w http.ResponseWriter, r *http.Request) {
	SwiftCode := mux.Vars(r)["swift-code"]

	result, err := h.Repo.DB.Exec(context.Background(), "DELETE FROM branch WHERE SwiftCode=$1", SwiftCode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't delete data: %v", err), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, `{"error": "Bank not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{"message": "Data deleted."}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate response"}`, http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
