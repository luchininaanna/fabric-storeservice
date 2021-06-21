package transport

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
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

type addFabricRequest struct {
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
	Cost   float32 `json:"cost"`
}

type updateFabricRequest struct {
	ID     string  `json:"ID"`
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
	Cost   float32 `json:"cost"`
}

type addFabricResponse struct {
	Id string `json:"id"`
}

func Router(db *sql.DB) http.Handler {
	srv := &server{
		repository.NewUnitOfWork(db),
		queryImpl.NewFabricQueryService(db),
	}
	r := mux.NewRouter()

	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/fabric", srv.addFabric).Methods(http.MethodPost)
	s.HandleFunc("/fabric", srv.updateFabric).Methods(http.MethodPut)
	s.HandleFunc("/fabrics", srv.getFabrics).Methods(http.MethodGet)
	return cmd.LogMiddleware(r)
}

func (s *server) addFabric(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Can't read request body with error")
		return
	}

	defer infrastructure.LogError(r.Body.Close())

	var request addFabricRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Can't parse json response with error")
		return
	}

	var h = command.NewAddFabricCommandHandler(s.unitOfWork)
	id, err := h.Handle(command.AddFabricCommand{
		Name:   request.Name,
		Amount: request.Amount,
		Cost:   request.Cost,
	})
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
		log.Error("Can't read request body with error")
	}

	defer infrastructure.LogError(r.Body.Close())

	var request updateFabricRequest
	err = json.Unmarshal(b, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Error("Can't parse json response with error")
	}

	var h = command.NewUpdateFabricCommandHandler(s.unitOfWork)

	fabricUid, err := uuid.Parse(request.ID)
	if err != nil {
		log.Error("Can't parse fabric uid")
	}

	err = h.Handle(command.UpdateFabricCommand{
		ID:     fabricUid,
		Name:   request.Name,
		Cost:   request.Cost,
		Amount: request.Amount,
	})
	if err != nil {
		http.Error(w, WrapError(err).Error(), http.StatusBadRequest)
		return
	}
}
