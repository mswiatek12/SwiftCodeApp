package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"SwiftCodeApp/models"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func readSheet(srv *sheets.Service, spreadsheetId string, readRange string, banks []models.Branch) []models.Branch {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data: %v", err)
	}

	if len(resp.Values) == 0 {
		log.Println("No data found.")
	} else {
		for _, row := range resp.Values[1:] {												// skipping the first row with column headers
			banks = append(banks, models.Branch {												// using fmt.Sprintf("%s", row[6]) to convert interface{} to string
				Address:     	fmt.Sprintf("%s", row[4]),
				BankName:    	fmt.Sprintf("%s", row[3]),
				CountryISO2: 	fmt.Sprintf("%s", row[0]),
				CountryName: 	fmt.Sprintf("%s", row[6]),
				SwiftCode: 		fmt.Sprintf("%s", row[1]),
				IsHeadquarter: 	strings.Contains(fmt.Sprintf("%s", row[1]), "XXX"), 		// "XXX" substring does not appear in other positions than the end so it works
			})
		}
	}
	return banks
}

func main() {
	ctx := context.Background()

	creds, err := google.FindDefaultCredentials(ctx, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to find default credentials: %v", err)
	}
	fmt.Println("Credentials succesfully read")
	srv, err := sheets.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	log.Println("Succesfully connected to Google Sheets Api")
	spreadsheetId := "1iFFqsu_xruvVKzXAadAAlDBpIuU51v-pfIEU5HeGa8w"
	readRange := "Sheet1!A1:G"
	var banks []models.Branch
	banks = readSheet(srv, spreadsheetId, readRange, banks)
	fmt.Print(banks[0])
}
