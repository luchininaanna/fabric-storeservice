package model

import (
	"github.com/google/uuid"
	"time"
)

type Fabric struct {
	ID        uuid.UUID
	Name      string
	Cost      float32
	Amount    float32
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type FabricRepository interface {
	Store(f Fabric) error
	Get(fabricUuid uuid.UUID) (*Fabric, error)
}

func NewFabric(fabricUuid uuid.UUID, name string, cost float32, amount float32, createdAt time.Time, updatedAt *time.Time) (Fabric, error) {

	if name == "" {
		return Fabric{}, FabricWithoutNameError
	}

	if cost <= 0 {
		return Fabric{}, InvalidFabricCostError
	}

	if amount < 0 {
		return Fabric{}, InvalidFabricAmountError
	}

	return Fabric{
		fabricUuid,
		name,
		cost,
		amount,
		createdAt,
		updatedAt,
	}, nil
}
