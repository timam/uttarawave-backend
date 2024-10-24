package response

import (
	"github.com/timam/uttarawave-backend/models"
	"time"
)

type CustomerResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type CustomerListResponse struct {
	Items      []CustomerItemResponse `json:"items"`
	Pagination PaginationInfo         `json:"pagination"`
}

type CustomerItemResponse struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Mobile               string    `json:"mobile"`
	Email                string    `json:"email"`
	Type                 string    `json:"type"`
	IdentificationNumber string    `json:"identificationNumber"`
	Address              Address   `json:"address"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type Address struct {
	ID         string  `json:"id"`
	House      string  `json:"house"`
	Road       string  `json:"road"`
	Block      string  `json:"block"`
	Area       string  `json:"area"`
	BuildingID *string `json:"buildingId,omitempty"`
}

func NewCustomerResponse(status int, message string, data interface{}) CustomerResponse {
	return CustomerResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewCustomerItemResponse(customer *models.Customer) CustomerItemResponse {
	response := CustomerItemResponse{
		ID:     customer.ID,
		Name:   customer.Name,
		Mobile: customer.Mobile,
		Type:   string(customer.Type),
		Address: Address{
			ID:         customer.Address.ID,
			House:      customer.Address.House,
			Road:       customer.Address.Road,
			Block:      customer.Address.Block,
			Area:       customer.Address.Area,
			BuildingID: customer.Address.BuildingID,
		},
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}

	if customer.Email != nil {
		response.Email = *customer.Email
	}

	if customer.IdentificationNumber != nil {
		response.IdentificationNumber = *customer.IdentificationNumber
	}

	return response
}

func NewCustomerListResponse(customers []models.Customer, total int64, page, size int) CustomerListResponse {
	customerResponses := make([]CustomerItemResponse, len(customers))
	for i, customer := range customers {
		customerResponses[i] = NewCustomerItemResponse(&customer)
	}

	return CustomerListResponse{
		Items: customerResponses,
		Pagination: PaginationInfo{
			Total: total,
			Page:  page,
			Size:  size,
		},
	}
}
