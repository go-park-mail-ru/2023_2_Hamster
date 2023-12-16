package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	generatedAuth "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc/generated"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	authHandler "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc"
	authRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/repository/postgresql"
	authUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/usecase"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	log := logger.NewLogger(ctx)
	db, err := postgresql.InitPostgresDB(ctx)
	if err != nil {
		log.Errorf("Error Initializing PostgreSQL database: %v", err)
		return
	}
	defer func() {
		db.Close()

		log.Info("Db closed without errors")
	}()

	log.Info("Db connection successfully")

	authRepo := authRep.NewRepository(db, *log)

	authUsecase := authUsecase.NewUsecase(authRepo, *log)

	service := authHandler.NewAuthGRPC(authUsecase, *log)

	srv, ok := net.Listen("tcp", ":8010")
	if ok != nil {
		log.Fatalln("can't listen port", err)
	}

	metricsMw := middleware.NewMetricsMiddleware()
	metricsMw.Register(middleware.ServiceAuthName)

	server := grpc.NewServer(grpc.UnaryInterceptor(metricsMw.ServerMetricsInterceptor))

	generatedAuth.RegisterAuthServiceServer(server, service)
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	http.Handle("/", r)
	httpSrv := http.Server{Handler: r, Addr: ":8011"}

	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil {
			fmt.Print(err)
		}
	}()

	fmt.Print("creator running on: ", srv.Addr())
	return server.Serve(srv)

}
