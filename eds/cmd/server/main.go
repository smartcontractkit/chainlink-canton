package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/smartcontractkit/chainlink-canton/eds/internal/api"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/config"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/disclosure"
	"github.com/smartcontractkit/chainlink-canton/eds/internal/ledger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	envConfig, err := config.LoadEnvironments(cfg.EnvironmentsConfigPath)
	if err != nil {
		log.Fatalf("failed to load environments config: %v", err)
	}

	ledgerClient, err := ledger.NewClient(
		fmt.Sprintf("%s:%d", cfg.LedgerAPIHost, cfg.LedgerAPIPort),
		cfg.JWTSecret,
		cfg.JWTAudience,
	)
	if err != nil {
		log.Fatalf("failed to create ledger client: %v", err)
	}
	defer ledgerClient.Close()

	disclosureSvc := disclosure.NewService(ledgerClient, envConfig)
	router := api.NewRouter(disclosureSvc, envConfig)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Println("shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	log.Printf("starting explicit disclosure server on %s:%d", cfg.Host, cfg.Port)
	log.Printf("loaded %d environments: %v", len(envConfig.Environments), envConfig.EnvironmentNames())
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
