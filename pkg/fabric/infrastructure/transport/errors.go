package transport

import (
	"errors"
	log "github.com/sirupsen/logrus"
	commonErrors "storeservice/pkg/common/errors"
	appErrors "storeservice/pkg/fabric/application/errors"
	modelErrors "storeservice/pkg/fabric/model"
)

var InternalError = commonErrors.InternalError
var FabricNotExistError = errors.New("fabric: fabric not exist")
var FabricWithoutNameError = errors.New("fabric: fabric without name")
var InvalidFabricCostError = errors.New("fabric: invalid fabric cost")
var InvalidFabricAmountError = errors.New("fabric: invalid fabric account")

func WrapError(err error) error {
	switch err {
	case nil:
		return nil
	case commonErrors.InternalError:
		return InternalError
	case appErrors.FabricNotExistError:
		return FabricNotExistError
	case modelErrors.FabricWithoutNameError:
		return FabricWithoutNameError
	case modelErrors.InvalidFabricCostError:
		return InvalidFabricCostError
	case modelErrors.InvalidFabricAmountError:
		return InvalidFabricAmountError
	default:
		log.Error(err)
		return InternalError
	}
}
