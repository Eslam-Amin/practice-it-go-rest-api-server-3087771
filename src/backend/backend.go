package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// RespondWithError writes an error response with the given code and message.
// It uses respondWithJSON to write the response as JSON.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON writes a JSON response with the given code and payload.
// It marshals the payload as JSON, sets the Content-Type header to application/json,
// writes the given HTTP status code, and then writes the JSON response.
// If json.Marshal returns an error, it will panic.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type App struct {
	DB     *sql.DB
	Port   string
	Router *mux.Router
}

// Initialize sets up the database connection and initializes the router
// by calling initializeRoutes.
func (app *App) Initialize() {
	db, err := sql.Open("sqlite3", "../../practiceit.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	app.DB = db
	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

// initializeRoutes sets up the routes for the application.
// It maps the /products endpoint to the getProducts and createProduct handlers,
// and the /products/{id} endpoint to the getProduct handler.
// It maps the /orders endpoint to the getAllOrders and createOrder handlers,
// and the /orders/{id} endpoint to the getOrder handler.
func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/products", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/products/{id}", app.getProduct).Methods("GET")

	app.Router.HandleFunc("/orders", app.getAllOrders).Methods("GET")
	app.Router.HandleFunc("/orders", app.createOrder).Methods("POST")
	app.Router.HandleFunc("/orders/{id}", app.getOrder).Methods("GET")
}

// CreateProduct creates a new product in the database.
// It reads the request body as JSON and attempts to create a new product
// with the given details. If the request body is invalid JSON, or if the
// product details are invalid, it returns a 400 error response. If the
// product creation fails for some reason, it returns a 500 error response.
// Otherwise, it returns a 201 response with the newly created product.
func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var productInput product
	json.Unmarshal(body, &productInput)
	prod, err := NewProduct(productInput.ProductCode, productInput.Name, productInput.Inventory, productInput.Price, productInput.Status)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	err = prod.createProduct(app.DB)
	if err != nil {
		fmt.Printf("CreateProduct Error: %v", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, prod)
}

// getProducts retrieves a list of all products from the database.
// If there is an error retrieving the products, it will return a 500 error response.
// Otherwise, it will return a 200 response with the list of products.
func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		fmt.Printf("GetProducts Error: %v", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

// getProduct retrieves a single product from the database.
// It expects the product ID as a URL parameter, and returns a 404 error response if the product is not found.
// If there is an error retrieving the product, it will return a 500 error response.
// Otherwise, it will return a 200 response with the product.
func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
		return
	}
	product, err := getProduct(app.DB, id)
	if err != nil {
		fmt.Printf("GetProduct Error: %v", err.Error())
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Product Not Found")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, product)
}

// createOrder creates a new order in the database.
// It expects a request body as JSON, and attempts to create a new order
// with the given details. If the request body is invalid JSON, or if the
// order details are invalid, it will return a 400 error response. If the
// order creation fails for some reason, it will return a 500 error response.
// Otherwise, it will return a 201 response with the newly created order.
// It also creates a new order item for each item in the order, and returns a 500
// error response if any of the order items cannot be created.
func (app *App) createOrder(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var orderInput order
	json.Unmarshal(body, &orderInput)
	ord, err := NewOrder(orderInput.CustomerName, orderInput.Total)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = ord.createOrder(app.DB)
	if err != nil {
		fmt.Printf("CreateOrder Error: %v", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	orderItems := []orderItem{}
	for _, item := range orderInput.Items {
		ordItem, err := NewOrderItem(ord.ID, item.ProductID, item.Quantity)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		err = ordItem.createOrderItem(app.DB)
		if err != nil {
			fmt.Printf("CreateOrderItem Error: %v", err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		orderItems = append(orderItems, *ordItem)
	}
	ord.Items = orderItems
	respondWithJSON(w, http.StatusCreated, ord)
}

// getAllOrders retrieves a list of all orders from the database.
// If there is an error retrieving the orders, it will return a 500 error response.
// Otherwise, it will return a 200 response with the list of orders.
func (app *App) getAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := getAllOrders(app.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, orders)
}

// getOrder retrieves a single order from the database.
// It expects the order ID as a URL parameter, and returns a 404 error response if the order is not found.
// If there is an error retrieving the order, it will return a 500 error response.
// Otherwise, it will return a 200 response with the order.
func (app *App) getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Order ID")
		return
	}
	order, err := getOrder(app.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Order Not Found")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, order)
}

// Run starts the HTTP server and begins listening for incoming requests.
// It prints a message to the console indicating that the server has started
// and is listening on the specified port. It then calls http.ListenAndServe to
// start the server, passing in the port and nil as the handler. If there
// is an error starting the server, it will be logged to the console using
// log.Fatal.
func (app *App) Run() {
	fmt.Printf("Server Started and listening on Port: %v", app.Port)
	log.Fatal(http.ListenAndServe(app.Port, nil))
}
