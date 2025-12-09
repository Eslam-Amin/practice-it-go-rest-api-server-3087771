package backend

import (
	"database/sql"
	"errors"
)

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

func NewOrder() *order {
	return &order{}
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
		err = ord.getOrderItems(db)
		if err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	return orders, nil
}

func getOrder(db *sql.DB, orderID int) (order, error) {
	var ord order
	err := db.QueryRow("Select * from orders where id = ?", orderID).Scan(&ord.ID, &ord.CustomerName, &ord.Total, &ord.Status)
	if err != nil {
		return ord, err
	}
	err = ord.getOrderItems(db)
	if err != nil {
		return ord, err
	}
	return ord, nil
}

func (ord *order) getOrderItems(db *sql.DB) error {
	rows, err := db.Query("Select * from order_items where order_id = ?", ord.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	orderItems := []orderItem{}
	for rows.Next() {
		var ordItem orderItem
		if err := rows.Scan(&ordItem.OrderID, &ordItem.ProductID, &ordItem.Quantity); err != nil {
			return err
		}
		orderItems = append(orderItems, ordItem)
	}
	ord.Items = orderItems
	return nil
}

func (ord *order) createOrder(db *sql.DB) error {
	res, err := db.Exec("INSERT into orders (customerName, totla, status) values (?, ?, ?) ", ord.CustomerName, ord.Total, ord.Status)
	if err != nil {
		return nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil
	}
	ord.ID = int(id)
	return nil
}

func (ordItem *orderItem) createOrderItem(db *sql.DB) error {
	_, err := db.Exec("INSERT into order_items (order_id, product_id, quantity) values (?, ?, ?) ", ordItem.OrderID, ordItem.ProductID, ordItem.Quantity)
	if err != nil {
		return nil
	}
	return nil
}
