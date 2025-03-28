package tests

import (
	"encoding/json"
	"context"
	"net/http/httptest"
	"testing"
	"net/http"
	"SwiftCodeApp/handlers"
	"SwiftCodeApp/models"
	"SwiftCodeApp/repository"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestDeleteBankBySwiftCode(t *testing.T) {

	repo := repository.NewBankRepository(testDB)

	bankHandler := handlers.NewBankHandler(repo)

	sqlStatement := `INSERT INTO branch (Address, BankName, CountryISO2, CountryName, IsHeadquarter, SwiftCode) 
        VALUES ($1, $2, $3, $4, $5, $6)`

	bankToDelete := models.Branch {
		Address:		"Freedom street 900 Bank",
		BankName: 		"Remote bank of Iceland",
		CountryISO2: 	"US",
		CountryName: 	"UNITED STATES",
		IsHeadquarter: 	true,
		SwiftCode: 		"REMITLY9XXX",
	}

	_, err := testDB.Exec(context.Background(), sqlStatement, bankToDelete.Address, bankToDelete.BankName, bankToDelete.CountryISO2, bankToDelete.CountryName, bankToDelete.IsHeadquarter, bankToDelete.SwiftCode)
	if err != nil {
		t.Fatalf("Error inserting data: %v", err)
	}

	req := httptest.NewRequest("DELETE", "/v1/swift-codes/" + bankToDelete.SwiftCode, nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/v1/swift-codes/{swift-code}", bankHandler.DeleteBankBySwiftCode).Methods("DELETE")

	router.ServeHTTP(rr, req)

	if rr.Code == http.StatusInternalServerError {
		t.Logf("Server error: %v", rr.Body.String())
	}	

	assert.Equal(t, http.StatusOK, rr.Code)

	expectedResponse := `{"message": "Data deleted."}`

	var responseData map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&responseData)
	if err != nil {
		t.Fatalf("Couldn't encode actual response: %v", err)
	}

	actualResponse, err := json.Marshal(responseData)
	if err != nil {
		t.Fatalf("Couldn't encode actual response: %v", err)
	}

	assert.JSONEq(t, expectedResponse, string(actualResponse), "Response body does not match")

	t.Run("Invalid delete", func(t *testing.T) {
		invalidSwiftCode := "INVALIDSWIFTCODE123"

		req := httptest.NewRequest("DELETE", "/v1/swift-codes/" + invalidSwiftCode, nil)
		rr := httptest.NewRecorder()

		router := mux.NewRouter()
		router.HandleFunc("/v1/swift-codes/{swift-code}", bankHandler.DeleteBankBySwiftCode).Methods("DELETE")

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)

		expectedErrorResponse := `{"error": "Bank not found"}`

		var errorResponseData map[string]interface{}

		err = json.NewDecoder(rr.Body).Decode(&errorResponseData)
		if err != nil {
			t.Fatalf("Couldn't decode error response: %v", err)
		}

		actualErrorResponse, err := json.Marshal(errorResponseData)
		if err != nil {
			t.Fatalf("Couldn't encode error response: %v", err)
		}

		assert.JSONEq(t, expectedErrorResponse, string(actualErrorResponse), "Error response body does not match")
	})
}