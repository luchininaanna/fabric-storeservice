package repository

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"storeservice/pkg/fabric/application/command"
	"storeservice/pkg/fabric/model"
)

type unitOfWork struct {
	db *sql.DB
}

func NewUnitOfWork(db *sql.DB) command.UnitOfWork {
	return &unitOfWork{db}
}

func (u *unitOfWork) Execute(job func(rp model.FabricRepository) error) error {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}

	err = job(&fabricRepository{tx})

	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Error(err2)
		}
	} else {
		err2 := tx.Commit()
		if err2 != nil {
			log.Error(err2)
		}
	}

	return err
}
