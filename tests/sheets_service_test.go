package tests

import (
	"testing"
	"SwiftCodeApp/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/api/sheets/v4"
)

type MockSheetsAPI struct {
	mock.Mock
}

func (m *MockSheetsAPI) SpreadsheetsValuesGet(spreadsheetId, readRange string) (*sheets.ValueRange, error) {
	args := m.Called(spreadsheetId, readRange)
	val, ok := args.Get(0).(*sheets.ValueRange)
	if !ok {
		return nil, args.Error(1)
	}
	return val, args.Error(1)
}

func TestReadSheet_Success(t *testing.T) {
	mockAPI := new(MockSheetsAPI)

	mockResponse := &sheets.ValueRange{
		Values: [][]interface{}{
			// CODE SKIPS FIRST LINE SO FIRST ENTRY WONT BE IN "BANKS" CHECK "ReadSheet" func
			{"PL", "LUV4REMITLY", "BIC11", "REMITLY HAS IT", "Street of wealthiest", "Krakow", "POLAND", "Europe/Warsaw"},
			{"PL", "LUV4REMITLY", "BIC11", "REMITLY HAS IT", "Street of wealthiest", "Krakow", "POLAND", "Europe/Warsaw"},
			{"US", "REMITLYGXXX", "BIC11", "SOME EXAMPLE", "Success st", "New York", "UNITED STATES", "Pacific/Easter"},
		},
	}

	mockAPI.On("SpreadsheetsValuesGet", mock.Anything, mock.Anything).Return(mockResponse, nil)

	sheetService := service.NewSheetsServiceWithAPI(mockAPI)

	banks, err := sheetService.ReadSheet("spreadsheet-id", "Sheet1!A1:G")

	assert.NoError(t, err)
	assert.Len(t, banks, 2)

	assert.Equal(t, "LUV4REMITLY", banks[0].SwiftCode)
	assert.Equal(t, "PL", banks[0].CountryISO2)
	assert.False(t, banks[0].IsHeadquarter)

	assert.Equal(t, "REMITLYGXXX", banks[1].SwiftCode)
	assert.Equal(t, "US", banks[1].CountryISO2)
	assert.True(t, banks[1].IsHeadquarter)

	mockAPI.AssertExpectations(t)
}

func TestReadSheet_Failure(t *testing.T) {
	mockAPI := new(MockSheetsAPI)

	mockAPI.On("SpreadsheetsValuesGet", mock.Anything, mock.Anything).Return(nil, assert.AnError)

	sheetService := service.NewSheetsServiceWithAPI(mockAPI)

	banks, err := sheetService.ReadSheet("spreadsheet-id", "Sheet1!A1:G")

	assert.Error(t, err)
	assert.Nil(t, banks)

	mockAPI.AssertExpectations(t)
}