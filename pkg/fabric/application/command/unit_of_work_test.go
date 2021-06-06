package command

import (
	"github.com/google/uuid"
	"storeservice/pkg/fabric/model"
)

type mockUnitOfWork struct {
	fabrics map[string]model.Fabric
}

func (m *mockUnitOfWork) Execute(f func(rp model.FabricRepository) error) error {
	return f(m)
}

func (m *mockUnitOfWork) Store(o model.Fabric) error {
	if m.fabrics == nil {
		m.fabrics = make(map[string]model.Fabric)
	}

	m.fabrics[o.ID.String()] = o
	return nil
}

func (m *mockUnitOfWork) Get(orderUuid uuid.UUID) (*model.Fabric, error) {
	if m.fabrics == nil {
		m.fabrics = make(map[string]model.Fabric)
	}

	o, ok := m.fabrics[orderUuid.String()]
	if ok {
		return &o, nil
	}

	return nil, nil
}
