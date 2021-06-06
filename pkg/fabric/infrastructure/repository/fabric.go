package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"storeservice/pkg/common/infrastructure"
	"storeservice/pkg/fabric/model"
)

type fabricRepository struct {
	tx *sql.Tx
}

func (or *fabricRepository) Store(fabric model.Fabric) error {
	_, err := or.tx.Exec(""+
		`INSERT INTO fabric (id, name, amount, cost)
		 VALUES (UUID_TO_BIN(?), ?, ?, ?)`, fabric.ID, fabric.Name, fabric.Amount, fabric.Cost)

	if err != nil {
		err = infrastructure.InternalError(err)
	}

	return err
}

func (or *fabricRepository) Get(fabricUuid uuid.UUID) (*model.Fabric, error) {
	fabricIdBin, err := fabricUuid.MarshalBinary()
	if err != nil {
		return nil, err
	}

	rows, err := or.tx.Query(""+
		`SELECT
		   BIN_TO_UUID(f.id) AS id,
		   f.name AS name,
		   f.amount AS amount,
		   f.cost AS cost,
		 FROM fabric f
		 WHERE f.id = ?`, fabricIdBin)

	if err != nil {
		return nil, err
	}

	if rows.Next() {
		order, err := parseFabric(rows)
		if err != nil {
			return nil, err
		}

		return order, nil
	}

	return nil, err
}

func parseFabric(r *sql.Rows) (*model.Fabric, error) {
	var fabricId string
	var name string
	var amount int
	var cost int

	err := r.Scan(&fabricId, &name, &amount, &cost)
	if err != nil {
		return nil, err
	}

	fabricUuid, err := uuid.Parse(fabricId)
	if err != nil {
		return nil, err
	}

	order, err := model.NewFabric(fabricUuid, name, amount, cost)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
