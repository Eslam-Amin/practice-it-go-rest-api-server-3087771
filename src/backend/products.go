package backend

import (
	"database/sql"

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
