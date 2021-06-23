package transport

import (
	"context"
	"database/sql"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	storeservice "storeservice/api"
	"storeservice/pkg/common/cmd"
	"storeservice/pkg/common/infrastructure/transport"
	"storeservice/pkg/fabric/application/command"
	"storeservice/pkg/fabric/application/query"
	queryImpl "storeservice/pkg/fabric/infrastructure/query"
	"storeservice/pkg/fabric/infrastructure/repository"
)

type server struct {
	unitOfWork command.UnitOfWork
	fqs        query.FabricQueryService
}

func (s *server) AddFabric(_ context.Context, request *storeservice.AddFabricRequest) (*storeservice.AddFabricResponse, error) {
	var h = command.NewAddFabricCommandHandler(s.unitOfWork)
	id, err := h.Handle(command.AddFabricCommand{
		Name:   request.Name,
		Amount: request.Amount,
		Cost:   request.Cost,
	})
	if err != nil {
		return nil, err
	}

	return &storeservice.AddFabricResponse{Id: id.String()}, nil
}

func (s *server) UpdateFabric(_ context.Context, request *storeservice.UpdateFabricRequest) (*empty.Empty, error) {
	var h = command.NewUpdateFabricCommandHandler(s.unitOfWork)

	fabricUid, err := uuid.Parse(request.Id)
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
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *server) GetFabrics(_ context.Context, empty *empty.Empty) (*storeservice.FabricsResponse, error) {
	fabrics, err := s.fqs.GetFabrics()
	if err != nil {
		return nil, err
	}

	var fabricsResponseList []*storeservice.FabricsResponse_FabricResponse
	for _, fabric := range fabrics {
		fabricsResponseList = append(fabricsResponseList, &storeservice.FabricsResponse_FabricResponse{
			FabricId: fabric.ID,
			Name:     fabric.Name,
			Amount:   fabric.Amount,
			Cost:     fabric.Cost,
		})
	}

	response := storeservice.FabricsResponse{
		Fabrics: fabricsResponseList,
	}

	return &response, nil
}

func Router(db *sql.DB) http.Handler {
	srv := Server(db)

	router := transport.NewServeMux()
	err := storeservice.RegisterStoreServiceHandlerServer(context.Background(), router, srv)
	if err != nil {
		log.Fatal(err)
	}

	return cmd.LogMiddleware(router)
}

func Server(db *sql.DB) storeservice.StoreServiceServer {
	return &server{
		repository.NewUnitOfWork(db),
		queryImpl.NewFabricQueryService(db),
	}
}
