package db

import "database/sql"

type Database struct {
	Repository *sql.DB
}

func InitDatabase(cfgPath string) (*Database, error) {
	db, err := sql.Open("mysql", "root:password@/mydb")
	if err != nil {
		return nil, err
	}

	repo := &Database{
		Repository: db,
	}

	return repo, nil
}

func (db *Database) GetCurrentOpenedShopsProducts() (*sql.Rows, error) {
	result, err := db.Repository.Query("SELECT * FROM Shop WHERE name = ?", "myshop")
	if err != nil {
		return nil, err
	}
	return result, nil
}
