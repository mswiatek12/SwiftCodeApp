package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"SwiftCodeApp/handlers"
	"SwiftCodeApp/models"
	"SwiftCodeApp/repository"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestInsertBank(t *testing.T) {
	t.Log("Starting TestInsertBank...")

	repo := repository.NewBankRepository(testDB)
	bankHandler := handlers.NewBankHandler(repo)

	testBank := models.Branch{
		Address:       "Golang is cool",
		BankName:      "Random Albanian Bank",
		CountryISO2:   "AL",
		CountryName:   "ALBANIA",
		IsHeadquarter: false,
		SwiftCode:     "ALBXEXAMPLE",
	}

	testBankToInsert, err := json.Marshal(testBank)
	if err != nil {
		t.Fatalf("Couldn't encode to JSON: %v", err)
	}

	req := httptest.NewRequest("POST", "/v1/swift-codes/"+testBank.SwiftCode, bytes.NewBuffer(testBankToInsert))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/v1/swift-codes/{swift-code}", bankHandler.InsertBank).Methods("POST")

	router.ServeHTTP(rr, req)

	if rr.Code == http.StatusInternalServerError {
		t.Logf("Server error: %v", rr.Body.String())
	}	

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseData map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&responseData)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	expectedResponse := `{"message": "Data saved."}`

	actualResponse, err := json.Marshal(responseData)
	if err != nil {
		t.Fatalf("Couldn't encode actual response: %v", err)
	}

	assert.JSONEq(t, expectedResponse, string(actualResponse), "Response body does not match")
	defer cleanUpTestData("ALBXEXAMPLE")
}