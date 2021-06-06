package model

import (
	"github.com/google/uuid"
)

type Fabric struct {
	ID     uuid.UUID
	Name   string
	Cost   int
	Amount int
}

type FabricRepository interface {
	Store(o Fabric) error
	Get(fabricUuid uuid.UUID) (*Fabric, error)
}

func NewFabric(fabricUuid uuid.UUID, name string, cost int, amount int) (Fabric, error) {

	if name == "" {
		return Fabric{}, FabricWithoutNameError
	}

	if cost <= 0 {
		return Fabric{}, InvalidFabricCostError
	}

	if amount <= 0 {
		return Fabric{}, InvalidFabricAmountError
	}

	return Fabric{
		fabricUuid,
		name,
		cost,
		amount,
	}, nil
}
