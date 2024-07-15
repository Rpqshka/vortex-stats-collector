package main

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	vortex "vortex-stats-collector"
	"vortex-stats-collector/pkg/handler"
	"vortex-stats-collector/pkg/repository"
	"vortex-stats-collector/pkg/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewClickhouseDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	//migrations
	m, err := migrate.New(
		"file://schema",
		fmt.Sprintf("clickhouse://%s:%s?username=%s&password=%s",
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
			os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")))
	if err != nil {
		logrus.Fatalf("Error creating migration instance: %v", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Error applying migrations: %v", err)
	}
	logrus.Debug("Migrations applied successfully")

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	//start server
	srv := new(vortex.Server)
	go func() {
		if err := srv.Run(os.Getenv("SERV_PORT"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Printf("Vortex App Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Printf("Vortex App Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}
