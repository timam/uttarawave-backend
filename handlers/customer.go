package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/pkg/response"
	"github.com/timam/uttarawave-backend/repositories"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type CustomerHandler struct {
	repo         repositories.CustomerRepository
	buildingRepo repositories.BuildingRepository
}

func NewCustomerHandler(
	cr repositories.CustomerRepository,
	br repositories.BuildingRepository) *CustomerHandler {
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
			c.JSON(http.StatusBadRequest, response.NewCustomerResponse(http.StatusBadRequest, "Invalid input", nil))
			return
		}

		if customerInput.Mobile == "" || customerInput.Name == "" {
			logger.Warn("Missing required fields")
			c.JSON(http.StatusBadRequest, response.NewCustomerResponse(http.StatusBadRequest, "Mobile and Name are required fields", nil))
			return
		}

		customer := models.Customer{
			ID:                   uuid.New().String(),
			Mobile:               customerInput.Mobile,
			Name:                 customerInput.Name,
			Email:                customerInput.Email,
			Type:                 models.CustomerType(customerInput.Type),
			IdentificationNumber: customerInput.IdentificationNumber,
			Address: models.Address{
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
			c.JSON(http.StatusInternalServerError, response.NewCustomerResponse(http.StatusInternalServerError, "Failed to save customer data", nil))
			return
		}

		logger.Info("Customer created successfully", zap.String("id", customer.ID), zap.String("mobile", customer.Mobile))

		customerItemResponse := response.NewCustomerItemResponse(&customer)
		c.JSON(http.StatusCreated, response.NewCustomerResponse(http.StatusCreated, "Customer created successfully", customerItemResponse))
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
			c.JSON(http.StatusInternalServerError, response.NewCustomerResponse(http.StatusInternalServerError, "Failed to get customer data", nil))
			return
		}
		if customer == nil {
			c.JSON(http.StatusNotFound, response.NewCustomerResponse(http.StatusNotFound, "Customer not found", nil))
			return
		}

		customerItemResponse := response.NewCustomerItemResponse(customer)
		c.JSON(http.StatusOK, response.NewCustomerResponse(http.StatusOK, "Customer retrieved successfully", customerItemResponse))
	}
}

func (h *CustomerHandler) GetAllCustomers() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

		customers, totalCount, err := h.repo.GetCustomersPaginated(page, pageSize)
		if err != nil {
			logger.Error("Failed to get customers", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewCustomerResponse(http.StatusInternalServerError, "Failed to get customers", nil))
			return
		}

		customerListResponse := response.NewCustomerListResponse(customers, int64(totalCount), page, pageSize)
		c.JSON(http.StatusOK, response.NewCustomerResponse(http.StatusOK, "Customers retrieved successfully", customerListResponse))
	}
}

func (h *CustomerHandler) UpdateCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		logger.Info("UpdateCustomer called", zap.String("id", id))

		if id == "" {
			logger.Warn("Customer ID not provided")
			c.JSON(http.StatusBadRequest, response.NewCustomerResponse(http.StatusBadRequest, "Customer ID must be provided", nil))
			return
		}

		var updateData models.Customer
		if err := c.ShouldBindJSON(&updateData); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, response.NewCustomerResponse(http.StatusBadRequest, "Invalid input", nil))
			return
		}

		existingCustomer, err := h.repo.GetCustomer(id)
		if err != nil {
			logger.Error("Failed to find customer", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewCustomerResponse(http.StatusInternalServerError, "Failed to find customer", nil))
			return
		}
		if existingCustomer == nil {
			logger.Warn("Customer not found", zap.String("id", id))
			c.JSON(http.StatusNotFound, response.NewCustomerResponse(http.StatusNotFound, "Customer not found", nil))
			return
		}

		if updateData.Mobile != "" && updateData.Mobile != existingCustomer.Mobile {
			c.JSON(http.StatusBadRequest, response.NewCustomerResponse(http.StatusBadRequest, "Mobile number cannot be updated", nil))
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
				c.JSON(http.StatusInternalServerError, response.NewCustomerResponse(http.StatusInternalServerError, "Failed to get building details", nil))
				return
			}
			existingCustomer.Address = building.Address
			existingCustomer.Address.ID = uuid.New().String()
			existingCustomer.Address.CustomerID = &existingCustomer.ID
			existingCustomer.Address.BuildingID = updateData.Address.BuildingID
		} else {
			if updateData.Address.House == "" || updateData.Address.Road == "" || updateData.Address.Block == "" || updateData.Address.Area == "" {
				c.JSON(http.StatusBadRequest, response.NewCustomerResponse(http.StatusBadRequest, "All address fields (House, Road, Block, Area) are required when BuildingID is not provided", nil))
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
			c.JSON(http.StatusInternalServerError, response.NewCustomerResponse(http.StatusInternalServerError, "Failed to update customer data", nil))
			return
		}

		logger.Info("Customer updated successfully", zap.String("id", existingCustomer.ID), zap.String("name", existingCustomer.Name))
		customerItemResponse := response.NewCustomerItemResponse(existingCustomer)
		c.JSON(http.StatusOK, response.NewCustomerResponse(http.StatusOK, "Customer updated successfully", customerItemResponse))
	}
}

func (h *CustomerHandler) DeleteCustomer() gin.HandlerFunc {
	return func(c *gin.Context) {
		mobile := c.Query("mobile")

		if mobile == "" {
			c.JSON(http.StatusBadRequest, response.NewCustomerResponse(http.StatusBadRequest, "Mobile number must be provided", nil))
			return
		}

		customer, err := h.repo.GetCustomerByMobile(mobile)
		if err != nil {
			logger.Error("Failed to find customer by mobile", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewCustomerResponse(http.StatusInternalServerError, "Failed to find customer by mobile", nil))
			return
		}
		if customer.ID == "" {
			c.JSON(http.StatusNotFound, response.NewCustomerResponse(http.StatusNotFound, "Customer not found", nil))
			return
		}

		if err := h.repo.DeleteCustomer(customer.ID); err != nil {
			logger.Error("Failed to delete customer by mobile", zap.Error(err))
			c.JSON(http.StatusInternalServerError, response.NewCustomerResponse(http.StatusInternalServerError, "Failed to delete customer", nil))
			return
		}

		c.JSON(http.StatusOK, response.NewCustomerResponse(http.StatusOK, "Customer deleted successfully", nil))
	}
}
