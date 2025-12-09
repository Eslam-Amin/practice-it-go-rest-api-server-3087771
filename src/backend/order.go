package backend

import "database/sql"

type order struct {
	ID           int         `json:"id"`
	CustomerName string      `json:"customerName"`
	Total        int         `json:"total"`
	Status       string      `json:"status"`
	Items        []orderItem `json:"items"`
}

type orderItem struct {
	OrderID   int `json:"orderId"`
	ProductID int `json:"productId"`
	Quantity  int `json:"quantity"`
}

func getAllOrders(db *sql.DB) ([]order, error) {
	rows, err := db.Query("Select * from orders")
	if err != nil {
		return nil, err
	}
	orders := []order{}
	defer rows.Close()
	for rows.Next() {
		var ord order
		if err := rows.Scan(&ord.ID, &ord.CustomerName, &ord.Total, &ord.Status); err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	return orders, nil
}
