package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-park-mail-ru/2023_2_Hamster/cmd/api/init/db/postgresql"
	"github.com/go-park-mail-ru/2023_2_Hamster/internal/common/logger"
	authHandler "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc"
	grpcAuth "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/delivery/grpc/generated"
	authRep "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/repository/postgresql"
	authUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/auth/usecase"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	connString := "host=0.0.0.0 port=5436 user=hamster dbname=HammyWallets password=2003 sslmode=disable"
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}

	authRepo := authRep.NewRepository(pool, *log)

	authUse := authUsecase.NewUsecase(authRepo, *log)

	service := authHandler.NewAuthGRPC(authUse, *log)

	listener, err := net.Listen("tcp", "127.0.0.1:8088")
	defer func() {
		if err := listener.Close(); err != nil {
			log.Errorf("Error while closing search tcp listener: %v", err)
		}
	}()

	if err != nil {
		log.Errorf("Can't listen port: %v", err)
		return
	}

	server := grpc.NewServer()

	// httpMetricsServer := &http.Server{
	// 	Handler:        promhttp.HandlerFor(promhttp.HandlerOpts{}),
	// 	Addr:           os.Getenv(config.SearchExporterListenParam),
	// 	MaxHeaderBytes: maxHeaderBytesHTTP,
	// 	ReadTimeout:    readTimeoutHTTP,
	// 	WriteTimeout:   writeTimeoutHTTP,
	// }

	// go func() {
	// 	if err := httpMetricsServer.ListenAndServe(); err != nil {
	// 		log.Errorf("Unable to start a http search metrics server:", err)
	// 	}
	// }()

	grpcAuth.RegisterAuthServiceServer(server, service)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-stop
		log.Info("Server question gracefully shutting down...")

		server.GracefulStop()
	}()

	log.Info("Starting grpc server search")
	if err := server.Serve(listener); err != nil {
		log.Errorf("Server question error: %v", err)
		os.Exit(1)
	}
	wg.Wait()

	return server.Serve(listener)
}
