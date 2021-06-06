package command

import (
	"github.com/google/uuid"
	"storeservice/pkg/fabric/application/errors"
	"storeservice/pkg/fabric/model"
)

type UpdateFabricCommand struct {
	ID     uuid.UUID
	Name   string
	Cost   int
	Amount int
}

type updateFabricCommandHandler struct {
	unitOfWork UnitOfWork
}

type UpdateFabricCommandHandler interface {
	Handle(command UpdateFabricCommand) error
}

func NewUpdateFabricCommandHandler(unitOfWork UnitOfWork) UpdateFabricCommandHandler {
	return &updateFabricCommandHandler{unitOfWork}
}

func (h *updateFabricCommandHandler) Handle(c UpdateFabricCommand) error {
	err := h.unitOfWork.Execute(func(rp model.FabricRepository) error {

		fabric, err := rp.Get(c.ID)
		if fabric == nil {
			return errors.FabricNotExistError
		}

		order, err := model.NewFabric(c.ID, c.Name, c.Cost, c.Amount)
		if err != nil {
			return err
		}

		return rp.Store(order)
	})

	return err
}
