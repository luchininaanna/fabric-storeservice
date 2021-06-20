package repository

import (
	"database/sql"
	"github.com/google/uuid"
	"storeservice/pkg/common/infrastructure"
	"storeservice/pkg/fabric/model"
	"time"
)

type fabricRepository struct {
	tx *sql.Tx
}

func (or *fabricRepository) Store(fabric model.Fabric) error {
	_, err := or.tx.Exec(""+
		`INSERT INTO fabric (id, name, amount, cost, created_at, updated_at)
		 VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?) `+
		"ON DUPLICATE KEY UPDATE `name` = ?, `amount` = ?, `cost` = ?, `created_at` = ?, `updated_at` = ?",
		fabric.ID, fabric.Name, fabric.Amount, fabric.Cost, fabric.CreatedAt, fabric.UpdatedAt,
		fabric.Name, fabric.Amount, fabric.Cost, fabric.CreatedAt, fabric.UpdatedAt)

	if err != nil {
		err = infrastructure.InternalError(err)
	}

	return err
}

func (or *fabricRepository) Get(fabricUuid uuid.UUID) (*model.Fabric, error) {
	rows, err := or.tx.Query(""+
		`SELECT
		   BIN_TO_UUID(f.id) AS id,
		   f.name AS name,
		   f.amount AS amount,
		   f.cost AS cost,
           f.created_at AS created_at,
		   f.updated_at AS updated_at
		 FROM fabric f
		 WHERE f.id = UUID_TO_BIN(?)`, fabricUuid)

	if err != nil {
		return nil, err
	}

	defer infrastructure.CloseRows(rows)

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
	var amount float32
	var cost float32
	var createdAt time.Time
	var updatedAtNullable sql.NullTime

	err := r.Scan(&fabricId, &name, &amount, &cost, &createdAt, &updatedAtNullable)
	if err != nil {
		return nil, err
	}

	fabricUuid, err := uuid.Parse(fabricId)
	if err != nil {
		return nil, err
	}

	var updatedAt *time.Time
	if updatedAtNullable.Valid {
		updatedAt = &updatedAtNullable.Time
	}

	order, err := model.NewFabric(fabricUuid, name, amount, cost, createdAt, updatedAt)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
