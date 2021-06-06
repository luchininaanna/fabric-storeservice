package command

import "storeservice/pkg/fabric/model"

type UnitOfWork interface {
	Execute(func(rp model.FabricRepository) error) error
}
