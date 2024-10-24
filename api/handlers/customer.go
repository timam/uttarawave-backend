package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	models2 "github.com/timam/uttarawave-backend/internals/models"
	repositories2 "github.com/timam/uttarawave-backend/internals/repositories"
	"github.com/timam/uttarawave-backend/internals/response"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CustomerHandler struct {
	repo         repositories2.CustomerRepository
	buildingRepo repositories2.BuildingRepository
}

func NewCustomerHandler(
	cr repositories2.CustomerRepository,
	br repositories2.BuildingRepository) *CustomerHandler {
	return &CustomerHandler{
		repo:         cr,
		buildingRepo: br,
	}
}

func (h *CustomerHandler) CreateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		var customerInput struct {
			Mobile               string  `json:"mobile"`
			Name                 string  `json:"name"`
			Email                *string `json:"email,omitempty"`
			Type                 string  `json:"type,omitempty"`
			IdentificationNumber *string `json:"identificationNumber,omitempty"`
			Flat                 string  `json:"flat,omitempty"`
			House                string  `json:"house"`
			Road                 string  `json:"road"`
			Block                string  `json:"block"`
			Area                 string  `json:"area"`
			City                 string  `json:"city,omitempty"`
		}

		if err := c.ShouldBindJSON(&customerInput); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
			return
		}

		if customerInput.Mobile == "" || customerInput.Name == "" {
			err := errors.New("mobile and name are required fields")
			logger.Warn("Missing required fields")
			response.Error(c, http.StatusBadRequest, "Missing required fields", err.Error())
			return
		}

		// Validate customer type
		if customerInput.Type == "" {
			err := errors.New("customer type is required")
			logger.Warn("Missing customer type")
			response.Error(c, http.StatusBadRequest, "Customer type is required", err.Error())
			return
		}

		if customerInput.Type != string(models2.Individual) && customerInput.Type != string(models2.Business) {
			err := errors.New("invalid customer type")
			logger.Warn("Invalid customer type", zap.String("type", customerInput.Type))
			response.Error(c, http.StatusBadRequest, "Invalid customer type", err.Error())
			return
		}

		customer := models2.Customer{
			ID:                   uuid.New().String(),
			Mobile:               customerInput.Mobile,
			Name:                 customerInput.Name,
			Email:                customerInput.Email,
			Type:                 models2.CustomerType(customerInput.Type),
			IdentificationNumber: customerInput.IdentificationNumber,
			Address: models2.Address{
				ID:    uuid.New().String(),
				Flat:  customerInput.Flat,
				House: customerInput.House,
				Road:  customerInput.Road,
				Block: customerInput.Block,
				Area:  customerInput.Area,
				City:  customerInput.City,
			},
		}

		customer.Address.CustomerID = &customer.ID

		err := h.repo.CreateCustomer(c.Request.Context(), &customer)
		if err != nil {
			logger.Error("Failed to save customer data", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to save customer data", err.Error())
			return
		}

		logger.Info("Customer created successfully", zap.String("id", customer.ID), zap.String("mobile", customer.Mobile))

		customerItemResponse := response.NewCustomerItemResponse(&customer)
		response.Success(c, http.StatusCreated, "Customer created successfully", customerItemResponse)
	}
}

func (h *CustomerHandler) GetCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		mobile := c.Query("mobile")

		if mobile == "" {
			h.GetAllCustomers()(c)
			return
		}

		customer, err := h.repo.GetCustomerByMobile(mobile)
		if err != nil {
			logger.Error("Failed to get customer data", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to get customer data", err.Error())
			return
		}
		if customer == nil {
			err := errors.New("customer not found")
			response.Error(c, http.StatusNotFound, "Customer not found", err.Error())
			return
		}

		customerItemResponse := response.NewCustomerItemResponse(customer)
		response.Success(c, http.StatusOK, "Customer retrieved successfully", customerItemResponse)
	}
}

func (h *CustomerHandler) GetAllCustomers() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

		customers, totalCount, err := h.repo.GetCustomersPaginated(page, pageSize)
		if err != nil {
			logger.Error("Failed to get customers", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to get customers", err.Error())
			return
		}

		customerListResponse := response.NewCustomerListResponse(customers, int64(totalCount), page, pageSize)
		response.Success(c, http.StatusOK, "Customers retrieved successfully", customerListResponse)
	}
}

