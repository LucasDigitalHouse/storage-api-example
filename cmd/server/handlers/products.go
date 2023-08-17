package handlers

import (
	"app/internal/products/storage"
	"app/pkg/web/request"
	"app/pkg/web/response"
	"errors"
	"net/http"
	"strconv"
)

// NewControllerProduct returns new ControllerProduct
func NewControllerProduct(storage storage.StorageProduct) *ControllerProduct {
	return &ControllerProduct{storage: storage}
}

// ControllerProduct is a controller for products
type ControllerProduct struct {
	// storage is a storage for products
	storage storage.StorageProduct
}

// GetOne returns one product by id
type ResponseProduct struct {
	Name    string	`json:"name"`
	Type	string	`json:"type"`
	Count	int		`json:"count"`
	Price	float64	`json:"price"`
}
type ResponseBody struct {
	Message string			 `json:"message"`
	Data    *ResponseProduct `json:"data"`
	Error   bool			 `json:"error"`
}
func (c *ControllerProduct) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		idParam, err := request.PathLastParam(r)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "invalid path param", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "parameter must be int", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}


		// process
		product, err := c.storage.GetOne(id)
		if err != nil {
			var code int; var body *ResponseBody
			switch {
			case errors.Is(err, storage.ErrStorageProductNotFound):
				code = http.StatusNotFound
				body = &ResponseBody{Message: "product not found", Data: nil, Error: true}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBody{Message: "internal error", Data: nil, Error: true}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseBody{
			Message: "success",
			Data: &ResponseProduct{			// serialization
				Name:   product.Name,
				Type:	product.Type,
				Count:	product.Count,
				Price:	product.Price,
			},
			Error: false,
		}

		response.JSON(w, code, body)
	}
}

// Store stores product
type RequestProductStore struct {
	Name    string	`json:"name"`
	Type	string	`json:"type"`
	Count	int		`json:"count"`
	Price	float64	`json:"price"`
}
type ResponseProductStore struct {
	Name    string	`json:"name"`
	Type	string	`json:"type"`
	Count	int		`json:"count"`
	Price	float64	`json:"price"`
}
type ResponseBodyStore struct {
	Message string					`json:"message"`
	Data    *ResponseProductStore	`json:"data"`
	Error   bool					`json:"error"`
}
func (c *ControllerProduct) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var req RequestProductStore
		err := request.JSON(r, &req)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "invalid json", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// process
		// -> deserialization
		product := &storage.Product{
			Name:   req.Name,
			Type:	req.Type,
			Count:	req.Count,
			Price:	req.Price,
		}
		err = c.storage.Store(product)
		if err != nil {
			var code int; var body *ResponseBody
			switch {
			case errors.Is(err, storage.ErrStorageProductNotUnique):
				code = http.StatusBadRequest
				body = &ResponseBody{Message: "product not unique", Data: nil, Error: true}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBody{Message: "internal error", Data: nil, Error: true}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusCreated
		body := &ResponseBodyStore{
			Message: "success",
			Data: &ResponseProductStore{	// serialization
				Name:   product.Name,
				Type:	product.Type,
				Count:	product.Count,
				Price:	product.Price,
			},
			Error: false,
		}

		response.JSON(w, code, body)
	}
}

// Update updates product
type RequestProductUpdate struct {
	Name    string	`json:"name"`
	Type	string	`json:"type"`
	Count	int		`json:"count"`
	Price	float64	`json:"price"`
}
type ResponseProductUpdate struct {
	Name    string	`json:"name"`
	Type	string	`json:"type"`
	Count	int		`json:"count"`
	Price	float64	`json:"price"`
}
type ResponseBodyUpdate struct {
	Message string					`json:"message"`
	Data    *ResponseProductUpdate	`json:"data"`
	Error   bool					`json:"error"`
}
func (c *ControllerProduct) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		idParam, err := request.PathLastParam(r)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "invalid path param", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}
		id, err := strconv.Atoi(idParam)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "parameter must be int", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// process
		// -> get searched product by id
		pr, err := c.storage.GetOne(id)
		if err != nil {
			var code int; var body *ResponseBodyUpdate
			switch {
			case errors.Is(err, storage.ErrStorageProductNotFound):
				code = http.StatusNotFound
				body = &ResponseBodyUpdate{Message: "product not found", Data: nil, Error: true}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBodyUpdate{Message: "internal error", Data: nil, Error: true}
			}

			response.JSON(w, code, body)
			return
		}
		// -- serialization
		product := &RequestProductUpdate{
			Name:   pr.Name,
			Type:	pr.Type,
			Count:	pr.Count,
			Price:	pr.Price,
		}

		// -> patch product to RequestProductUpdate(filled with original data)
		err = request.JSON(r, product)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "invalid json", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}
		// -- deserialization
		prUpdate := &storage.Product{
			ID:		id,
			Name:   product.Name,
			Type:	product.Type,
			Count:	product.Count,
			Price:	product.Price,
		}
		// -- update product
		err = c.storage.Update(prUpdate)
		if err != nil {
			var code int; var body *ResponseBody
			switch {
			case errors.Is(err, storage.ErrStorageProductNotFound):
				code = http.StatusNotFound
				body = &ResponseBody{Message: "product not found", Data: nil, Error: true}
			case errors.Is(err, storage.ErrStorageProductNotUnique):
				code = http.StatusBadRequest
				body = &ResponseBody{Message: "product not unique", Data: nil, Error: true}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBody{Message: "internal error", Data: nil, Error: true}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseBodyUpdate{
			Message: "success",
			Data: &ResponseProductUpdate{	// serialization
				Name:   prUpdate.Name,
				Type:	prUpdate.Type,
				Count:	prUpdate.Count,
				Price:	prUpdate.Price,
			},
			Error: false,
		}

		response.JSON(w, code, body)
	}
}

// Delete deletes product by id
type ResponseBodyDelete struct {
	Message string	`json:"message"`
	Data    any		`json:"data"`
	Error   bool	`json:"error"`
}
func (c *ControllerProduct) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		idParam, err := request.PathLastParam(r)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "invalid path param", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}
		id, err := strconv.Atoi(idParam)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "parameter must be int", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// process
		// -> delete product by id
		err = c.storage.Delete(id)
		if err != nil {
			var code int; var body *ResponseBody
			switch {
			case errors.Is(err, storage.ErrStorageProductNotFound):
				code = http.StatusNotFound
				body = &ResponseBody{Message: "product not found", Data: nil, Error: true}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBody{Message: "internal error", Data: nil, Error: true}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseBodyDelete{
			Message: "success",
			Data:    nil,
			Error:   false,
		}

		response.JSON(w, code, body)
	}
}