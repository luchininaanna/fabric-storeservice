package command

import (
	"github.com/google/uuid"
	"storeservice/pkg/fabric/application/errors"
	"storeservice/pkg/fabric/model"
	"time"
)

type UpdateFabricCommand struct {
	ID     uuid.UUID
	Name   string
	Cost   float32
	Amount float32
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
		if err != nil {
			return err
		}
		if fabric == nil {
			return errors.FabricNotExistError
		}

		t := time.Now()
		fabric.UpdatedAt = &t
		fabric.Name = c.Name
		fabric.Cost = c.Cost
		fabric.Amount = c.Amount

		return rp.Store(*fabric)
	})

	return err
}
