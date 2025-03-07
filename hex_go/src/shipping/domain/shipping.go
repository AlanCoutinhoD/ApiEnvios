package domain

import (
	"errors"
	"time"
)

type Shipping struct {
	ID        int64     `json:"id"`
	IdUser    string    `json:"idUser"`
	IdProduct int64     `json:"idProduct"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

// Validar los datos del envío
func (s *Shipping) Validate() error {
	if s.IdUser == "" {
		return errors.New("el ID de usuario es requerido")
	}
	if s.IdProduct <= 0 {
		return errors.New("el ID de producto debe ser mayor a 0")
	}
	if s.Quantity <= 0 {
		return errors.New("la cantidad debe ser mayor a 0")
	}
	return nil
}

type ShippingRepository interface {
	Save(shipping *Shipping) error
	GetByID(id int64) (*Shipping, error)
	GetByUserID(userID int64) ([]*Shipping, error)
}

type MessageQueue interface {
	DeclareQueue(queueName string) error
	PublishMessage(queueName string, body string, headers map[string]interface{}) error
	ConsumeMessages(queueName string, handler func([]byte, map[string]interface{}) error) error
}
