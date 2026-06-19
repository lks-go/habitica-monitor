// Command monitor — точка входа сервиса мониторинга Habitica.
// Собирает зависимости (composition root), поднимает HTTP-сервер
// и фоновый планировщик снапшотов.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/habitica-monitor/configs"
	"github.com/example/habitica-monitor/internal/application"
	"github.com/example/habitica-monitor/internal/infrastructure/habitica"
	"github.com/example/habitica-monitor/internal/infrastructure/sqlite"
	httpiface "github.com/example/habitica-monitor/internal/interfaces/http"
)

func main() {
	cfg := configs.Load()

	// --- Infrastructure ---
	db, err := sqlite.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close()

	userRepo := sqlite.NewUserRepository(db)
	statsRepo := sqlite.NewStatsRepository(db)
	habiticaClient := habitica.NewClient(cfg.XClient)

	// --- Application ---
	userSvc := application.NewUserService(userRepo)
	statsSvc := application.NewStatsService(statsRepo)
	snapshotSvc := application.NewSnapshotService(userRepo, statsRepo, habiticaClient, cfg.SnapshotInterval)

	// --- Lifecycle ---
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Фоновый планировщик снапшотов.
	go snapshotSvc.Run(ctx)

	// --- HTTP ---
	handler := httpiface.NewHandler(userSvc, statsSvc)
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpiface.CORS(cfg.CORSOrigin, handler.Routes()),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("http: слушаю %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown: завершаю работу...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
	os.Exit(0)
}
