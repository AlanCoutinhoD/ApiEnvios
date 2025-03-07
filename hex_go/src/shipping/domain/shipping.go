package domain

import (
	
)

type Shipping struct {
	ID        int64     `json:"id"`
	IdUser    string    `json:"idUser"`
	IdProduct int64     `json:"idProduct"`
	Quantity  int       `json:"quantity"`
}

// Validar los datos del envío


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
