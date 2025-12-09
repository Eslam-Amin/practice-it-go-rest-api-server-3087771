package backend

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type product struct {
	ID          int     `json:"id"`
	ProductCode string  `json:"productCode"`
	Name        string  `json:"name"`
	Inventory   int     `json:"inventory"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
}

func NewProduct(productCode, name string, inventory int, price float64, status string) (*product, error) {
	if productCode == "" || name == "" || inventory == 0 || price == 0 || status == "" {
		return nil, errors.New("some fields are empty")
	}
	return &product{ProductCode: productCode, Name: name, Inventory: inventory, Price: price, Status: status}, nil
}

func getProducts(db *sql.DB) ([]product, error) {
	rows, err := db.Query("SELECT * FROM products")

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	products := []product{}
	for rows.Next() {
		var prod product
		if err := rows.Scan(&prod.ID, &prod.ProductCode, &prod.Name, &prod.Inventory, &prod.Price, &prod.Status); err != nil {
			return nil, err
		}
		products = append(products, prod)

	}
	return products, nil
}

func getProduct(db *sql.DB, id int) (product, error) {
	var prod product
	err := db.QueryRow("SELECT * FROM products WHERE id = ?", id).Scan(&prod.ID, &prod.ProductCode, &prod.Name, &prod.Inventory, &prod.Price, &prod.Status)
	return prod, err
}

func (prod *product) createProduct(db *sql.DB) error {
	res, err := db.Exec("INSERT into products (productCode, name, inventory, price, status) values (?, ?, ?, ?, ?)",
		prod.ProductCode, prod.Name, prod.Inventory, prod.Price, prod.Status)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	prod.ID = int(id)
	return nil
}
