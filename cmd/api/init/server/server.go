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
	if err := godotenv.Load(); err != nil {
		return cfg, err
	}

	host, existHost := os.LookupEnv("SERVER_HOST")
	port, existPort := os.LookupEnv("SERVER_PORT")
	if !existHost || !existPort {
		return cfg, errors.New("existHost or existPort is Empty")
	}

	cfg = cfgServer{
		ServerPort: host,
		ServerHost: port,
	}

	return cfg, nil
}

func (s *Server) Run(handler http.Handler) error {
	cfgSer, err := initServerConfigFromEnv()
	if err != nil {
		return err
	}

	s.httpServer = &http.Server{
		Addr:    cfgSer.ServerHost + ":" + cfgSer.ServerPort,
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
