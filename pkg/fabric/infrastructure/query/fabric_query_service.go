package query

import (
	"database/sql"
	"storeservice/pkg/common/infrastructure"
	"storeservice/pkg/fabric/application/query"
	"storeservice/pkg/fabric/application/query/data"
)

func NewFabricQueryService(db *sql.DB) query.FabricQueryService {
	return &fabricQueryService{db: db}
}

type fabricQueryService struct {
	db *sql.DB
}

func (qs *fabricQueryService) GetFabric(id string) (*data.FabricData, error) {
	rows, err := qs.db.Query(""+
		getSelectOrderSQL()+
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

	return nil, nil // not found
}

func (qs *fabricQueryService) GetFabrics() ([]data.FabricData, error) {
	rows, err := qs.db.Query(getSelectOrderSQL())

	if err != nil {
		return nil, infrastructure.InternalError(err)
	}
	defer infrastructure.CloseRows(rows)

	var fabrics []data.FabricData
	if rows.Next() {
		fabric, err := parseFabric(rows)
		if err != nil {
			return nil, infrastructure.InternalError(err)
		}

		fabrics = append(fabrics, *fabric)
	}

	return fabrics, nil // not found
}

func getSelectOrderSQL() string {
	return "" +
		`SELECT
		   BIN_TO_UUID(f.id) AS id,
           f.name AS name,
           f.amount AS amount
           f.cost AS cost
		FROM fabric f`
}

func parseFabric(r *sql.Rows) (*data.FabricData, error) {
	var fabricId string
	var name string
	var amount int
	var cost int

	err := r.Scan(&fabricId, &name, &amount, &cost)
	if err != nil {
		return nil, err
	}

	return &data.FabricData{
		ID:     fabricId,
		Name:   name,
		Amount: amount,
		Cost:   cost,
	}, nil
}
