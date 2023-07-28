package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gosidekick/goconfig"
	"github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/core/service"
	storagemysql "github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/driven/storage/mysql"
	"github.com/izzdalfk/norma-research-pi-server-umkm-app/internal/driver/rest"
	"github.com/jmoiron/sqlx"
)

func main() {
	var cfg config
	if err := goconfig.Parse(&cfg); err != nil {
		log.Fatalf("unable to parse app config: %v", err)
	}
	// init. mysql storage
	dbConn, err := sqlx.ConnectContext(context.Background(), "mysql", cfg.SQLDSN)
	handleError(err, fmt.Sprintf("unable to connect to mysql database due: %v", err))

	strg, err := storagemysql.NewStorage(storagemysql.StorageConfig{
		DBClient: dbConn,
	})
	handleError(err, fmt.Sprintf("unable to initialize mysql storage due: %v", err))

	// init. service
	svc, err := service.NewService(service.ServiceConfig{
		Storage:        strg,
		SupportService: &mockSupportService{},
	})
	handleError(err, fmt.Sprintf("unable to initialize core service due: %v", err))

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// init. server API handler
	api, err := rest.NewAPI(rest.APIConfig{
		Service: svc,
	})
	handleError(err, fmt.Sprintf("unable to initialize rest api due: %v", err))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: api.Handler(),
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func handleError(err error, message string) {
	if err != nil {
		log.Fatalf(message)
	}
}

type config struct {
	SQLDSN string `cfg:"db_sqldsn" cfgRequired:"true" cfgDefault:"root:test1234@tcp(localhost:23306)/umkm?timeout=5s"`
}

type mockSupportService struct{}

func (m *mockSupportService) CalculateDeliveryPrice(ctx context.Context) (float64, error) {
	return 0, nil
}

func (m *mockSupportService) PickupDelivery(ctx context.Context) (bool, error) {
	return true, nil
}
