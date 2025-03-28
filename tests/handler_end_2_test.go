package tests

import (
	"encoding/json"
	"context"
	"net/http/httptest"
	"testing"
	"log"
	"net/http"
	"SwiftCodeApp/handlers"
	"SwiftCodeApp/models"
	"SwiftCodeApp/repository"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetBankByISO2Code(t *testing.T) {
	t.Log("Starting TestGetBankByISO2Code...")

	repo := repository.NewBankRepository(testDB)
	bankHandler := handlers.NewBankHandler(repo)

	sqlStatement := `INSERT INTO branch (Address, BankName, CountryISO2, CountryName, IsHeadquarter, SwiftCode) 
        VALUES ($1, $2, $3, $4, $5, $6)`

	testBank := []models.Branch {
		{
			Address:       "That is an example address",
			BankName:      "Golden Sachs",
			CountryISO2:   "RO",
			CountryName:   "ROMANIA",
			IsHeadquarter: false,
			SwiftCode:     "KSZKSZ12345",
		},
		{
			Address:       "Nicaraghua 777",
			BankName:      "Simunic",
			CountryISO2:   "RO",
			CountryName:   "ROMANIA",
			IsHeadquarter: true,
			SwiftCode:     "ILUVGOLAXXX",
		},
	}

	_, err := testDB.Exec(context.Background(), sqlStatement, testBank[0].Address, testBank[0].BankName, testBank[0].CountryISO2, testBank[0].CountryName, testBank[0].IsHeadquarter, testBank[0].SwiftCode)
    if err != nil {
        t.Fatalf("Error inserting data: %v", err)
    }
	_, err = testDB.Exec(context.Background(), sqlStatement, testBank[1].Address, testBank[1].BankName, testBank[1].CountryISO2, testBank[1].CountryName, testBank[1].IsHeadquarter, testBank[1].SwiftCode)
    if err != nil {
        t.Fatalf("Error inserting data: %v", err)
    }

	req := httptest.NewRequest("GET", "/v1/swift-codes/country/" + testBank[0].CountryISO2, nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
    router.HandleFunc("/v1/swift-codes/country/{countryISO2code}", bankHandler.GetBanksByISO2Code).Methods("GET")

    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)

	var responseData map[string]interface{}

	err = json.NewDecoder(rr.Body).Decode(&responseData)
	if err != nil {
		t.Fatalf("Error decoding response body: %v", err)
	}

	expectedResponse := map[string]interface{} {
		"countryISO2": "RO",
		"countryName": "ROMANIA",
		"swiftCodes": []map[string]interface{} {
			{
				"address":       	"That is an example address",
				"bankName":      	"Golden Sachs",
				"countryISO2":   	"RO",
				"isHeadquarter": 	false,
				"swiftCode":     	"KSZKSZ12345",
			},
			{
				"address":       	"Nicaraghua 777",
				"bankName":      	"Simunic",
				"countryISO2":   	"RO",
				"isHeadquarter": 	true,
				"swiftCode":     	"ILUVGOLAXXX",
			},
		},
	}

	branches, ok := responseData["swiftCodes"].([]interface{})
	if !ok {
		t.Fatalf("Expected 'SwiftCodes' to be a slice but got: %T", responseData["swiftCodes"])
	}

	//converting from branches([]interface{}) to branchesAsMaps(map[string]interface{}) : - )

	var branchesAsMaps []map[string]interface{}
	for _, branch := range branches {
		if branchMap, ok := branch.(map[string]interface{});
		ok {
			branchesAsMaps = append(branchesAsMaps, branchMap)
		} else {
			t.Fatalf("Expected 'branch' to be a map but got: %T", branch)
		}
	}

	log.Printf("Actual response: %v", responseData)

	assert.Equal(t, expectedResponse["countryISO2"], responseData["countryISO2"])
	assert.Equal(t, expectedResponse["countryName"], responseData["countryName"])
	assert.Len(t, responseData["swiftCodes"], 2)
	assert.Equal(t, expectedResponse["swiftCodes"], branchesAsMaps)
	defer cleanUpTestData("KSZKSZ12345")
	defer cleanUpTestData("ILUVGOLAXXX")
}
