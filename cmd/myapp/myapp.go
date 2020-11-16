package main

import (
	"log"
	"os"
	"text/template"
	"time"

	"github.com/Muha113/golang-shop/pkg/db"
	"github.com/Muha113/golang-shop/pkg/models"
)

type ViewData struct {
	Shops []*models.Shop
}

func startAt(at time.Time, every time.Duration, exec chan<- bool) {
	t1, err := time.Parse(db.TimeLayout, time.Now().Format(db.TimeLayout))
	if err != nil {
		log.Fatal(err)
	}

	diff := at.Sub(t1)
	if diff < 0 {
		diff = (24 * time.Hour) + diff
	}

	for {
		<-time.After(diff)
		exec <- true
		diff = 3 * time.Second
	}
}

func main() {
	repo, err := db.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	exec := make(chan bool, 1)

	at, err := time.Parse(db.TimeLayout, "03:00:00")
	if err != nil {
		log.Fatal(err)
	}
	every := 24 * time.Hour

	go startAt(at, every, exec)

	for {
		select {
		case <-exec:
			res, err := repo.GetCurrentOpenedShopsProducts()
			if err != nil {
				log.Fatal(err)
			}

			tmpl, err := template.ParseFiles("./assets/template.xml")
			if err != nil {
				log.Fatal(err)
			}

			data := ViewData{
				Shops: res,
			}

			file, err := os.Create("./output/output.xml")
			if err != nil {
				log.Fatal(err)
			}

			err = tmpl.Execute(file, data)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
