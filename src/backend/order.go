package backend

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
