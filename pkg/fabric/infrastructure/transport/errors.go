package transport

import (
	"errors"
	log "github.com/sirupsen/logrus"
	commonErrors "storeservice/pkg/common/errors"
	appErrors "storeservice/pkg/fabric/application/errors"
)

var FabricNotExistError = errors.New("fabric: fabric not exist")
var InternalError = commonErrors.InternalError

func WrapError(err error) error {
	switch err {
	case nil:
		return nil
	case appErrors.FabricNotExistError:
		return FabricNotExistError
	case commonErrors.InternalError:
		return InternalError
	default:
		log.Error(err)
		return InternalError
	}
}
