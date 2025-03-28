package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"SwiftCodeApp/models"
	"SwiftCodeApp/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
    err := godotenv.Load("../config/.env")
    if err != nil {
        log.Fatalf("Unable to load .env file: %v", err)
    }
    
    testDB, err = pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
    if err != nil {
        log.Fatalf("Failed to connect to test database: %v", err)
    }
    
    defer testDB.Close()

    code := m.Run()

    os.Exit(code)
}

func TestInsertBankDb(t *testing.T) {
	t.Log("Starting TestInsertBankDb...")

	repo := repository.NewBankRepository(testDB)
	repo.CreateRelation()
	testBank := models.Branch{
		Address:       "Pawia 17 st",
		BankName:      "I love Remitly",
		CountryISO2:   "PL",
		CountryName:   "POLAND",
		IsHeadquarter: true,
		SwiftCode:     "ILUVREMITLY",
	}

	err := repo.InsertBankDb(testBank)
	assert.NoError(t, err)

	var count int
	err = testDB.QueryRow(context.Background(), "SELECT COUNT(*) FROM branch WHERE SwiftCode = $1", testBank.SwiftCode).Scan(&count)
	assert.NoError(t, err)

	t.Logf("Bank count in DB: %d", count)
	assert.Equal(t, 1, count)

	_, err = testDB.Exec(context.Background(), "DELETE FROM branch WHERE SwiftCode = $1", testBank.SwiftCode)
	assert.NoError(t, err)

	t.Log("Cleanup successful")
}
