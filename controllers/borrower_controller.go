package controllers

import (
	"amartha/models"
	"amartha/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BorrowerController struct {
	borrowerService *services.BorrowerService
}

func NewBorrowerController(borrowerService *services.BorrowerService) *BorrowerController {
	return &BorrowerController{
		borrowerService: borrowerService,
	}
}

// GetBorrowers handles GET /api/v1/borrowers
func (c *BorrowerController) GetBorrowers(ctx *gin.Context) {
	borrowers, err := c.borrowerService.GetBorrowers()
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

// GetBorrowerByID handles GET /api/v1/borrowers/:id
func (c *BorrowerController) GetBorrowerByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Borrower ID is required",
		})
		return
	}

	borrower, err := c.borrowerService.GetBorrowerByID(id)
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

// CreateBorrower handles POST /api/v1/borrowers
func (c *BorrowerController) CreateBorrower(ctx *gin.Context) {
	var borrower models.Borrower
	if err := ctx.ShouldBindJSON(&borrower); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := c.borrowerService.CreateBorrower(&borrower); err != nil {
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

// UpdateBorrower handles PUT /api/v1/borrowers/:id
func (c *BorrowerController) UpdateBorrower(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Borrower ID is required",
		})
		return
	}

	var borrower models.Borrower
	if err := ctx.ShouldBindJSON(&borrower); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := c.borrowerService.UpdateBorrower(&borrower); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update borrower",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":    borrower,
		"message": "Borrower updated successfully",
	})
}

// DeleteBorrower handles DELETE /api/v1/borrowers/:id
func (c *BorrowerController) DeleteBorrower(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Borrower ID is required",
		})
		return
	}

	if err := c.borrowerService.DeleteBorrower(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to delete borrower",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Borrower deleted successfully",
	})
}
