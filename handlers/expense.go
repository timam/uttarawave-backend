package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"github.com/timam/uttarawave-backend/repositories"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type ExpenseHandler struct {
	expenseRepo repositories.ExpenseRepository
}

func NewExpenseHandler(er repositories.ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{
		expenseRepo: er,
	}
}

func (h *ExpenseHandler) CreateExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		var expense models.Expense
		if err := c.ShouldBindJSON(&expense); err != nil {
			logger.Error("Failed to bind JSON", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		expense.ID = uuid.New().String()
		expense.PaidAt = time.Now()

		err := h.expenseRepo.CreateExpense(c.Request.Context(), &expense)
		if err != nil {
			logger.Error("Failed to create expense", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
			return
		}

		c.JSON(http.StatusCreated, expense)
	}
}

// Add other expense handler methods here (Get, GetAll, etc.)
