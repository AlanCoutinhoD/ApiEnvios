package application

import "demo/src/shipping/domain"

type ShippingUseCase struct {
	repository   domain.ShippingRepository
	messageQueue domain.MessageQueue
}

func NewShippingUseCase(repo domain.ShippingRepository, mq domain.MessageQueue) *ShippingUseCase {
	return &ShippingUseCase{
		repository:   repo,
		messageQueue: mq,
	}
}

func (uc *ShippingUseCase) CreateShipping(shipping *domain.Shipping) error {
	// Guardar en base de datos
	err := uc.repository.Save(shipping)
	if err != nil {
		return err
	}

	// Enviar mensaje a la cola
	headers := map[string]interface{}{
		"idUser": shipping.IdUser,
	}

	err = uc.messageQueue.PublishMessage(
		"orderNotification_queue",
		shipping.IdUser,
		headers,
	)
	if err != nil {
		return err
	}

	return nil
}
