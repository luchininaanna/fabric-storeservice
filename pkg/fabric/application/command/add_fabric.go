package command

import (
	"github.com/google/uuid"
	"storeservice/pkg/fabric/model"
	"time"
)

type AddFabricCommand struct {
	Name   string
	Amount float32
	Cost   float32
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
	var fabricId *uuid.UUID
	err := h.unitOfWork.Execute(func(rp model.FabricRepository) error {

		fabric, err := model.NewFabric(uuid.New(), c.Name, c.Cost, c.Amount, time.Now(), nil)
		if err != nil {
			return err
		}

		fabricId = &fabric.ID

		return rp.Store(fabric)
	})

	return fabricId, err
}
