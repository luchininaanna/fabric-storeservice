package query

import (
	"database/sql"
	"storeservice/pkg/common/infrastructure"
	"storeservice/pkg/fabric/application/errors"
	"storeservice/pkg/fabric/application/query"
	"storeservice/pkg/fabric/application/query/data"
	"time"
)

func NewFabricQueryService(db *sql.DB) query.FabricQueryService {
	return &fabricQueryService{db: db}
}

type fabricQueryService struct {
	db *sql.DB
}

func (qs *fabricQueryService) GetFabric(id string) (*data.FabricData, error) {
	rows, err := qs.db.Query(""+
		getSelectFabricSQL()+
		"WHERE f.id = UUID_TO_BIN(?)", id)

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	if rows.Next() {
		fabric, err := parseFabric(rows)
		if err != nil {
			return nil, infrastructure.InternalError(err)
		}

		return fabric, nil
	}

	return nil, errors.FabricNotExistError
}

func (qs *fabricQueryService) GetFabrics() ([]data.FabricData, error) {
	rows, err := qs.db.Query(getSelectFabricSQL())

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	var fabrics []data.FabricData
	for rows.Next() {
		fabric, err := parseFabric(rows)
		if err != nil {
			return nil, infrastructure.InternalError(err)
		}

		fabrics = append(fabrics, *fabric)
	}

	return fabrics, nil
}

func getSelectFabricSQL() string {
	return `SELECT
		   BIN_TO_UUID(f.id) AS id,
           f.name AS name,
           f.amount AS amount,
           f.cost AS cost,
           f.created_at AS created_at,
		   f.updated_at AS updated_at
		FROM fabric f`
}

func parseFabric(r *sql.Rows) (*data.FabricData, error) {
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

	var updatedAt *time.Time
	if updatedAtNullable.Valid {
		updatedAt = &updatedAtNullable.Time
	}

	return &data.FabricData{
		ID:        fabricId,
		Name:      name,
		Amount:    amount,
		Cost:      cost,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
