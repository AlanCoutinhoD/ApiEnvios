package infrastructure

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQService struct {
	channel *amqp.Channel
}

func NewRabbitMQService(channel *amqp.Channel) *RabbitMQService {
	return &RabbitMQService{
		channel: channel,
	}
}

func (s *RabbitMQService) DeclareQueue(queueName string) error {
	_, err := s.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("error declarando cola %s: %v", queueName, err)
	}
	return nil
}

func (s *RabbitMQService) PublishMessage(queueName string, body string, headers map[string]interface{}) error {
	// Asegurar que la cola existe
	if err := s.DeclareQueue(queueName); err != nil {
		return err
	}

	// Confirmar publicaciones
	if err := s.channel.Confirm(false); err != nil {
		return fmt.Errorf("error activando confirmaciones de publicación: %v", err)
	}

	confirms := s.channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	err := s.channel.Publish(
		"",        // exchange
		queueName, // routing key
		true,      // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:        []byte(body),
			Headers:     headers,
			DeliveryMode: amqp.Persistent, // Mensajes persistentes
		})

	if err != nil {
		return fmt.Errorf("error publicando mensaje: %v", err)
	}

	// Esperar confirmación
	if confirmed := <-confirms; !confirmed.Ack {
		return fmt.Errorf("mensaje no confirmado por el broker")
	}

	log.Printf("✅ Mensaje publicado exitosamente en la cola: %s", queueName)
	return nil
}

func (s *RabbitMQService) ConsumeMessages(queueName string, handler func([]byte, map[string]interface{}) error) error {
	// Asegurar que la cola existe
	if err := s.DeclareQueue(queueName); err != nil {
		return err
	}

	msgs, err := s.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("error configurando consumidor: %v", err)
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body, msg.Headers); err != nil {
				log.Printf("❌ Error procesando mensaje: %v", err)
				msg.Nack(false, true) // Rechazar mensaje y reencolar
			} else {
				msg.Ack(false) // Confirmar procesamiento exitoso
			}
		}
	}()

	log.Printf("✅ Consumidor iniciado para la cola: %s", queueName)
	return nil
}
