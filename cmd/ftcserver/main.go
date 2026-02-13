package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rbrabson/ftcstanding/database"
	"github.com/rbrabson/ftcstanding/query"
	"github.com/rbrabson/ftcstanding/server"
	"github.com/spf13/cobra"
)

var (
	port       int
	seasonFlag string
)

// setLogLevelFromEnv sets the log level from the LOG_LEVEL environment variable.
func setLogLevelFromEnv() slog.Level {
	levelStr := os.Getenv("LOG_LEVEL")

	var logLevel slog.Level
	switch strings.ToLower(levelStr) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))
	return logLevel
}

var rootCmd = &cobra.Command{
	Use:     "ftcserver",
	Short:   "FTC Standing HTTP API server",
	Long:    "HTTP REST API server for FTC (FIRST Tech Challenge) standing data including teams, events, matches, awards, and rankings.",
	Example: "  # Start the server on default port 8080\n  ftcserver\n\n  # Start the server on a custom port\n  ftcserver --port 3000\n\n  # Specify a season (optional, can still be provided in API paths)\n  ftcserver --season 2024",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Determine season if provided
		season := seasonFlag
		if season == "" {
			season = os.Getenv("FTC_SEASON")
		}

		var err error
		db, err := database.Init(season)
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		defer db.Close()

		query.Init(db)

		httpServer := server.NewServer(db)

		addr := fmt.Sprintf(":%d", port)
		srv := &http.Server{
			Addr:         addr,
			Handler:      httpServer,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		}

		go func() {
			slog.Info("Starting FTC API server", "address", addr)
			slog.Info("API documentation available at http://localhost:" + fmt.Sprint(port) + "/v1/{season}/{resource}")
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				slog.Error("Server failed to start", "error", err)
				os.Exit(1)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		slog.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}

		slog.Info("Server exited")
		return nil
	},
}

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	rootCmd.Flags().StringVarP(&seasonFlag, "season", "s", "", "Default season year (defaults to FTC_SEASON environment variable)")
}

func main() {
	godotenv.Load()
	setLogLevelFromEnv()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
