package controllers

import (
	"net/http"

	"github.com/satryarangga/amartha-loan-engine/models"
	"github.com/satryarangga/amartha-loan-engine/services"

	"github.com/gin-gonic/gin"
)

type BorrowerController struct {
	borrowerService *services.BorrowerServiceImpl
}

func NewBorrowerController(borrowerService *services.BorrowerServiceImpl) *BorrowerController {
	return &BorrowerController{
		borrowerService: borrowerService,
	}
}

// GetBorrowers godoc
// @Summary Get all borrowers
// @Description Retrieve a list of all borrowers
// @Tags borrowers
// @Accept json
// @Produce json
// @Success 200 {object} models.Borrower "Success"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /borrowers [get]
func (c *BorrowerController) GetBorrowers(ctx *gin.Context) {
	borrowers, err := c.borrowerService.GetBorrowers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch borrowers",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  borrowers,
		"count": len(borrowers),
	})
}

// GetBorrowerByID godoc
// @Summary Get borrower by ID
// @Description Retrieve a specific borrower by their ID
// @Tags borrowers
// @Accept json
// @Produce json
// @Param id path string true "Borrower ID"
// @Success 200 {object} models.Borrower "Success"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Router /borrowers/{id} [get]
func (c *BorrowerController) GetBorrowerByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Borrower ID is required",
		})
		return
	}

	borrower, err := c.borrowerService.GetBorrowerByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   "Borrower not found",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": borrower,
	})
}

// CreateBorrower godoc
// @Summary Create a new borrower
// @Description Create a new borrower with the provided information
// @Tags borrowers
// @Accept json
// @Produce json
// @Param borrower body models.BorrowerRequest true "Borrower object"
// @Success 201 {object} models.Borrower "Created"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Router /borrowers [post]
func (c *BorrowerController) CreateBorrower(ctx *gin.Context) {
	var borrower models.Borrower
	if err := ctx.ShouldBindJSON(&borrower); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := c.borrowerService.CreateBorrower(ctx, &borrower); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create borrower",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data":    borrower,
		"message": "Borrower created successfully",
	})
}
