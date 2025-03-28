package repository

import (
	"SwiftCodeApp/models"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BankRepository struct {
	DB *pgxpool.Pool
}

func NewBankRepository(db *pgxpool.Pool) *BankRepository {
	return &BankRepository{DB: db}
}

func (repo *BankRepository) CreateRelation() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS branch (
		Address TEXT,
		BankName TEXT,
		CountryISO2 TEXT,
		CountryName TEXT,
		IsHeadquarter BOOLEAN,
		SwiftCode TEXT UNIQUE
	);`

	_, err := repo.DB.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Printf("Error creating table: %v", err)
	}

	return err
}

func (repo *BankRepository) InsertBankDb(bank models.Branch) error {
	sqlStatement := `INSERT INTO branch (Address, BankName, CountryISO2, CountryName, IsHeadquarter, SwiftCode)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (SwiftCode) DO NOTHING;`

	_, err := repo.DB.Exec(context.Background(), sqlStatement,
		bank.Address,
		bank.BankName,
		bank.CountryISO2,
		bank.CountryName,
		bank.IsHeadquarter,
		bank.SwiftCode,
	)

	return err
}
