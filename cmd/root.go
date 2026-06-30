/*
 * Copyright 2026 MuixStudio
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/muixstudio/clio/internal/api"
	"github.com/muixstudio/clio/internal/driver"
	"github.com/muixstudio/clio/internal/driver/config"
	"github.com/muixstudio/clio/internal/infra/database"
	"github.com/muixstudio/clio/internal/infra/logger"
	"github.com/muixstudio/clio/internal/persistence/gormstore"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "zus",
	Short:         "zus",
	Long:          `zus`,
	RunE:          run,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func run(cmd *cobra.Command, args []string) error {
	cfgPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("failed to get config flag: %w", err)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	logger := logger.New(cfg.Log)

	db, err := database.NewDB(cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() {
		if err := database.CloseDB(db); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	store, err := gormstore.NewGormStore(db)
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	log.Println("persistence: PostgreSQL")

	d := driver.New(
		cfg,
		driver.WithStore(store),
		driver.WithLogger(logger),
	)

	engine, err := api.NewApiEngine(d)

	if oidc := d.OIDCStrategy(); oidc != nil {
		oidc.RegisterRoutes(engine)
	}

	if err != nil {
		return fmt.Errorf("failed to initialize engine: %w", err)
	}

	dispatchCtx, stopDispatch := context.WithCancel(context.Background())
	defer stopDispatch()
	if err := d.StartAlertDispatch(dispatchCtx); err != nil {
		return fmt.Errorf("failed to start alert dispatch: %w", err)
	}
	defer d.StopAlertDispatch()

	//Create HTTP server
	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: engine,
	}

	// Start server in goroutine
	go func() {
		log.Printf("⇨ Gin server starting on %s", cfg.Addr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	//Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	//Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")

	return nil
}

func init() {
	rootCmd.Flags().StringP("config", "c", "", "config path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
