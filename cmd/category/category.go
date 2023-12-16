package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	generatedCategory "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/delivery/grpc/generated"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	categoryHandler "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/delivery/grpc"
	categoryRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/repository/postgres"
	categoryUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/category/usecase"
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

	categoryRepo := categoryRep.NewRepository(db, *log)

	categoryUsecase := categoryUsecase.NewUsecase(categoryRepo, *log)

	service := categoryHandler.NewCategoryGRPC(categoryUsecase, *log)

	srv, ok := net.Listen("tcp", ":8030")
	if ok != nil {
		log.Fatalln("can't listen port", err)
	}

	metricsMw := middleware.NewMetricsMiddleware()
	metricsMw.Register(middleware.ServiceCategoryName)

	server := grpc.NewServer(grpc.UnaryInterceptor(metricsMw.ServerMetricsInterceptor))

	generatedCategory.RegisterCategoryServiceServer(server, service)
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	http.Handle("/", r)
	httpSrv := http.Server{Handler: r, Addr: ":8031"}

	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil {
			fmt.Print(err)
		}
	}()

	fmt.Print("creator running on: ", srv.Addr())
	return server.Serve(srv)

}
