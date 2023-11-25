package question

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
	questionRepository "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/question/repository/postgresql"
	questionUsecase "github.com/go-park-mail-ru/2023_2_Hamster/internal/microservices/question/usecase"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
)

const (
	maxHeaderBytesHTTP = 1 << 20
	readTimeoutHTTP    = 10 * time.Second
	writeTimeoutHTTP   = 10 * time.Second
)

func main() {
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

	connString := "host=HammyWallet_QUESTION port=5436 user=hamster dbname=HammyWallets password=56748 sslmode=disable"
	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}

	questionRepo := questionRepository.NewRepository(pool)

	questionUsecase := questionUsecase.NewUsecase(questionRepo)

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

	questionProto.RegistrationQuestionServer(server, questionGRPS.NewQuestionGRPS(questionUsecase, log))

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
}
