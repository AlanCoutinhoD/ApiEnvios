package infrastructure

import (
	"database/sql"
	"demo/src/shipping/application"
	"demo/src/shipping/infrastructure/controllers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

var db *sql.DB
var conn *amqp.Connection
var ch *amqp.Channel

func Init(router *gin.Engine) {
	var err error

	// Inicializar conexiones
	db, err = InitMySQL()
	if err != nil {
		log.Fatal(err)
	}

	conn, ch, err = InitRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}

	// Inicializar componentes de Shipping
	shippingRepo := NewMySQLRepository(db)
	messagingService := NewRabbitMQService(ch)
	shippingUseCase := application.NewShippingUseCase(shippingRepo, messagingService)
	createShippingController := controllers.NewCreateShippingController(shippingUseCase)
	shippingRouter := NewShippingRouter(createShippingController)

	// Configurar rutas
	shippingRouter.SetupRoutes(router)
}
