package tests

import (
	"testing"
	"SwiftCodeApp/service"
)

func TestIsValidSwiftCode(t *testing.T) {
	swiftCodes := []struct {
		swiftCode string
		expected  bool
	}{
		{"ABDCEFGHIJK", true},
		{"ALBANIAaXXX", true},
		{"TOOSHORT", false},
		{"MAYBETOOLONG", false},
		{"BINGO123XXX", true},
		{"A@@LOFB1XXX", false},
	}

	for _, test := range swiftCodes {
		t.Run(test.swiftCode, func(t *testing.T) {
			actual := service.IsValidSwiftCode(test.swiftCode)
			if actual != test.expected {
				t.Errorf("Expected %v, but got %v for SwiftCode %s", test.expected, actual, test.swiftCode)
			}
		})
	}
}

func TestIsHeadquarterSwiftCode(t *testing.T) {
	swiftCodes := []struct {
		swiftCode string
		expected  bool
	}{
		{"ABDCEFGHIJK", false},
		{"ALBANIAaXXX", true},
		{"BINGO123XXX", true},
		{"A@@LOFB1XXX", false},
		{"TOOSHORT", false},
		{"MAYBETOOLONG", false},
	}

	for _, test := range swiftCodes {
		t.Run(test.swiftCode, func(t *testing.T) {
			actual := service.IsHeadquarter(test.swiftCode)
			if actual != test.expected {
				t.Errorf("Expected %v, but got %v for SwiftCode %s", test.expected, actual, test.swiftCode)
			}
		})
	}
}
