package handler

import (
	"encoding/json"
	"github/rossi1/go-api-microservice-example/domain/entity"
	"github/rossi1/go-api-microservice-example/service"
	"net/http"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs"

	"github.com/gofrs/uuid"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

var _, _ = uuid.NewV4()

type Service interface {
	CreateCategory(name string, product []entity.Product) (*entity.Category, error)
	DeleteCategory(categoryID string) error
	UpdateCategory(name, categoryID string) (*entity.Category, error)
	GetCategory(categoryID string) (*entity.Category, error)
	CreateProduct(product entity.Product, categoryID string) (*entity.Product, error)
	GetAllCategories(r *http.Request) ([]entity.Category, error)
	DeleteProduct(productID string) error
	GetProduct(productID string) (*entity.Product, error)
	UpdateProduct(productID string, params map[string]interface{}) (*entity.Product, error)
	GetAllProducts(categoryID string, r *http.Request) ([]entity.Product, error)
}

type Handler struct {
	Service
}

func NewHandler(db *gorm.DB, search interface{}, pub interface{}, topic string) *Handler {
	srv := service.New(db, search, pub, topic)
	return &Handler{srv}
}

// getCategory godoc
// @Summary get a category item
// @Description get string by categoryId
// @ID get-string-by-uuid
// @Accept  json
// @Produce  json
// @Param categoryId path uuid.UUID true "category Id"
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /category/{categoryId} [get]
func (h *Handler) getCategoryAPIView(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "categoryId")

	result, err := h.GetCategory(categoryID)
	if err != nil {
		ErrorResponse(w, err.Error(), 400)
		return
	}
	Response(w, SuccessResponse{
		Message: "record fetched",
		Data:    result,
	}, http.StatusOK)
}

func (h *Handler) getAllCategoriesAPIView(w http.ResponseWriter, r *http.Request) {
	result, err := h.GetAllCategories(r)
	if err != nil {
		ErrorResponse(w, err.Error(), 500)
		return
	}
	Response(w, SuccessResponse{
		Message: "records fetched",
		Data:    result,
	}, http.StatusOK)
}

func (h *Handler) deleteCategoryAPIView(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "categoryId")

	err := h.DeleteCategory(categoryID)

	if err != nil {
		ErrorResponse(w, err.Error(), 500)
		return
	}
	Response(w, SuccessResponse{
		Message: "record deleted",
		Data:    nil,
	}, http.StatusNoContent)
}

func (h *Handler) createCategoryAPIView(w http.ResponseWriter, r *http.Request) {
	var (
		req     CategorySerializer
		product []entity.Product
	)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, err.Error(), 400)
		return
	}

	defer r.Body.Close()

	if len(req.Product) > 0 {
		for _, val := range req.Product {
			product = append(product, entity.Product{Name: val.Name, Description: val.Description, Discount: val.Discount,
				Tax: val.Tax, BarCode: val.BarCode, Expires: val.Expires, Weight: val.Weight, Image: val.Image})
		}
	}

	result, err := h.CreateCategory(req.Name, product)

	if err != nil {
		ErrorResponse(w, err.Error(), 500)
		return
	}
	Response(w, SuccessResponse{
		Message: "record created",
		Data:    result,
	}, http.StatusCreated)
}

func (h *Handler) updateCategoryAPIView(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "categoryId")

	var req CategorySerializer

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, err.Error(), 400)
		return
	}
	defer r.Body.Close()

	result, err := h.UpdateCategory(req.Name, categoryID)
	if err != nil {
		ErrorResponse(w, err.Error(), 404)
		return
	}

	Response(w, SuccessResponse{
		Message: "record updated",
		Data:    result,
	}, http.StatusOK)

}

func (h *Handler) getAllProductsAPIView(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "categoryId")
	result, err := h.GetAllProducts(categoryID, r)
	if err != nil {
		ErrorResponse(w, err.Error(), 500)
		return
	}
	Response(w, SuccessResponse{
		Message: "records fetched",
		Data:    result,
	}, http.StatusOK)
}

func (h *Handler) deleteProductAPIView(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")

	err := h.DeleteProduct(productID)

	if err != nil {
		ErrorResponse(w, err.Error(), 500)
		return
	}
	Response(w, SuccessResponse{
		Message: "record deleted",
		Data:    nil,
	}, http.StatusNoContent)
}

func (h *Handler) getProductAPIView(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")

	result, err := h.GetProduct(productID)
	if err != nil {
		ErrorResponse(w, err.Error(), 404)
		return
	}
	Response(w, SuccessResponse{
		Message: "record fetched",
		Data:    result,
	}, http.StatusOK)
}

func (h *Handler) createProductAPIView(w http.ResponseWriter, r *http.Request) {

	categoryID := chi.URLParam(r, "categoryId")

	var (
		product entity.Product
	)

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		ErrorResponse(w, err.Error(), 400)
		return
	}

	defer r.Body.Close()

	result, err := h.CreateProduct(product, categoryID)

	if err != nil {
		ErrorResponse(w, err.Error(), 500)
		return
	}
	Response(w, SuccessResponse{
		Message: "record created",
		Data:    result,
	}, http.StatusCreated)
}

func (h *Handler) updateProductAPIView(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")

	var (
		req ProductSerializer
	)

	obj := map[string]interface{}{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ErrorResponse(w, err.Error(), 400)
		return
	}

	defer r.Body.Close()

	obj["name"] = req.Name
	obj["tax"] = req.Tax
	obj["description"] = req.Description
	obj["weight"] = req.Weight
	obj["expires"] = req.Expires
	obj["bar_code"] = req.BarCode
	obj["discount"] = req.Discount
	obj["image"] = req.Image

	result, err := h.UpdateProduct(productID, obj)
	if err != nil {
		ErrorResponse(w, err.Error(), 404)
		return
	}

	Response(w, SuccessResponse{
		Message: "record updated",
		Data:    result,
	}, http.StatusOK)

}