func (h *CustomerHandler) UpdateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		logger.Info("UpdateCustomer called", zap.String("id", id))

		if id == "" {
			err := errors.New("customer ID must be provided")
			logger.Warn("Customer ID not provided")
			response.Error(c, http.StatusBadRequest, "Customer ID must be provided", err.Error())
			return
		}

		var updateData models2.Customer
		if err := c.ShouldBindJSON(&updateData); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			response.Error(c, http.StatusBadRequest, "Invalid input", err.Error())
			return
		}

		existingCustomer, err := h.repo.GetCustomer(id)
		if err != nil {
			logger.Error("Failed to find customer", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to find customer", err.Error())
			return
		}
		if existingCustomer == nil {
			err := errors.New("customer not found")
			logger.Warn("Customer not found", zap.String("id", id))
			response.Error(c, http.StatusNotFound, "Customer not found", err.Error())
			return
		}

		if updateData.Mobile != "" && updateData.Mobile != existingCustomer.Mobile {
			err := errors.New("mobile number cannot be updated")
			response.Error(c, http.StatusBadRequest, "Mobile number cannot be updated", err.Error())
			return
		}

		existingCustomer.Name = updateData.Name
		existingCustomer.Email = updateData.Email
		existingCustomer.Type = updateData.Type
		existingCustomer.IdentificationNumber = updateData.IdentificationNumber

		if updateData.Address.BuildingID != nil {
			building, err := h.buildingRepo.GetBuildingByID(c.Request.Context(), *updateData.Address.BuildingID)
			if err != nil {
				logger.Error("Failed to get building details", zap.Error(err))
				response.Error(c, http.StatusInternalServerError, "Failed to get building details", err.Error())
				return
			}
			existingCustomer.Address = building.Address
			existingCustomer.Address.ID = uuid.New().String()
			existingCustomer.Address.CustomerID = &existingCustomer.ID
			existingCustomer.Address.BuildingID = updateData.Address.BuildingID
		} else {
			if updateData.Address.House == "" || updateData.Address.Road == "" || updateData.Address.Block == "" || updateData.Address.Area == "" {
				err := errors.New("all address fields (House, Road, Block, Area) are required when BuildingID is not provided")
				response.Error(c, http.StatusBadRequest, "Invalid address input", err.Error())
				return
			}
			existingCustomer.Address = updateData.Address
			existingCustomer.Address.ID = uuid.New().String()
			existingCustomer.Address.CustomerID = &existingCustomer.ID
			existingCustomer.Address.BuildingID = nil
		}

		err = h.repo.UpdateCustomer(existingCustomer)
		if err != nil {
			logger.Error("Failed to update customer data", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to update customer data", err.Error())
			return
		}

		logger.Info("Customer updated successfully", zap.String("id", existingCustomer.ID), zap.String("name", existingCustomer.Name))
		customerItemResponse := response.NewCustomerItemResponse(existingCustomer)
		response.Success(c, http.StatusOK, "Customer updated successfully", customerItemResponse)
	}
}

func (h *CustomerHandler) DeleteCustomer() gin.HandlerFunc {
	//TODO: if customer have any associated subscription, or any device that is not returned; customer cant be deleted.
	return func(c *gin.Context) {
		mobile := c.Query("mobile")

		if mobile == "" {
			err := errors.New("mobile number must be provided")
			response.Error(c, http.StatusBadRequest, "Mobile number must be provided", err.Error())
			return
		}

		customer, err := h.repo.GetCustomerByMobile(mobile)
		if err != nil {
			logger.Error("Failed to find customer by mobile", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to find customer by mobile", err.Error())
			return
		}
		if customer == nil {
			err := errors.New("customer not found")
			response.Error(c, http.StatusNotFound, "Customer not found", err.Error())
			return
		}

		if err := h.repo.DeleteCustomer(customer.ID); err != nil {
			logger.Error("Failed to delete customer by mobile", zap.Error(err))
			response.Error(c, http.StatusInternalServerError, "Failed to delete customer", err.Error())
			return
		}

		response.Success(c, http.StatusOK, "Customer deleted successfully", nil)
	}
}
