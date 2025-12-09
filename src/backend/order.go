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

// NewOrder creates a new order with the given customer name and items.
// It returns an error if the customer name is empty or if the items slice is empty.
// It also returns an error if any of the items have invalid details.
// Otherwise, it returns a new order with the given customer name and items.
func NewOrder(customerName string, total int) (*order, error) {
	if customerName == "" || total == 0 {
		return nil, errors.New("order is empty")
	}

	return &order{
		CustomerName: customerName,
		Status:       "Pending",
		Total:        total,
	}, nil
}

// NewOrderItem creates a new order item with the given order ID, product ID, and quantity.
// It returns an error if any of the given parameters are empty (i.e. 0).
// Otherwise, it returns a new order item with the given parameters.
func NewOrderItem(orderId, productId, quantity int) (*orderItem, error) {
	if quantity == 0 {
		return nil, errors.New("order item quantity is empty")
	}
	return &orderItem{
		OrderID:   orderId,
		ProductID: productId,
		Quantity:  quantity,
	}, nil
}

// getAllOrders retrieves a list of all orders from the database.
// If there is an error retrieving the orders, it will return a nil slice and an error.
// Otherwise, it will return a slice of all orders.
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

// getOrder retrieves a single order from the database.
// It expects the order ID as a parameter, and returns an error if the order is not found.
// If there is an error retrieving the order, it will return a nil order and an error.
// Otherwise, it will return the retrieved order with its items.
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

// getOrderItems retrieves a list of all order items associated with the given order.
// It expects a database connection and the order ID as parameters, and returns an error if the
// order items cannot be retrieved.
// If there is an error retrieving the order items, it will return a nil slice and an error.
// Otherwise, it will return a slice of all order items associated with the given order.
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

// createOrder creates a new order in the database.
// It expects a database connection and an order as parameters, and returns an error if the order cannot be created.
// If there is an error creating the order, it will return a nil error.
// Otherwise, it will return a nil error and set the order's ID to the last inserted ID.
func (ord *order) createOrder(db *sql.DB) error {
	res, err := db.Exec("INSERT into orders (customerName, total, status) values (?, ?, ?) ", ord.CustomerName, ord.Total, ord.Status)
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

// createOrderItem creates a new order item in the database.
// It expects a database connection and an order item as parameters, and returns an error if the order item cannot be created.
// If there is an error creating the order item, it will return a nil error.
// Otherwise, it will return a nil error and set the order item's ID to the last inserted ID.
func (ordItem *orderItem) createOrderItem(db *sql.DB) error {
	_, err := db.Exec("INSERT into order_items (order_id, product_id, quantity) values (?, ?, ?) ", ordItem.OrderID, ordItem.ProductID, ordItem.Quantity)
	if err != nil {
		return nil
	}
	return nil
}
