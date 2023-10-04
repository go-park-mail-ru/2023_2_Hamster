package server

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Server struct {
	httpServer *http.Server
}

type cfgServer struct {
	ServerPort string
	ServerHost string
}

func initServerConfigFromEnv() (cfgServer, error) {
	var cfg = cfgServer{}
	if err := godotenv.Load(".env"); err != nil {
		return cfg, err
	}

	host, existHost := os.LookupEnv("SERVER_HOST")
	port, existPort := os.LookupEnv("SERVER_PORT")
	if !existHost || !existPort {
		return cfg, errors.New("existHost or existPort is Empty")
	}

	cfg = cfgServer{
		ServerHost: host,
		ServerPort: port,
	}

	return cfg, nil
}

func (s *Server) Run(handler http.Handler) error {
	cfgSer, err := initServerConfigFromEnv()
	if err != nil {
		return err
	}

	addr := cfgSer.ServerHost + ":" + cfgSer.ServerPort
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	if s.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
