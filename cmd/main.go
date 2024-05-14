package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"psPro-task/internal/app"
	"psPro-task/internal/config"
	"psPro-task/internal/delivery"
	"psPro-task/internal/repository"
	"psPro-task/internal/service"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	/*
		if err := gotenv.Load(); err != nil {
			logrus.Fatalf("Ошибка при получении переменных окружения %s", err.Error())
		}
	*/
	dbCfg := config.GetDBConfig()
	db, err := repository.OpenDB(dbCfg)
	if err != nil {
		logrus.Fatalf("Ошибка при подклюении к базе данных: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := delivery.NewHandler(services)

	srv := new(app.Server)
	go func() {
		if err := srv.Run(os.Getenv("HTTP_PORT"), handlers.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Ошибка при работе http-сервера: %s", err.Error())
		}
	}()

	logrus.Println("Сервер запуска команд начал работы")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("Сервер запуска команд завершил работу")
	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Произошла ошибка при завершении работы сервера: %s", err.Error())
	}
	if err = db.Close(); err != nil {
		logrus.Errorf("Ошибка при отсоединении от базы данных: %s", err.Error())
	}
}
