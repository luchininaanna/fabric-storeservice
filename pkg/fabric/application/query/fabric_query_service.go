package query

import "storeservice/pkg/fabric/application/query/data"

type FabricQueryService interface {
	GetFabric(id string) (*data.FabricData, error)
	GetFabrics() ([]data.FabricData, error)
}
