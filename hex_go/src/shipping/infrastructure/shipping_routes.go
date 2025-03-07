package infrastructure

import (
	"demo/src/shipping/infrastructure/controllers"

	"github.com/gin-gonic/gin"
)

type ShippingRouter struct {
	createShippingController *controllers.CreateShippingController
}

func NewShippingRouter(createShippingController *controllers.CreateShippingController) *ShippingRouter {
	return &ShippingRouter{
		createShippingController: createShippingController,
	}
}

func (sr *ShippingRouter) SetupRoutes(router *gin.Engine) {
	shippingGroup := router.Group("/shipping")
	{
		shippingGroup.POST("", sr.createShippingController.Execute)
	}
}
