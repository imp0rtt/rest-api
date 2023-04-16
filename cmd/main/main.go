package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest-api/internal/config"
	"rest-api/internal/user"
	"rest-api/internal/user/db"
	"rest-api/pkg/client/mongodb"
	"rest-api/pkg/logging"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Infof("Create router")
	router := httprouter.New()

	cfg := config.GetConfig()
	cfgMongo := cfg.MongoDB
	mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username, cfgMongo.Password, cfgMongo.Database, cfgMongo.AuthDB)
	if err != nil {
		panic(err)
	}
	//brew services start mongodb-community@6.0
	//brew services stop mongodb-community@6.0

	storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)
	service := user.NewService(storage, logger)
	logger.Info("Register user handler")
	handler := user.NewHandler(logger, service)
	handler.Register(router)
	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("Start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("socket path %s", socketPath)

		logger.Info("listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		logger.Infof("Server is listening unix socket: %s", socketPath)

	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s: %s", cfg.Listen.BindIp, cfg.Listen.Port))
		logger.Infof("Server is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
