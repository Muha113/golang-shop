package db

import (
	"database/sql"
	"sync"
	"time"

	"github.com/Muha113/golang-shop/pkg/models"
	"github.com/Muha113/golang-shop/pkg/parser"

	_ "github.com/go-sql-driver/mysql"
)

const TimeLayout = "15:04:05"

type Database struct {
	Repository *sql.DB
	mutex      *sync.Mutex
}

func InitDatabase() (*Database, error) {
	db, err := sql.Open("mysql", "root:password@/mydb")
	if err != nil {
		return nil, err
	}

	repo := &Database{
		Repository: db,
		mutex:      &sync.Mutex{},
	}

	return repo, nil
}

func (db *Database) GetCurrentOpenedShopsProducts() ([]*models.Shop, error) {
	db.mutex.Lock()

	defer db.mutex.Unlock()

	// now := time.Now().Format(TimeLayout)

	resultShops, err := db.Repository.Query("SELECT * FROM Shop WHERE CAST(? AS time) BETWEEN open AND close", "20:00:00")
	if err != nil {
		return nil, err
	}

	defer resultShops.Close()

	result := make([]*models.Shop, 0)

	var tmpOpenTime []uint8
	var tmpCloseTime []uint8

	for resultShops.Next() {
		tmpShop := &models.Shop{}

		err = resultShops.Scan(&tmpShop.ID, &tmpShop.Name, &tmpShop.URL, &tmpOpenTime, &tmpCloseTime)
		if err != nil {
			return nil, err
		}

		tmpShop.Open, err = time.Parse(TimeLayout, string(tmpOpenTime))
		if err != nil {
			return nil, err
		}

		tmpShop.Close, err = time.Parse(TimeLayout, string(tmpCloseTime))
		if err != nil {
			return nil, err
		}

		resultProducts, err := db.Repository.Query("SELECT * FROM Product WHERE shop_id = ?", tmpShop.ID)
		if err != nil {
			return nil, err
		}
		defer resultProducts.Close()

		var rowsAmount int
		err = db.Repository.QueryRow("SELECT COUNT(*) FROM Product WHERE shop_id = ?", tmpShop.ID).Scan(&rowsAmount)
		if err != nil {
			return nil, err
		}

		tmpShop.Products = make([]*models.Product, 0, rowsAmount)

		resChan := make(chan *models.Product, rowsAmount)
		group := sync.WaitGroup{}
		group.Add(rowsAmount)

		for resultProducts.Next() {
			tmpProduct := &models.Product{}

			err = resultProducts.Scan(&tmpProduct.ID, &tmpProduct.ShopID, &tmpProduct.Name, &tmpProduct.Description, &tmpProduct.Price)
			if err != nil {
				return nil, err
			}

			go func(res *models.Product, s *sync.WaitGroup, c chan<- *models.Product) {
				res.Description = parser.RemoveTagsHTML(res.Description)
				c <- tmpProduct
				s.Done()

			}(tmpProduct, &group, resChan)
		}

		group.Wait()
		close(resChan)

		for v := range resChan {
			tmpShop.Products = append(tmpShop.Products, v)
		}

		result = append(result, tmpShop)
	}

	return result, nil
}
