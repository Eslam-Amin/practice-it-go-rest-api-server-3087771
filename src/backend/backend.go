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

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

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

func (app *App) Initialize() {
	db, err := sql.Open("sqlite3", "../../practiceit.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	app.DB = db
	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

func (app *App) initializeRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
	app.Router.HandleFunc("/products", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/products/{id}", app.getProduct).Methods("GET")
}

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

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		fmt.Printf("GetProducts Error: %v", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, products)
}

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
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, product)
}

func (app *App) Run() {
	fmt.Printf("Server Started and listening on Port: %v", app.Port)
	log.Fatal(http.ListenAndServe(app.Port, nil))
}
