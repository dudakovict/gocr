package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/dudakovict/gocr/conf"
	"github.com/dudakovict/gocr/ocr"
)

func main() {
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags|log.Lmicroseconds)

	if err := run(logger); err != nil {
		logger.Fatalf("Error starting the server: %s", err)
	}
}

func run(logger *log.Logger) error {
	cfg := conf.Load()
	logger.Printf("Loaded Configuration: %+v\n", cfg)

	ocrProcessor := ocr.NewOCRProcessor(logger)
	defer ocrProcessor.Close()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("public")))
	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		ocrProcessor.UploadHandler(w, r, cfg)
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		TLSConfig:    &tls.Config{MinVersion: tls.VersionTLS12},
	}

	var wg sync.WaitGroup
	wg.Add(1)

	serverErrors := make(chan error, 1)

	go func() {
		defer wg.Done()
		logger.Printf("Server listening on :%d...\n", cfg.Port)
		serverErrors <- server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case <-sigCh:
		logger.Println("Received termination signal. Initiating graceful shutdown...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Printf("Error during server shutdown: %s", err)
		}

		logger.Println("Server gracefully stopped.")
	}

	wg.Wait()

	return nil
}
