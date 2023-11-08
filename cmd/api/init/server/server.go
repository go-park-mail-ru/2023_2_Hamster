package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	maxHeaderBytes = 1 << 20
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
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

func (s *Server) Init(handler http.Handler) error {
	cfgSer, err := initServerConfigFromEnv()
	if err != nil {
		return err
	}

	addr := cfgSer.ServerHost + ":" + cfgSer.ServerPort
	s.httpServer = &http.Server{
		Addr:           addr,
		Handler:        handler,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	}

	return nil
}

func (s *Server) Run() error {
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
