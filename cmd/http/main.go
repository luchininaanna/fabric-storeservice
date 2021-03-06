package main

import (
	"context"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"net/http"
	"storeservice/pkg/common/cmd"
	"storeservice/pkg/fabric/infrastructure/transport"
)

const appID = "store"

type config struct {
	cmd.WebConfig
	cmd.DatabaseConfig
}

func main() {
	var conf config
	if err := envconfig.Process(appID, &conf); err != nil {
		log.Fatal(err)
	}

	cmd.SetupLogger()

	killSignalChan := cmd.GetKillSignalChan()
	srv := startServer(&conf)

	cmd.WaitForKillSignal(killSignalChan)
	log.Fatal(srv.Shutdown(context.Background()))
}

func startServer(conf *config) *http.Server {
	log.WithFields(log.Fields{"port": conf.ServerPort}).Info("starting the store server")
	db := cmd.CreateDBConnection(conf.DatabaseConfig)

	router := transport.Router(db)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", conf.ServerPort), Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
		log.Fatal(db.Close())
	}()

	return srv
}
