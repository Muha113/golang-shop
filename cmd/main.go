package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Muha113/golang-shop/pkg/db"
	"github.com/Muha113/golang-shop/pkg/models"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	repo, err := db.InitDatabase("./config/config.json")
	if err != nil {
		log.Fatal(err)
	}

	res, err := repo.GetCurrentOpenedShopsProducts()
	if err != nil {
		log.Fatal(err)
	}

	for res.Next() {
		s := models.Shop{}
		var tmpOpen []uint8
		var tmpClose []uint8
		err = res.Scan(&s.ID, &s.Name, &s.URL, &tmpOpen, &tmpClose)
		if err != nil {
			log.Fatal(err)
		}

		s.Open, err = time.Parse("2006-01-02 15:04:05", string(tmpOpen))
		if err != nil {
			log.Fatal(err)
		}

		s.Close, err = time.Parse("2006-01-02 15:04:05", string(tmpClose))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(s)
	}

	// fmt.Println()

	// fmt.Println(result.LastInsertId()) // id добавленного объекта
	// fmt.Println(result.RowsAffected()) // количество затронутых строк
}
