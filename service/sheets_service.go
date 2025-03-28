package service

import (
	"fmt"
	"log"
	"strings"
	"context"
	"regexp"
	"SwiftCodeApp/models"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func IsValidSwiftCode(swiftCode string) bool {
	if len(swiftCode) != 11 {
		return false
	}

	match, _ := regexp.MatchString(`^[A-Za-z0-9]+$`, swiftCode)
	return match
}

func IsHeadquarter(swiftCode string) bool {
	return strings.HasSuffix(swiftCode, "XXX") && IsValidSwiftCode(swiftCode)
}

type SheetsAPI interface {
	SpreadsheetsValuesGet(spreadsheetId, readRange string) (*sheets.ValueRange, error)
}

type SheetsService struct {
	api SheetsAPI
}

func NewSheetsService(ctx context.Context) (*SheetsService, error) {
	creds, err := google.FindDefaultCredentials(ctx, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	srv, err := sheets.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	return &SheetsService{api: &RealSheetsApi{srv}}, nil
}

func NewSheetsServiceWithAPI(api SheetsAPI) *SheetsService {
	return &SheetsService{api: api}
}

type RealSheetsApi struct {
	srv *sheets.Service
}

func (r *RealSheetsApi) SpreadsheetsValuesGet(spreadsheetId, readRange string) (*sheets.ValueRange, error) {
	return r.srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
}

func (s *SheetsService) ReadSheet(spreadsheetId string, readRange string) ([]models.Branch, error) {
	var banks []models.Branch

	resp, err := s.api.SpreadsheetsValuesGet(spreadsheetId, readRange)
	if err != nil {
		return nil, err
	}

	if len(resp.Values) == 0 {
		log.Println("No data found.")
		return banks, nil
	}

	for _, row := range resp.Values[1:] {
		banks = append(banks, models.Branch{
			Address:       fmt.Sprintf("%s", row[4]),
			BankName:      fmt.Sprintf("%s", row[3]),
			CountryISO2:   strings.ToUpper(fmt.Sprintf("%s", row[0])),
			CountryName:   strings.ToUpper(fmt.Sprintf("%s", row[6])),
			SwiftCode:     fmt.Sprintf("%s", row[1]),
			IsHeadquarter: IsHeadquarter(fmt.Sprintf("%s", row[1])),
		})
	}
	return banks, nil
}
