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
// @Param loan body models.LoanRequest true "Loan object"
// @Success 201 {object} models.Loan "Created"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Router /loans [post]
func (c *LoanController) CreateLoan(ctx *gin.Context) {
	var loan models.LoanRequest
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

// GetLoanByID godoc
// @Summary Get loan by ID
// @Description Retrieve a specific loan by its ID
// @Tags loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID"
// @Success 200 {object} models.Loan "Success"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 404 {object} models.ErrorResponse "Not Found"
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
