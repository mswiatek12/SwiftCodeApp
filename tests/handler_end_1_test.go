package tests

import (
	"SwiftCodeApp/handlers"
	"SwiftCodeApp/models"
	"SwiftCodeApp/repository"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func cleanUpTestData(SwiftCode string) {
	_, err := testDB.Exec(context.Background(), "DELETE FROM branch WHERE swiftCode = $1", SwiftCode)
	if err != nil {
		log.Fatalf("Error cleaning up test data: %v", err)
	}
	log.Println("Data cleaned.")
}
func TestGetBankBySwiftCode(t *testing.T) {
    t.Log("Starting TestGetBankBySwiftCode...")

    testSwiftCode := "EXAMPLUKXXX"
    repo := repository.NewBankRepository(testDB)
    bankHandler := handlers.NewBankHandler(repo)

    sqlStatement := `INSERT INTO branch (Address, BankName, CountryISO2, CountryName, IsHeadquarter, SwiftCode) 
        VALUES ($1, $2, $3, $4, $5, $6)`

    testBank := models.Branch{
        Address:       "Aguaricona 3",
        BankName:      "Simunic",
        CountryISO2:   "UK",
        CountryName:   "UNITED KINGDOM",
        IsHeadquarter: true,
        SwiftCode:     testSwiftCode,
    }

	testBranch := models.Branch {
		Address:		"Mill Hill East 21",
		BankName:		"Deloitte",
		CountryISO2: 	"UK",
		IsHeadquarter: 	false,
		SwiftCode: 		"EXAMPLUK133",
	}

    _, err := testDB.Exec(context.Background(), sqlStatement, testBank.Address, testBank.BankName, testBank.CountryISO2, testBank.CountryName, testBank.IsHeadquarter, testSwiftCode)
    if err != nil {
        t.Fatalf("Error inserting data: %v", err)
    }
	_, err = testDB.Exec(context.Background(), sqlStatement, testBranch.Address, testBranch.BankName, testBranch.CountryISO2, testBranch.CountryName, testBranch.IsHeadquarter, testBranch.SwiftCode)
    if err != nil {
        t.Fatalf("Error inserting data: %v", err)
    }

    t.Logf("Request URL: /v1/swift-codes/%s", testSwiftCode)

    req := httptest.NewRequest("GET", "/v1/swift-codes/"+testSwiftCode, nil)
    rr := httptest.NewRecorder()

    router := mux.NewRouter()
    router.HandleFunc("/v1/swift-codes/{swift-code}", bankHandler.GetBankBySwiftCode).Methods("GET")

    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

    var responseData map[string]interface{}
    err = json.NewDecoder(rr.Body).Decode(&responseData)
    if err != nil {
        t.Fatalf("Error decoding response body: %v", err)
    }

	expectedBranch := map[string]interface{}{
		"address":       "Mill Hill East 21",
		"bankName":      "Deloitte",
		"countryISO2":   "UK",
		"isHeadquarter": false,
		"swiftCode":     "EXAMPLUK133",
	}

	expectedResponse := map[string]interface{} {
		"address":       "Aguaricona 3",
		"bankName":      "Simunic",
		"countryISO2":   "UK",
		"countryName":   "UNITED KINGDOM",
		"isHeadquarter": true,
		"swiftCode":     testSwiftCode,
		"branches": []map[string]interface{}{
			{
				"address":       "Mill Hill East 21",
				"bankName":      "Deloitte",
				"countryISO2":   "UK",
				"isHeadquarter": false,
				"swiftCode":     "EXAMPLUK133",
			},
		},
	}

	branches, ok := responseData["branches"].([]interface{})
	if !ok {
		t.Fatalf("Expected 'branches' to be a slice but got: %T", responseData["branches"])
	}
    assert.Equal(t, expectedResponse["address"], responseData["address"])
    assert.Equal(t, expectedResponse["bankName"], responseData["bankName"])
    assert.Equal(t, expectedResponse["countryISO2"], responseData["countryISO2"])
    assert.Equal(t, expectedResponse["countryName"], responseData["countryName"])
    assert.Equal(t, expectedResponse["isHeadquarter"], responseData["isHeadquarter"])
    assert.Equal(t, expectedResponse["swiftCode"], responseData["swiftCode"])
	assert.Len(t, branches, 1)
	assert.Equal(t, expectedBranch, branches[0])
	
    defer cleanUpTestData("EXAMPLUKXXX")
	defer cleanUpTestData("EXAMPLUK133")
}