package transport

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"storeservice/pkg/common/cmd"
	"storeservice/pkg/common/infrastructure"
	"storeservice/pkg/fabric/application/command"
	"storeservice/pkg/fabric/application/query"
	queryImpl "storeservice/pkg/fabric/infrastructure/query"
	"storeservice/pkg/fabric/infrastructure/repository"
)

type server struct {
	unitOfWork command.UnitOfWork
	fqs        query.FabricQueryService
}

func Router(db *sql.DB) http.Handler {
	srv := &server{
		repository.NewUnitOfWork(db),
		queryImpl.NewFabricQueryService(db),
	}
	r := mux.NewRouter()

	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/order", srv.addFabric).Methods(http.MethodPost)
	s.HandleFunc("/order}", srv.updateFabric).Methods(http.MethodPut)
	s.HandleFunc("/order/{ID:[0-9a-zA-Z-]+}", srv.getFabrics).Methods(http.MethodGet)
	return cmd.LogMiddleware(r)
}

type addFabricResponse struct {
	Id string `json:"id"`
}

func (s *server) addFabric(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Can't read request body with error")
		return
	}

	defer infrastructure.LogError(r.Body.Close())

	var addFabricCommand command.AddFabricCommand
	err = json.Unmarshal(b, &addFabricCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Can't parse json response with error")
		return
	}

	var h = command.NewAddFabricCommandHandler(s.unitOfWork)
	id, err := h.Handle(addFabricCommand)
	if err != nil {
		http.Error(w, WrapError(err).Error(), http.StatusBadRequest)
		return
	}

	RenderJson(w, &addFabricResponse{id.String()})
}

func (s *server) updateFabric(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Can't read request body with error")
	}

	defer r.Body.Close()

	var updateFabricCommand command.UpdateFabricCommand
	err = json.Unmarshal(b, &updateFabricCommand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal("Can't parse json response with error")
	}

	var h = command.NewUpdateFabricCommandHandler(s.unitOfWork)
	err = h.Handle(updateFabricCommand)
	if err != nil {
		http.Error(w, WrapError(err).Error(), http.StatusBadRequest)
		return
	}
}
