package main

import (
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
	storeservice "storeservice/api"
	"storeservice/pkg/common/cmd"
	transportUtils "storeservice/pkg/common/infrastructure/transport"
	"storeservice/pkg/fabric/infrastructure/transport"
)

const appID = "store"

type config struct {
	cmd.GRPCConfig
	cmd.DatabaseConfig
}

func main() {
	var c config
	if err := envconfig.Process(appID, &c); err != nil {
		log.Fatal(err)
	}

	cmd.SetupLogger()

	killSignalChan := cmd.GetKillSignalChan()
	startServer(killSignalChan, c)
}

func startServer(killSignalCh <-chan os.Signal, c config) {
	l := transportUtils.ListenTCP(c.ServerPort)
	defer transportUtils.CloseService(l, "socket")

	db := cmd.CreateDBConnection(c.DatabaseConfig)
	defer transportUtils.CloseService(db, "db connection")

	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	storeservice.RegisterStoreServiceServer(grpcServer, transport.Server(db))

	go func() {
		log.WithField("port", c.ServerPort).Info("starting the server")
		if err := grpcServer.Serve(l); err != nil {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	cmd.WaitForKillSignal(killSignalCh)
}
