package controllers

import (
	"demo/src/shipping/application"
	"demo/src/shipping/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateShippingController struct {
	useCase *application.ShippingUseCase
}

func NewCreateShippingController(useCase *application.ShippingUseCase) *CreateShippingController {
	return &CreateShippingController{
		useCase: useCase,
	}
}

func (csc *CreateShippingController) Execute(ctx *gin.Context) {
	var shipping domain.Shipping

	if err := ctx.ShouldBindJSON(&shipping); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := csc.useCase.CreateShipping(&shipping)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, shipping)
}
