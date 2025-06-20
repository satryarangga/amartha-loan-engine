package controllers

import (
	"amartha/models"
	"amartha/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoanController struct {
	loanService *services.LoanService
}

func NewLoanController(loanService *services.LoanService) *LoanController {
	return &LoanController{
		loanService: loanService,
	}
}

// CreateLoan handles POST /api/v1/loans
func (c *LoanController) CreateLoan(ctx *gin.Context) {
	var loan models.Loan
	if err := ctx.ShouldBindJSON(&loan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := c.loanService.CreateLoan(&loan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create loan",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data":    loan,
		"message": "Loan created successfully",
	})
}

// GetLoans handles GET /api/v1/loans
func (c *LoanController) GetLoans(ctx *gin.Context) {
	loans, err := c.loanService.GetLoans()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch loans",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  loans,
		"count": len(loans),
	})
}

// GetLoanByID handles GET /api/v1/loans/:id
func (c *LoanController) GetLoanByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Loan ID is required",
		})
		return
	}

	loan, err := c.loanService.GetLoanByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Loan not found",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": loan,
	})
}

// UpdateLoan handles PUT /api/v1/loans/:id
func (c *LoanController) UpdateLoan(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Loan ID is required",
		})
		return
	}

	var loan models.Loan
	if err := ctx.ShouldBindJSON(&loan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := c.loanService.UpdateLoan(&loan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update loan",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    loan,
		"message": "Loan updated successfully",
	})
}

// DeleteLoan handles DELETE /api/v1/loans/:id
func (c *LoanController) DeleteLoan(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Loan ID is required",
		})
		return
	}

	if err := c.loanService.DeleteLoan(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to delete loan",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Loan deleted successfully",
	})
}
