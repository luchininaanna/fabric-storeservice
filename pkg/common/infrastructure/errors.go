package infrastructure

import (
	log "github.com/sirupsen/logrus"
	"storeservice/pkg/common/errors"
)

func InternalError(e error) error {
	log.Error(e)
	return errors.InternalError
}
