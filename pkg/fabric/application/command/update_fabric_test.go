package command

import (
	"github.com/google/uuid"
	"storeservice/pkg/fabric/application/errors"
	"testing"
)

func TestUpdateNotExistingFabric(t *testing.T) {
	uow := &mockUnitOfWork{}
	h := updateFabricCommandHandler{uow}
	_, err := h.Handle(UpdateFabricCommand{
		uuid.New(),
		"silk",
		77,
		5,
	})

	if err != errors.FabricNotExistError {
		t.Error("Update not existing fabric")
	}
}
