package command

import (
	"github.com/google/uuid"
	"storeservice/pkg/fabric/model"
)

type AddFabricCommand struct {
	Name   string
	Amount int
	Cost   int
}

type addFabricCommandHandler struct {
	unitOfWork UnitOfWork
}

type AddFabricCommandHandler interface {
	Handle(command AddFabricCommand) (*uuid.UUID, error)
}

func NewAddFabricCommandHandler(unitOfWork UnitOfWork) AddFabricCommandHandler {
	return &addFabricCommandHandler{unitOfWork}
}

func (h *addFabricCommandHandler) Handle(c AddFabricCommand) (*uuid.UUID, error) {
	var orderId *uuid.UUID
	err := h.unitOfWork.Execute(func(rp model.FabricRepository) error {

		order, err := model.NewFabric(uuid.New(), c.Name, c.Cost, c.Amount)
		if err != nil {
			return err
		}

		return rp.Store(order)
	})

	return orderId, err
}
