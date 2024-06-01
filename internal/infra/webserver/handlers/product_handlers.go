package handlers

import (
	"encoding/json"
	"github.com/DieegoAlves/API/internal/dto"
	"github.com/DieegoAlves/API/internal/entity"
	"github.com/DieegoAlves/API/internal/infra/database"
	entityPkg "github.com/DieegoAlves/API/pkg/entity"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{ProductDB: db}
}

// CreateProduct godoc
// @Summary			Create Product
// @Description		Create a New Product
// @Tags			products
// @Accept			json
// @Produce			json
// @Param			request		body		dto.CreateProductInput		true 	"product request"
// @Success 		201
// @Failure			500 		{object}	Error
// @Failure 		401			{object}	Error
// @Router			/products	[post]
// @Security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ProductDB.CreateProduct(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetProduct godoc
// @Summary			Get a Product
// @Description		Show a specific product in Database
// @Tags			products
// @Accept			json
// @Produce			json
// @Param			id				path		string			true	"product id"	Format(uuid)
// @Success 		200				{object}	entity.Product
// @Failure			404
// @Failure 		500				{object}	Error
// @Router			/products/{id}	[get]
// @Security 		ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
// @Summary			Update a Product
// @Description		Update a Product in Database
// @Tags			products
// @Accept			json
// @Produce			json
// @Param			id				path		string						true	"product id"	Format(uuid)
// @Param			request			body		dto.CreateProductInput		true 	"product request"
// @Success 		200
// @Failure 		400
// @Failure			404
// @Failure			500 			{object}	Error
// @Router			/products/{id}	[put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.ProductDB.UpdateProduct(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteProduct godoc
// @Summary			Delete a Product
// @Description		Delete a Product in Database
// @Tags			products
// @Accept			json
// @Produce			json
// @Param			id				path		string		true	"product id"	Format(uuid)
// @Success 		200
// @Failure 		400
// @Failure			404
// @Failure			500 			{object}	Error
// @Router			/products/{id}	[delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.ProductDB.DeleteProduct(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetAllProducts godoc
// @Summary			List Products
// @Description		Show a list of all products in Database
// @Tags			products
// @Accept			json
// @Produce			json
// @Param			page		query		string	false	"page number"
// @Param			limit		query		string	false	"limit"
// @Success 		200			{array}		entity.Product
// @Failure			404 		{object}	Error
// @Failure 		500			{object}	Error
// @Router			/products	[get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// convert string to int
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0 // default value
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0 // default value
	}

	sort := r.URL.Query().Get("sort")

	if page == "" || limit == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(products)
	if err != nil {
		return
	}
}
