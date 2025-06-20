package controllers

import (
	"net/http"

	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/satryarangga/amartha-loan-engine/services"

	"github.com/gin-gonic/gin"
)

type LoanController struct {
	loanService *services.LoanServiceImpl
}

func NewLoanController(loanService *services.LoanServiceImpl) *LoanController {
	return &LoanController{
		loanService: loanService,
	}
}

// CreateLoan godoc
// @Summary Create a new loan
// @Description Create a new loan with automatic schedule generation
// @Tags loans
// @Accept json
// @Produce json
// @Param loan body models.Loan true "Loan object"
// @Success 201 {object} map[string]interface{} "Created"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /loans [post]
func (c *LoanController) CreateLoan(ctx *gin.Context) {
	var loan models.Loan
	if err := ctx.ShouldBindJSON(&loan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := c.loanService.CreateLoan(ctx, &loan); err != nil {
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

// GetLoans godoc
// @Summary Get all loans
// @Description Retrieve a list of all loans
// @Tags loans
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /loans [get]
func (c *LoanController) GetLoans(ctx *gin.Context) {
	loans, err := c.loanService.GetLoans(ctx)
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

// GetLoanByID godoc
// @Summary Get loan by ID
// @Description Retrieve a specific loan by its ID
// @Tags loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Router /loans/{id} [get]
func (c *LoanController) GetLoanByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Loan ID is required",
		})
		return
	}

	loan, err := c.loanService.GetLoanByID(ctx, id)
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

// UpdateLoan godoc
// @Summary Update a loan
// @Description Update an existing loan's information
// @Tags loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID"
// @Param loan body models.Loan true "Updated loan object"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /loans/{id} [put]
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

	if err := c.loanService.UpdateLoan(ctx, &loan); err != nil {
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

// DeleteLoan godoc
// @Summary Delete a loan
// @Description Delete a loan by its ID
// @Tags loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID"
// @Success 200 {object} map[string]interface{} "Success"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Router /loans/{id} [delete]
func (c *LoanController) DeleteLoan(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Loan ID is required",
		})
		return
	}

	// Note: DeleteLoan method was removed from service as it's not in CommonRepository
	ctx.JSON(http.StatusNotImplemented, gin.H{
		"error": "Delete functionality not implemented",
	})
}
